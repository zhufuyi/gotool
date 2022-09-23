package templates

import (
	"embed"

	"github.com/zhufuyi/goctl/pkg/replacer"
)

// Replacers 模板名称对应的接口
var Replacers = map[string]replacer.Replacer{}

// Template 模板信息
type Template struct {
	Name     string
	FS       embed.FS
	FilePath string
}

// Init 初始化模板
func Init(template *Template, groupModuleNames []string) {
	setReplacers(template)

	// 各个子模块共有接口
	for _, name := range groupModuleNames {
		if template.Name == name {
			continue
		}
		Replacers[name] = Replacers[template.Name]
	}
}

func setReplacers(template *Template) {
	var err error
	if _, ok := Replacers[template.Name]; ok {
		panic(template.Name + " already exists")
	}
	Replacers[template.Name], err = replacer.NewWithFS(template.FilePath, template.FS)
	if err != nil {
		panic(err)
	}
}
