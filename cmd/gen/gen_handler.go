package gen

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/zhufuyi/goctl/pkg/replacer"
	"github.com/zhufuyi/goctl/pkg/sql2code"
	"github.com/zhufuyi/goctl/pkg/sql2code/parser"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
)

// HandlerCommand generate handler code
func HandlerCommand() *cobra.Command {
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
		Use:   "handler",
		Short: "Generate handler code",
		Long: `generate handler code.

Examples:
  # generate handler code.
  goctl gen handler --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate handler code and embed 'gorm.model'.
  goctl gen handler --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate handler code and specify the output path, Note: if the file already exists in the path, it will replace the original file directly.
  goctl gen handler --project-name=yourProjectName  --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenHandlerCommand(projectName, ModuleHandler, codes, outPath)
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

func runGenHandlerCommand(projectName string, moduleName string, codes map[string]string, outPath string) error {
	r := templates.Replacers[moduleName]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{"internal/model", "internal/cache", "internal/dao",
		"internal/ecode", "internal/handler", "internal/routers"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "swagger_types.go", "http_systemCode.go",
		"grpc_systemCode.go", "grpc_userExample.go", "routers.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addHandlerFields(projectName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutDir(outPath, "gen_"+moduleName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", moduleName, r.GetOutPath())
	return nil
}

func addHandlerFields(projectName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, modelFile)...)
	fields = append(fields, addTheDeleteFields(r, daoFile)...)
	fields = append(fields, addTheDeleteFields(r, handlerFile)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // 替换handler/userExample.go文件内容
			Old: handlerFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeHandler]),
		},
		{
			Old: selfPackageName + "/" + r.GetBasePath(),
			New: projectName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: projectName,
		},
		{
			Old: "userExampleNO = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(1000)),
		},
		{
			Old: projectName + "/pkg",
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
