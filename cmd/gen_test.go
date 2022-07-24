package cmd

import (
	"testing"
)

func Test_runGenApiCommand(t *testing.T) {
	//dstPath := "C:\\Users\\zys\\Desktop\\demo"
	dstPath := ""
	projectName := "demo"
	apiName := "user"
	apiTemplate.srcPath = "../templates/api"
	srcPathAbs, dstPathAbs, err := setTemplateDefault(apiTemplate, dstPath, projectName)
	if err != nil {
		t.Fatal(err)
	}

	err = runGenApiCommand(&genApiOptions{
		srcPath:     srcPathAbs,
		dstPath:     dstPathAbs,
		apiName:     apiName,
		projectName: projectName,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_runGenWebCommand(t *testing.T) {
	//dstPath := "C:\\Users\\zys\\Desktop\\demo"
	dstPath := ""
	projectName := "demo"
	apiName := "user"
	webTemplate.srcPath = "../templates/web"
	srcPathAbs, dstPathAbs, err := setTemplateDefault(webTemplate, dstPath, projectName)
	if err != nil {
		t.Fatal(err)
	}

	err = runGenWebCommand(&genWebOptions{
		srcPath:     srcPathAbs,
		dstPath:     dstPathAbs,
		apiName:     apiName,
		projectName: projectName,
	})
	if err != nil {
		t.Fatal(err)
	}
}
