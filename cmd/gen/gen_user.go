package gen

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/goctl/pkg/replace"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
)

// UserCommand generate user code
func UserCommand() *cobra.Command {
	var (
		projectName string // 项目名称
		outPath     string // 输出目录
	)

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user code",
		Long: `generate user code.

Examples:
  # generate user code, including registration, login and logout api
  goctl gen user -p yourProjectName
  goctl gen user -p yourProjectName -o /tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runGenUserCommand(templates.Handers[GenTypeUser], projectName, outPath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")

	return cmd
}

func runGenUserCommand(handler replace.Handler, projectName string, outPath string) error {
	// 设置模板信息
	templateIgnoreFiles := []string{}     // 忽略处理的文件
	replacementFields := []replace.Field{ // 替换字段
		{
			Old:             "github.com/zhufuyi/goctl/templates/user",
			New:             projectName,
			IsCaseSensitive: false, // 第一个字母区分大小写
		},
	}

	if handler == nil {
		return errors.New("handler is nil")
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
