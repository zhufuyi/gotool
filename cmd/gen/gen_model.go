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
  # generate model code.
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate model code and embed 'gorm.Model''.
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate model code and content mysql code.
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --project-name=yourProjectName

  # generate model code and specify the output directory, Note: if the file already exists in the path, it will replace the original file directly.
  goctl gen model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenModelCommand(projectName, ModuleModel, codes, outPath)
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

func runGenModelCommand(projectName string, moduleName string, codes map[string]string, outPath string) error {
	r := templates.Replacers[moduleName]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{"internal/model"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{}              // 指定子目录下忽略处理的目录
	ignoreFiles := []string{}             // 指定子目录下忽略处理的文件
	if projectName == "" {
		ignoreFiles = append(ignoreFiles, "init.go")
	}

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addModelFields(projectName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutDir(outPath, "gen_"+moduleName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", moduleName, r.GetOutPath())
	return nil
}

func addModelFields(projectName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, modelFile)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
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
				Old: selfPackageName + "/" + r.GetBasePath(),
				New: projectName,
			},
			{
				Old: "github.com/zhufuyi/sponge",
				New: projectName,
			},
			{
				Old: projectName + "/pkg",
				New: "github.com/zhufuyi/sponge/pkg",
			},
		}...)
	}

	return fields
}
