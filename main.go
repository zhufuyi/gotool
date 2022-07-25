package main

import (
	"embed"
	"os"

	"github.com/zhufuyi/goctl/cmd"
	"github.com/zhufuyi/goctl/global"
	"github.com/zhufuyi/goctl/util/template"
)

// api模板目录
//go:embed templates/api
var apiFS embed.FS

// web模板目录
//go:embed templates/web
var webFS embed.FS

func initTemplates() {
	var err error

	global.ApiTemplate, err = template.New("templates/api", apiFS)
	if err != nil {
		panic(err)
	}

	global.WebTemplate, err = template.New("templates/web", webFS)
	if err != nil {
		panic(err)
	}
}

func main() {
	initTemplates()

	rootCMD := cmd.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
