package gen

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/goctl/pkg/replacer"
	"github.com/zhufuyi/goctl/pkg/sql2code"
	"github.com/zhufuyi/goctl/pkg/sql2code/parser"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
)

// ModelCommand generate model code
func ModelCommand() *cobra.Command {
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
		Use:   "model",
		Short: "Generate model code",
		Long: `generate model code.

Examples:
  # generate model code
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate model code and content mysql code
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --project-name=yourProjectName

  # generate model code and embed 'gorm.Model''
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate model code and specify the output directory
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenModelCommand(projectName, GenTypeModel, codes, outPath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, E.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embedded", "e", false, "whether to embed gorm.Model")

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name, whether to generate the mysql contention code")

	return cmd
}

func runGenModelCommand(projectName string, genType string, codes map[string]string, outPath string) error {
	r := templates.Replacers[genType]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	ignoreFiles := []string{"conf.go", "conf.yml", "conf_test.go"} // 忽略处理的文件
	if projectName == "" {
		ignoreFiles = append(ignoreFiles, "init.go")
	}
	fields := addModelFields(projectName, r, codes)

	r.SetIgnoreFiles(ignoreFiles...)
	r.SetReplacementFields(fields)
	_ = r.SetOutPath(outPath, "gen_"+genType)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", genType, r.GetOutPath())
	return nil
}

func addModelFields(projectName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, "model/userExample.go")...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: "// todo generate model codes to here",
			New: codes[parser.CodeTypeModel],
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	if projectName != "" {
		fields = append(fields, []replacer.Field{
			{
				Old: "github.com/zhufuyi/goctl/templates/model",
				New: projectName,
			},
			{
				Old: "projectExample",
				New: projectName,
			},
		}...)
	}

	return fields
}
