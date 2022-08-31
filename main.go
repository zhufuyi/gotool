package main

import (
	"embed"
	"math/rand"
	"os"
	"time"

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

// handler模板目录
//
//go:embed templates/handler
var handlerFS embed.FS

// http服务模板目录
//
//go:embed templates/http_server
var httpFS embed.FS

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
			Name:     gen.GenTypeHandler,
			FS:       handlerFS,
			FilePath: "templates/handler",
		},
		{
			Name:     gen.GenTypeHTTP,
			FS:       httpFS,
			FilePath: "templates/http_server",
		},
	})

	rand.Seed(time.Now().UnixNano())
	
	rootCMD := cmd.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
