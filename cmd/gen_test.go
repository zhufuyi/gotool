package cmd

import (
	"embed"
	"testing"

	"github.com/zhufuyi/goctl/util/template"
)

// api模板目录
//go:embed gen.go
var apiFS embed.FS

// web模板目录
//go:embed gen.go
var webFS embed.FS

// user模板目录
//go:embed gen.go
var userFS embed.FS

func Test_runGenApiCommand(t *testing.T) {
	handler, err := template.New("../templates/api", apiFS)
	if err != nil {
		t.Fatal(err)
	}
	projectName := "demo"
	apiName := "user"
	outPath := "C:\\Users\\zys\\Desktop\\demo"

	err = runGenApiCommand(handler, projectName, apiName, outPath)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_runGenWebCommand(t *testing.T) {
	handler, err := template.New("../templates/web", webFS)
	if err != nil {
		t.Fatal(err)
	}
	projectName := "demo"
	apiName := "user"
	outPath := "C:\\Users\\zys\\Desktop\\demo"

	err = runGenApiCommand(handler, projectName, apiName, outPath)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_runGenUserCommand(t *testing.T) {
	projectName := "demo"
	handler, err := template.New("../templates/user", userFS)
	if err != nil {
		t.Fatal(err)
	}
	outPath := "C:\\Users\\zys\\Desktop\\demo"

	err = runGenUserCommand(handler, projectName, outPath)
	if err != nil {
		t.Fatal(err)
	}
}
