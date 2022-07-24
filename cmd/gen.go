package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/pkg/gofile"
)

// api 模板默认值
var apiTemplate = &genTemplate{
	srcPath:      "templates/api",
	pkgName:      "github.com/zhufuyi/goctl/templates/api",
	filterFiles:  []string{"dao.go", "common_code.go", "service.go", "routers.go", "global.go"},
	outputPrefix: "api_",
}

// web 模板默认值
var webTemplate = &genTemplate{
	srcPath:      "templates/web",
	pkgName:      "github.com/zhufuyi/goctl/templates/web",
	filterFiles:  []string{},
	outputPrefix: "web_",
}

var (
	// api的值和web的值共同字段
	templateApiName = "UserExample" // 公共模板的接口名称

	// api的值和web的值不一样字段
	templateSrcPath     = ""         // 模板代码目录
	templatePkgName     = ""         // 模板代码中import包名，需要替换为项目名
	templateFilterFiles = []string{} // 忽略处理的文件
	defaultOutputPrefix = ""         // 输出目录前缀
)

type genTemplate struct {
	srcPath      string
	pkgName      string
	filterFiles  []string
	outputPrefix string
}

func genGinApiCommand() *cobra.Command {
	var (
		resourceArg string
		dstPath     string // 输出目标目录
		apiName     string // 接口名称
		projectName string // 项目名称
	)

	cmd := &cobra.Command{
		Use:   "gen <resource>",
		Short: "Generate gin api code",
		Long: `generate gin api code.

Examples:
    # generate api code
    goctl gen api -p apiExample -a user
    goctl gen api -p apiExample -a user -o /tmp

    # generate web server code
    goctl gen web -p webExample -a user
    goctl gen web -p webExample -a user -o /tmp
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must specify the type of resource to gen. use 'goctl resources' for a complete list of supported resources")
			}
			resourceArg = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			switch resourceArg {
			case api:
				srcPathAbs, dstPathAbs, err := setTemplateDefault(apiTemplate, dstPath, projectName)
				if err != nil {
					return err
				}
				err = runGenApiCommand(&genApiOptions{
					srcPath:     srcPathAbs,
					dstPath:     dstPathAbs,
					apiName:     apiName,
					projectName: projectName,
				})
				if err != nil {
					return err
				}

			case web:
				srcPathAbs, dstPathAbs, err := setTemplateDefault(webTemplate, dstPath, projectName)
				if err != nil {
					return err
				}
				err = runGenWebCommand(&genWebOptions{
					srcPath:     srcPathAbs,
					dstPath:     dstPathAbs,
					apiName:     apiName,
					projectName: projectName,
				})
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown resource name '%s'. Use \"goctl resources\" for a complete list of supported resources.\n", resourceArg)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&apiName, "apiName", "a", "", "api name")
	cmd.MarkFlagRequired("apiName")
	cmd.Flags().StringVarP(&projectName, "projectName", "p", "", "project name")
	cmd.MarkFlagRequired("projectName")
	cmd.Flags().StringVarP(&dstPath, "out", "o", "", "export the code path")

	return cmd
}

// ---------------------------------------------------------------------------------------

type genApiOptions struct {
	srcPath     string
	dstPath     string
	apiName     string
	projectName string
}

func (g *genApiOptions) checkValid() error {
	var err error

	if g.apiName == "" {
		return errors.New("param 'apiName' is empty")
	}
	if g.projectName == "" {
		return errors.New("param 'projectName' is empty")
	}

	if g.srcPath == "" {
		return errors.New("param 'srcPath' is empty")
	} else {
		g.srcPath, err = filepath.Abs(g.srcPath)
		if err != nil {
			return err
		}
	}

	if g.dstPath == "" {
		return errors.New("param 'dstPath' is empty")
	}

	return nil
}

func runGenApiCommand(opts *genApiOptions) error {
	err := opts.checkValid()
	if err != nil {
		return err
	}

	kvs := getReplaceParams(opts.apiName, opts.projectName)
	files, err := gofile.ListFiles(opts.srcPath)
	if err != nil {
		return err
	}
	srcFiles := filterFiles(files, templateFilterFiles...)
	for _, srcFile := range srcFiles {
		dstFile := strings.ReplaceAll(srcFile, opts.srcPath, opts.dstPath)
		err = replaceFile(srcFile, dstFile, kvs)
		if err != nil {
			return err
		}
	}

	fmt.Printf("'%s' api generate successfully, output = %s\n\n", opts.apiName, opts.dstPath)
	return nil
}

// ------------------------------------------------------------------------------------------

type genWebOptions struct {
	srcPath     string
	dstPath     string
	apiName     string
	projectName string
}

func (g *genWebOptions) checkValid() error {
	var err error

	if g.apiName == "" {
		return errors.New("param 'apiName' is empty")
	}
	if g.projectName == "" {
		return errors.New("param 'projectName' is empty")
	}

	if g.srcPath == "" {
		return errors.New("param 'srcPath' is empty")
	} else {
		g.srcPath, err = filepath.Abs(g.srcPath)
		if err != nil {
			return err
		}
	}

	if g.dstPath == "" {
		return errors.New("param 'dstPath' is empty")
	}

	return nil
}

func runGenWebCommand(opts *genWebOptions) error {
	err := opts.checkValid()
	if err != nil {
		return err
	}

	kvs := getReplaceParams(opts.apiName, opts.projectName)
	files, err := gofile.ListFiles(opts.srcPath)
	if err != nil {
		return err
	}
	srcFiles := filterFiles(files, templateFilterFiles...)
	for _, srcFile := range srcFiles {
		dstFile := strings.ReplaceAll(srcFile, opts.srcPath, opts.dstPath)
		err = replaceFile(srcFile, dstFile, kvs)
		if err != nil {
			return err
		}
	}

	fmt.Printf("'%s' web generate successfully, output = %s\n\n", opts.projectName, opts.dstPath)
	return nil
}

// ------------------------------------------------------------------------------------------

func setTemplateDefault(t *genTemplate, dstPath string, projectName string) (string, string, error) {
	templateSrcPath = t.srcPath
	templatePkgName = t.pkgName
	templateFilterFiles = t.filterFiles
	defaultOutputPrefix = t.outputPrefix

	srcPathAbs, dstPathAbs, err := getDir(templateSrcPath, projectName)
	if err != nil {
		return "", "", err
	}
	srcPath := srcPathAbs
	if dstPath == "" {
		dstPath = dstPathAbs
	} else {
		var err error
		dstPath, err = filepath.Abs(dstPath)
		if err != nil {
			return "", "", err
		}
	}

	return srcPath, dstPath, nil
}

// 获取原始目录和输出的目标目录
func getDir(path string, projectName string) (string, string, error) {
	srcPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}
	dstPath := fmt.Sprintf("%s%s%s", gofile.GetRunPath(), getDelimiter(), projectName+"_"+defaultOutputPrefix+time.Now().Format("0102150405"))
	return srcPath, dstPath, nil
}

// 过滤文件
func filterFiles(filePaths []string, files ...string) []string {
	out := []string{}
	for _, file := range filePaths {
		_, filename := filepath.Split(file)
		isExist := false
		for _, v := range files {
			if filename == v {
				isExist = true
				break
			}
		}
		if !isExist {
			out = append(out, file)
		}
	}
	return out
}

// 替换文件路径、名称、内容
func replaceFile(oldFilePath string, newFilePath string, kvs map[string]string) error {
	// 创建目录
	path, newFilename := filepath.Split(newFilePath)
	err := os.MkdirAll(path, 0666)
	if err != nil {
		return err
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(oldFilePath)
	if err != nil {
		return err
	}

	// 替换文件名和文件内容
	for oldStr, newStr := range kvs {
		newFilename = strings.ReplaceAll(newFilename, oldStr, newStr)
		newFilePath = path + getDelimiter() + newFilename
		data = bytes.ReplaceAll(data, []byte(oldStr), []byte(newStr))
	}

	return ioutil.WriteFile(newFilePath, data, 0666)
}

// 替换文件参数，区分大小写
func getReplaceParams(apiName string, projectName string) map[string]string {
	kv := make(map[string]string)
	// 默认旧包名改为项目名称
	kv[templatePkgName] = projectName
	// 默认api名称改为新api名
	oldStr := templateApiName
	newStr := apiName

	// 把第一个字母转为大写
	oldStr = strings.ToUpper(oldStr[:1]) + oldStr[1:]
	newStr = strings.ToUpper(newStr[:1]) + newStr[1:]
	kv[oldStr] = newStr

	// 把第一个字母转为小写
	oldStr = strings.ToLower(oldStr[:1]) + oldStr[1:]
	newStr = strings.ToLower(newStr[:1]) + newStr[1:]
	kv[oldStr] = newStr

	return kv
}

// 根据系统类型获取分隔符
func getDelimiter() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}

	return delimiter
}
