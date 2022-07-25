package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/global"
	"github.com/zhufuyi/goctl/util/template"
)

func genGinCommand() *cobra.Command {
	var (
		resourceArg string
		outPath     string // 输出目录
		apiName     string // 接口名称
		projectName string // 项目名称
	)

	cmd := &cobra.Command{
		Use:   "gen <resource>",
		Short: "Generate gin api code",
		Long: `generate gin api code.

Examples:
    # generate api code
    goctl gen web -p yourProjectName -a yourApiName
    goctl gen web -p yourProjectName -a yourApiName -o /tmp

    # generate web server code
    goctl gen api -p yourProjectName -a yourApiName
    goctl gen api -p yourProjectName -a yourApiName -o /tmp
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must specify the type of resource to gen. use 'goctl resources' for a complete list of supported resources")
			}
			resourceArg = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			switch resourceArg {
			case api:
				err := runGenApiCommand(apiName, projectName, outPath)
				if err != nil {
					return err
				}

			case web:
				err := runGenWebCommand(apiName, projectName, outPath)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown resource name '%s'. Use \"goctl resources\" for a complete list of supported resources.\n", resourceArg)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&apiName, "apiName", "a", "", "api name")
	cmd.MarkFlagRequired("apiName")
	cmd.Flags().StringVarP(&projectName, "projectName", "p", "", "project name")
	cmd.MarkFlagRequired("projectName")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")

	return cmd
}

func runGenApiCommand(apiName string, projectName string, outPath string) error {
	templateIgnoreFiles := []string{"dao.go", "common_code.go", "service.go", "routers.go", "global.go"}
	templateApiName := "UserExample"
	templatePkgName := "github.com/zhufuyi/goctl/templates/api"

	fields := []template.Field{
		{
			Old:          templateApiName,
			New:          apiName,
			IsNeedCovert: true,
		}, {
			Old:          templatePkgName,
			New:          projectName,
			IsNeedCovert: false,
		},
	}
	global.ApiTemplate.SetIgnoreFiles(templateIgnoreFiles...)
	global.ApiTemplate.SetReplacementFields(fields)
	err := global.ApiTemplate.SetOutPath(outPath, apiName+"_"+api)
	if err != nil {
		return err
	}
	err = global.ApiTemplate.SaveFiles()
	if err != nil {
		return err
	}

	fmt.Printf("api '%s' generate successfully, output = %s\n\n", apiName, global.ApiTemplate.GetOutPath())

	return nil
}

func runGenWebCommand(apiName string, projectName string, outPath string) error {
	templateApiName := "UserExample"
	templatePkgName := "github.com/zhufuyi/goctl/templates/web"

	fields := []template.Field{
		{
			Old:          templateApiName,
			New:          apiName,
			IsNeedCovert: true,
		}, {
			Old:          templatePkgName,
			New:          projectName,
			IsNeedCovert: false,
		},
	}
	global.WebTemplate.SetReplacementFields(fields)
	err := global.WebTemplate.SetOutPath(outPath, projectName+"_"+web)
	if err != nil {
		return err
	}
	err = global.WebTemplate.SaveFiles()
	if err != nil {
		return err
	}

	fmt.Printf("web server '%s' generate successfully, output = %s\n\n", projectName, global.WebTemplate.GetOutPath())
	return nil
}
