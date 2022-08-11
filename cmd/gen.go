package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/global"
	"github.com/zhufuyi/goctl/util/template"
)

const (
	genTypeApi  = "api"
	genTypeWeb  = "web"
	genTypeUser = "user"
)

func genGinCommand() *cobra.Command {
	var (
		resourceArg string
		outPath     string // 输出目录
		apiName     string // 接口名称
		projectName string // 项目名称
	)

	cmd := &cobra.Command{
		Use:   "gen <type>",
		Short: "Generate web service code",
		Long: `generate web service code.

Examples:
    # list generate types
    goctl gen list

    # generate web service code
    goctl gen web -p yourProjectName -a yourApiName
    goctl gen web -p yourProjectName -a yourApiName -o /tmp

    # generate api code
    goctl gen api -p yourProjectName -a yourApiName
    goctl gen api -p yourProjectName -a yourApiName -o /tmp

    # generate user code, including registration, login and logout api
    goctl gen user -p yourProjectName
    goctl gen user  -p yourProjectName -o /tmp
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
			case "list":
				fmt.Println(string(ListTypeNames(genTypeWeb, genTypeApi, genTypeUser)))

			case genTypeApi:
				err := runGenApiCommand(global.ApiTemplater, projectName, apiName, outPath)
				if err != nil {
					return err
				}

			case genTypeWeb:
				err := runGenWebCommand(global.WebTemplater, projectName, apiName, outPath)
				if err != nil {
					return err
				}

			case genTypeUser:
				err := runGenUserCommand(global.UserTemplater, projectName, outPath)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown resource name '%s'. Use \"goctl resources\" for a complete list of supported resources.\n", resourceArg)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "projectName", "p", "", "project name")
	//cmd.MarkFlagRequired("projectName")
	cmd.Flags().StringVarP(&apiName, "apiName", "a", "", "api name")
	//cmd.MarkFlagRequired("apiName")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")

	return cmd
}

func checkParams(projectName string, apiName string) error {
	if projectName == "" {
		return errors.New("required flag(s) 'projectName' not set")
	}
	if apiName == "" {
		return errors.New("required flag(s) 'apiName' not set")
	}
	return nil
}

func runGenApiCommand(handler template.Handler, projectName string, apiName string, outPath string) error {
	if err := checkParams(projectName, apiName); err != nil {
		return err
	}

	// 设置模板信息
	templateIgnoreFiles := []string{"dao.go", "common_code.go", "service.go", "routers.go", "global.go", "conf.go"} // 忽略处理的文件
	fields := []template.Field{
		{
			Old:          "UserExample",
			New:          apiName,
			IsNeedCovert: true,
		}, {
			Old:          "github.com/zhufuyi/goctl/templates/api",
			New:          projectName,
			IsNeedCovert: false,
		},
	}

	handler.SetIgnoreFiles(templateIgnoreFiles...)
	handler.SetReplacementFields(fields)
	if err := handler.SetOutPath(outPath, apiName+"_"+genTypeApi); err != nil {
		return err
	}
	if err := handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate api '%s' code successfully, output = %s\n\n", apiName, handler.GetOutPath())
	return nil
}

func runGenWebCommand(handler template.Handler, projectName string, apiName string, outPath string) error {
	if err := checkParams(projectName, apiName); err != nil {
		return err
	}

	// 设置模板信息
	templateIgnoreFiles := []string{} // 忽略处理的文件
	fields := []template.Field{       // 替换字段
		{
			Old:          "UserExample",
			New:          apiName,
			IsNeedCovert: true,
		}, {
			Old:          "github.com/zhufuyi/goctl/templates/web",
			New:          projectName,
			IsNeedCovert: false,
		},
	}

	handler.SetIgnoreFiles(templateIgnoreFiles...)
	handler.SetReplacementFields(fields)
	if err := handler.SetOutPath(outPath, projectName+"_"+genTypeWeb); err != nil {
		return err
	}
	if err := handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate web server '%s' code successfully, output = %s\n\n", projectName, handler.GetOutPath())
	return nil
}

func runGenUserCommand(handler template.Handler, projectName string, outPath string) error {
	// 设置模板信息
	templateIgnoreFiles := []string{}      // 忽略处理的文件
	replacementFields := []template.Field{ // 替换字段
		{
			Old:          "github.com/zhufuyi/goctl/templates/user",
			New:          projectName,
			IsNeedCovert: false, // 不区分大小写
		},
	}

	handler.SetIgnoreFiles(templateIgnoreFiles...)
	handler.SetReplacementFields(replacementFields)
	if err := handler.SetOutPath(outPath, projectName); err != nil {
		return err
	}
	if err := handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate project '%s' code successfully, output = %s\n\n", projectName, handler.GetOutPath())
	return nil
}

// ------------------------------------------------------------------------------------------

// ListTypeNames 类型名称列表
func ListTypeNames(names ...string) []byte {
	content := []string{fmt.Sprintf("%d types are supported:\n", len(names))}
	for _, name := range names {
		content = append(content, "    "+name+"\n")
	}

	return []byte(strings.Join(content, ""))
}
