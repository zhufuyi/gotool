package main

import (
	"embed"
	"os"

	"github.com/zhufuyi/goctl/cmd"
	"github.com/zhufuyi/goctl/global"
	"github.com/zhufuyi/goctl/utils/template"
)

// api模板目录
//
//go:embed templates/api
var apiFS embed.FS

// web模板目录
//
//go:embed templates/web
var webFS embed.FS

// user模板目录
//
//go:embed templates/user
var userFS embed.FS

func initTemplates() {
	var err error

	global.ApiTemplater, err = template.New("templates/api", apiFS)
	if err != nil {
		panic(err)
	}

	global.WebTemplater, err = template.New("templates/web", webFS)
	if err != nil {
		panic(err)
	}

	global.UserTemplater, err = template.New("templates/user", userFS)
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
