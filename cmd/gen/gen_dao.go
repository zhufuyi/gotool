package gen

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/goctl/pkg/sql2code/parser"

	"github.com/zhufuyi/goctl/pkg/replacer"
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
			Package:  "model",
			JSONTag:  true,
			GormType: true,
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
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenDaoCommand(projectName, GenTypeDao, codes, outPath)
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

func runGenDaoCommand(projectName string, genType string, codes map[string]string, outPath string) error {
	r := templates.Replacers[genType]
	if r == nil {
		return errors.New("r is nil")
	}

	// 设置模板信息
	ignoreFiles := []string{"conf.go", "conf.yml", "conf_test.go", "init.go"} // 忽略处理的文件
	fields := addDAOFields(projectName, r, codes)

	r.SetIgnoreFiles(ignoreFiles...)
	r.SetReplacementFields(fields)
	_ = r.SetOutPath(outPath, "gen_"+genType)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", genType, r.GetOutPath())
	return nil
}

func addDAOFields(projectName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, "model/userExample.go")...)
	fields = append(fields, addTheDeleteFields(r, "dao/userExample.go")...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: "// todo generate model codes to here",
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: "// todo generate the update fields code to here",
			New: codes[parser.CodeTypeDAO],
		},
		{
			Old: "github.com/zhufuyi/goctl/templates/dao",
			New: projectName,
		},
		{
			Old: "projectExample",
			New: projectName,
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
