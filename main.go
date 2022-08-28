package main

import (
	"embed"
	"os"

	"github.com/zhufuyi/goctl/cmd/gen"

	"github.com/zhufuyi/goctl/cmd"
	"github.com/zhufuyi/goctl/templates"
)

// model模板目录
//
//go:embed templates/model
var modelFS embed.FS

// dao模板目录
//
//go:embed templates/dao
var daoFS embed.FS

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

func main() {
	// 初始化模板文件
	templates.Init([]templates.Template{
		{
			Name:     gen.GenTypeModel,
			FS:       modelFS,
			FilePath: "templates/model",
		},
		{
			Name:     gen.GenTypeDao,
			FS:       daoFS,
			FilePath: "templates/dao",
		},
		{
			Name:     gen.GenTypeApi,
			FS:       apiFS,
			FilePath: "templates/api",
		},
		{
			Name:     gen.GenTypeWeb,
			FS:       webFS,
			FilePath: "templates/web",
		},
		{
			Name:     gen.GenTypeUser,
			FS:       userFS,
			FilePath: "templates/user",
		},
	})

	rootCMD := cmd.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
