package templates

import (
	"embed"

	"github.com/zhufuyi/goctl/pkg/replace"
)

// Handers 名称对应模板处理接口
var Handers = map[string]replace.Handler{}

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
		if _, ok := Handers[v.Name]; ok {
			panic(v.Name + " already exists")
		}
		Handers[v.Name], err = replace.New(v.FilePath, v.FS)
		if err != nil {
			panic(err)
		}
	}
}
