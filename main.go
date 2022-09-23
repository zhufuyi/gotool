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

// 微服务模板目录
//
//go:embed templates/sponge
var microServiceFS embed.FS

// http服务模板目录
//
////go:embed templates/http_server
//var httpFS embed.FS

func main() {
	// 初始化模板文件
	templates.Init(&templates.Template{
		Name:     gen.ModuleSponge,
		FS:       microServiceFS,
		FilePath: "templates/sponge",
	}, gen.MicroServiceGroupModules)

	//templates.Init(&templates.Template{
	//	Name:     gen.ModuleHTTP,
	//	FS:       httpFS,
	//	FilePath: "templates/http_server",
	//}, gen.MicroServiceGroupModules)

	rand.Seed(time.Now().UnixNano())

	rootCMD := cmd.NewRootCMD()
	if err := rootCMD.Execute(); err != nil {
		rootCMD.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
