package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/util/template"
)

func replaceCommand() *cobra.Command {
	var (
		srcPath  string   // 源目录
		oldValue []string // 旧字段
		newValue []string // 新字段
	)

	cmd := &cobra.Command{
		Use:   "replace <path> <old...> <new...>",
		Short: "Replace fields in path files",
		Long: `replace fields in path files.

Examples:
    goctl replace -p /tmp -o oldField -n newField
    goctl replace -p /tmp -o oldField1 -n newField1 -o oldField2 -n newField2

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runReplaceCommand(srcPath, oldValue, newValue)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&srcPath, "path", "p", "", "source path or file")
	cmd.MarkFlagRequired("path")

	cmd.Flags().StringArrayVarP(&oldValue, "old", "o", nil, "old value, one by one corresponding to the 'new' fields")
	cmd.MarkFlagRequired("old")
	cmd.Flags().StringArrayVarP(&newValue, "new", "n", nil, "new value, one by one corresponding to the 'old' fields")
	cmd.MarkFlagRequired("new")

	return cmd
}

func runReplaceCommand(srcPath string, oldValues []string, newValues []string) error {
	if len(oldValues) != len(newValues) {
		return errors.New("len(old) must be equal to len(new)")
	}

	handler, err := template.NewSrc(srcPath)
	if err != nil {
		return err
	}

	// 设置模板信息
	templateIgnoreFiles := []string{} // 忽略处理的文件
	fields := []template.Field{}
	for i, old := range oldValues {
		fields = append(fields, template.Field{
			Old:          old,
			New:          newValues[i],
			IsNeedCovert: false,
		})
	}

	handler.SetIgnoreFiles(templateIgnoreFiles...)
	handler.SetReplacementFields(fields)
	if err := handler.SetOutPath("", "replace"); err != nil {
		return err
	}
	if err := handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("replace successfully, output = %s\n\n", handler.GetOutPath())
	return nil
}
