package gen

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/global"
	"github.com/zhufuyi/goctl/utils/template"
)

// ApiCommand generate api code
func ApiCommand() *cobra.Command {
	var (
		apiName     string // 接口名称
		projectName string // 项目名称
		outPath     string // 输出目录
	)

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Generate api code",
		Long: `generate api code.

Examples:
  # generate api code
  goctl gen api -p yourProjectName -a yourApiName
  goctl gen api -p yourProjectName -a yourApiName -o /tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := runGenApiCommand(global.ApiTemplater, projectName, apiName, outPath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&apiName, "api-name", "a", "", "api name")
	cmd.MarkFlagRequired("api-name")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")

	return cmd
}

func runGenApiCommand(handler template.Handler, projectName string, apiName string, outPath string) error {
	// 设置模板信息
	templateIgnoreFiles := []string{"dao.go", "common_code.go", "service.go", "routers.go", "global.go", "conf.go"} // 忽略处理的文件
	fields := []template.Field{
		{
			Old:             "UserExample",
			New:             apiName,
			IsCaseSensitive: true, // 第一个字母不区分大小写，UserExample和userExample都满足匹配条件
		}, {
			Old:             "github.com/zhufuyi/goctl/templates/api",
			New:             projectName,
			IsCaseSensitive: false, // 第一个字母区分大小写
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
