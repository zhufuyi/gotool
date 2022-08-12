package template

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zhufuyi/pkg/gofile"
)

// Handler 模板导出文件处理接口
type Handler interface {
	SetReplacementFields(fields []Field)
	SetIgnoreFiles(filenames ...string)
	SetOutPath(absPath string, name string) error
	GetOutPath() string
	SaveFiles() error
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	path              string            // 模板目录路径(不包含.或..)是否实际路径
	fs                embed.FS          // 模板目录对应二进制对象
	isActual          bool              // fs字段是否来源实际路径，如果为true，使用ioutil操作文件，如果为false使用fs操作文件
	files             []string          // 模板文件列表
	ignoreFiles       []string          // 忽略替换的文件列表
	replacementFields map[string]string // 从模板文件转为新文件需要替换的字符，key是模板字符串，value是新文件字符串
	outPath           string            // 默认输出目录名称后缀
}

// New 实例化
func New(path string, fs embed.FS) (*TemplateInfo, error) {
	files, err := listFiles(path, fs)
	if err != nil {
		return nil, err
	}

	return &TemplateInfo{
		path:              path,
		fs:                fs,
		isActual:          false,
		files:             files,
		replacementFields: make(map[string]string),
	}, nil
}

// NewSrc 实例化
func NewSrc(path string) (*TemplateInfo, error) {
	files, err := gofile.ListFiles(path)
	if err != nil {
		return nil, err
	}

	path, _ = filepath.Abs(path)
	return &TemplateInfo{
		path:              path,
		isActual:          true,
		files:             files,
		replacementFields: make(map[string]string),
	}, nil
}

// Field 替换字段信息
type Field struct {
	Old             string // 模板字段
	New             string // 新字段
	IsCaseSensitive bool   // 第一个字母是否区分大小写
}

// SetReplacementFields 设置替换字段
func (t *TemplateInfo) SetReplacementFields(fields []Field) {
	for _, v := range fields {
		t.setReplacementField(v.Old, v.New, v.IsCaseSensitive)
	}
}

// SetIgnoreFiles 设置忽略处理的文件
func (t *TemplateInfo) SetIgnoreFiles(filenames ...string) {
	t.ignoreFiles = append(t.ignoreFiles, filenames...)
}

// SetOutPath 设置输出目录路径，优先使用absPath绝对路径，如果absPath为空，自动在当前目录根据参数name和时间生成绝对路径
func (t *TemplateInfo) SetOutPath(absPath string, name string) error {
	if absPath != "" {
		abs, err := filepath.Abs(absPath)
		t.outPath = abs
		return err
	}

	t.outPath = getRunPath() + getPathDelimiter() + name + "_" + time.Now().Format("0102150405")
	return nil
}

// GetOutPath 获取输出目录路径
func (t *TemplateInfo) GetOutPath() string {
	return t.outPath
}

// SaveFiles 导出文件
func (t *TemplateInfo) SaveFiles() error {
	if t.outPath == "" {
		t.outPath = getRunPath() + getPathDelimiter() + "template_" + time.Now().Format("0102150405")
	}

	for _, file := range t.files {
		if t.isIgnoreFile(file) {
			continue
		}

		// 从二进制读取模板文件内容使用embed.FS，如果要从指定目录读取使用ioutil.ReadFile
		var data []byte
		var err error
		if t.isActual {
			data, err = ioutil.ReadFile(file)
		} else {
			data, err = t.fs.ReadFile(file)
		}
		if err != nil {
			return err
		}

		// 获取新文件路径
		newFilePath := t.getNewFilePath(file)
		dir, filename := filepath.Split(newFilePath)

		// 替换文本内容和文件名称
		for oldStr, newStr := range t.replacementFields {
			// 替换文件内容
			data = bytes.ReplaceAll(data, []byte(oldStr), []byte(newStr))

			// 替换文件名
			if strings.Contains(filename, oldStr) {
				newFilePath = dir + strings.ReplaceAll(filename, oldStr, newStr)
			}
		}

		// 保存文件
		err = saveToNewFile(newFilePath, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateInfo) isIgnoreFile(file string) bool {
	isIgnore := false
	_, filename := filepath.Split(file)
	for _, v := range t.ignoreFiles {
		if filename == v {
			isIgnore = true
			break
		}
	}
	return isIgnore
}

// SetReplacementField 设置替换字段，参数isNeedCovert表示是否匹配第一个字母大小写
func (t *TemplateInfo) setReplacementField(oldField string, newField string, isNeedCovert bool) {
	if isNeedCovert && isFirstAlphabet(oldField) {
		// 把第一个字母转为大写
		oldField = strings.ToUpper(oldField[:1]) + oldField[1:]
		newField = strings.ToUpper(newField[:1]) + newField[1:]
		t.replacementFields[oldField] = newField

		// 把第一个字母转为小写
		oldField = strings.ToLower(oldField[:1]) + oldField[1:]
		newField = strings.ToLower(newField[:1]) + newField[1:]
		t.replacementFields[oldField] = newField
	} else {
		t.replacementFields[oldField] = newField
	}
}

func (t *TemplateInfo) getNewFilePath(file string) string {
	var newFilePath string
	if t.isActual {
		newFilePath = t.outPath + strings.Replace(file, t.path, "", 1)
	} else {
		newFilePath = t.outPath + strings.Replace(file, t.path, "", 1)
	}

	if runtime.GOOS == "windows" {
		newFilePath = strings.ReplaceAll(newFilePath, "/", "\\")
	}

	return newFilePath
}

func saveToNewFile(filePath string, data []byte) error {
	// 创建目录
	dir, _ := filepath.Split(filePath)
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}

	// 保存文件
	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// ListFiles 遍历指定目录下所有文件，返回文件的绝对路径
func listFiles(path string, fs embed.FS) ([]string, error) {
	files := []string{}
	err := walkDir(path, &files, fs)
	return files, err
}

// 通过迭代方式遍历文件
func walkDir(dirPath string, allFiles *[]string, fs embed.FS) error {
	files, err := fs.ReadDir(dirPath) // 读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + "/" + file.Name()
		if file.IsDir() {
			walkDir(deepFile, allFiles, fs)
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

// 根据系统类型获取分隔符
func getPathDelimiter() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}

	return delimiter
}

// 获取程序执行的绝对路径
func getRunPath() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Println("os.Executable error.", err.Error())
		return ""
	}

	return filepath.Dir(dir)
}

// 判断字符串第一个字符是字母
func isFirstAlphabet(str string) bool {
	if len(str) == 0 {
		return false
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return true
	}

	return false
}
