package templates

import (
	"embed"

	"github.com/zhufuyi/goctl/pkg/replacer"
)

// Replacers 各个模板对应的接口
var Replacers = map[string]replacer.Replacer{}

// Template 模板信息
type Template struct {
	Name     string
	FS       embed.FS
	FilePath string
}

// Init 初始化模板
func Init(templates []Template) {
	var err error
	for _, v := range templates {
		if _, ok := Replacers[v.Name]; ok {
			panic(v.Name + " already exists")
		}
		Replacers[v.Name], err = replacer.NewWithFS(v.FilePath, v.FS)
		if err != nil {
			panic(err)
		}
	}
}
