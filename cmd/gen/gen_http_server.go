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

// HTTPCommand generate http code
func HTTPCommand() *cobra.Command {
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
		Use:   "http",
		Short: "Generate http code",
		Long: `generate http code.

Examples:
  # generate http code.
  goctl gen http --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate http code and embed 'gorm.model'.
  goctl gen http --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate http code and specify the output path, Note: if the file already exists in the path, it will replace the original file directly.
  goctl gen http --project-name=yourProjectName  --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenHTTPCommand(projectName, ModuleHTTP, sqlArgs.DBDsn, codes, outPath)
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

func runGenHTTPCommand(projectName string, moduleName string, dbDSN string, codes map[string]string, outPath string) error {
	r := templates.Replacers[moduleName]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{} // 只处理的子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{"sponge/api", "sponge/pkg",
		"sponge/third_party", "internal/service", "sponge/test"} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"swagger.json", "swagger.yaml", "protoc.sh",
		"proto-doc.sh", "grpc.go", "grpc_option.go", "LICENSE", "doc.go",
		"grpc_userExample.go", "grpc_systemCode.go", "proto.html"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addHTTPFields(projectName, r, dbDSN, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutDir(outPath, "gen_"+moduleName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", moduleName, r.GetOutPath())
	return nil
}

func addHTTPFields(projectName string, r replacer.Replacer, dbDSN string, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, modelFile)...)
	fields = append(fields, addTheDeleteFields(r, daoFile)...)
	fields = append(fields, addTheDeleteFields(r, handlerFile)...)
	fields = append(fields, addTheDeleteFields(r, mainFile)...)
	fields = append(fields, addTheDeleteGrpcFields(r, mainFile)...)
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
		{ // 替换main.go文件内容
			Old: mainFileMark,
			New: httpServerRegisterCode,
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
			Old: projectName + "/pkg",
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old: "sponge api docs",
			New: projectName + " api docs",
		},
		{
			Old: `"sponge"`,
			New: "\"" + projectName + "\"",
		},
		{
			Old: "userExampleNO = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(1000)),
		},
		{
			Old: "name: \"userExample\"",
			New: "name: " + "\"" + projectName + "\"",
		},
		{
			Old: "go.mod.bak",
			New: "go.mod",
		},
		{
			Old: "go.sum.bak",
			New: "go.sum",
		},
		{
			Old: "root:123456@(192.168.3.37:3306)/account",
			New: dbDSN,
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
