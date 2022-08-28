package gen

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/goctl/pkg/replace"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
)

// WebCommand generate web code
func WebCommand() *cobra.Command {
	var (
		apiName     string // 接口名称
		projectName string // 项目名称
		outPath     string // 输出目录
	)

	cmd := &cobra.Command{
		Use:   "web",
		Short: "Generate web code",
		Long: `generate web code.

Examples:
  # generate web code
  goctl gen web -p yourProjectName -a yourApiName
  goctl gen web -p yourProjectName -a yourApiName -o /tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runGenWebCommand(templates.Handers[GenTypeWeb], projectName, apiName, outPath)
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

func runGenWebCommand(handler replace.Handler, projectName string, apiName string, outPath string) error {
	// 设置模板信息
	templateIgnoreFiles := []string{} // 忽略处理的文件
	fields := []replace.Field{        // 替换字段
		{
			Old:             "UserExample",
			New:             apiName,
			IsCaseSensitive: true, // 第一个字母不区分大小写，UserExample和userExample都满足匹配条件
		}, {
			Old:             "github.com/zhufuyi/goctl/templates/web",
			New:             projectName,
			IsCaseSensitive: false, // 第一个字母区分大小写
		},
	}

	if handler == nil {
		return errors.New("handler is nil")
	}
	handler.SetIgnoreFiles(templateIgnoreFiles...)
	handler.SetReplacementFields(fields)
	if err := handler.SetOutPath(outPath, projectName+"_"+GenTypeWeb); err != nil {
		return err
	}
	if err := handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate web server '%s' code successfully, output = %s\n\n", projectName, handler.GetOutPath())
	return nil
}
