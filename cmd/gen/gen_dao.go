package gen

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/goctl/pkg/sql2code/parser"

	"github.com/zhufuyi/goctl/pkg/replace"
	"github.com/zhufuyi/goctl/pkg/sql2code"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
)

// DaoCommand generate dao code
func DaoCommand() *cobra.Command {
	var (
		projectName string // 项目名称
		outPath     string // 输出目录
		sqlArgs     = sql2code.Args{
			Package:       "model",
			JsonTag:       true,
			JsonNamedType: 0,
		}
	)

	cmd := &cobra.Command{
		Use:   "dao",
		Short: "Generate dao code",
		Long: `generate dao code.

Examples:
  # generate dao code
  goctl gen dao --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate dao code and embed 'gorm.model'
  goctl gen dao --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate dao code and specify the output path
  goctl gen dao --project-name=yourProjectName  --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {

			codes, err := sql2code.GetCodes(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenDaoCommand(GenTypeDao, projectName, codes, outPath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, E.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embedded", "e", false, "whether to embed gorm.Model")

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output path")

	return cmd
}

func runGenDaoCommand(genType string, projectName string, codes map[string]string, outPath string) error {
	// 设置模板信息
	ignoreFiles := []string{"conf.go", "conf.yml", "conf_test.go", "init.go"} // 忽略处理的文件
	fields := []replace.Field{
		{
			Old: "github.com/zhufuyi/goctl/templates/dao",
			New: projectName,
		},
		{
			Old: "// todo generate update fields code",
			New: codes[parser.CodeTypeDAO],
		},
	}

	handler := templates.Handers[genType]
	if handler == nil {
		return errors.New("handler is nil")
	}
	content, err := handler.ReadFile("model/userExample.go")
	if err != nil {
		return err
	}

	fields = append(fields,
		replace.Field{ // 替换文件内容
			Old: string(content),
			New: codes[parser.CodeTypeModel],
		},
		replace.Field{ // 在替换文件内容之后，才能使得修改文件名生效
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	)

	handler.SetIgnoreFiles(ignoreFiles...)
	handler.SetReplacementFields(fields)
	if err = handler.SetOutPath(outPath, "gen_"+genType); err != nil {
		return err
	}
	if err = handler.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", genType, handler.GetOutPath())
	return nil
}
