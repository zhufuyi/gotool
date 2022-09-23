package gen

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/zhufuyi/goctl/pkg/replacer"
	"github.com/zhufuyi/goctl/pkg/sql2code"
	"github.com/zhufuyi/goctl/pkg/sql2code/parser"
	"github.com/zhufuyi/goctl/templates"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/pkg/gofile"
)

// ServiceCommand generate service code
func ServiceCommand() *cobra.Command {
	var (
		projectName string // 项目名称，项目下面有多个子服务
		serverName  string // 服务名称
		outPath     string // 输出目录
		sqlArgs     = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Generate grpc service code",
		Long: `generate grpc service code.

Examples:
  # generate service code.
  goctl gen service --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate service code and embed 'gorm.model'.
  goctl gen service --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embedded

  # generate service code and specify the output path, Note: if the file already exists in the path, it will replace the original file directly.
  goctl gen service --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

  # generate protobuf code and specify the service name, Note: If there are multiple services in one project, you need to specify the service name, 
  # the default is that the server name equals the project name.
  goctl gen proto --project-name=yourProjectName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenServiceCommand(projectName, serverName, ModuleService, codes, outPath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")

	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, E.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embedded", "e", false, "whether to embed gorm.Model")

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output path")

	return cmd
}

func runGenServiceCommand(projectName string, serverName string, moduleName string, codes map[string]string, outPath string) error {
	r := templates.Replacers[moduleName]
	if r == nil {
		return errors.New("replacer is nil")
	}

	if serverName == "" {
		serverName = projectName
	}

	// 设置模板信息
	subDirs := []string{"internal/model", "internal/cache", "internal/dao",
		"internal/ecode", "internal/service", "api/userExample"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "http_systemCode.go", "http_userExample.go",
		"grpc_systemCode.go", "service.go", "userExample.pb.go", "userExample.pb.validate.go", "userExample_grpc.pb.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addServiceFields(projectName, serverName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutDir(outPath, "gen_"+moduleName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", moduleName, r.GetOutPath())
	return nil
}

func addServiceFields(projectName string, serverName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, modelFile)...)
	fields = append(fields, addTheDeleteFields(r, daoFile)...)
	fields = append(fields, addTheDeleteFields(r, protoFile)...)
	fields = append(fields, addTheDeleteFields(r, serviceFile)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // 替换v1/userExample.proto文件内容
			Old: protoFileMark,
			New: codes[parser.CodeTypeProto],
		},
		{ // 替换service/userExample_test.go文件内容
			Old: serviceFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeService]),
		},
		{
			Old: selfPackageName + "/" + r.GetBasePath(),
			New: projectName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: projectName,
		},
		// 替换目录名称
		{
			Old: strings.Join([]string{"api", "userExample", "v1"}, gofile.GetPathDelimiter()),
			New: strings.Join([]string{"api", serverName, "v1"}, gofile.GetPathDelimiter()),
		},
		{
			Old: "api/userExample/v1",
			New: fmt.Sprintf("api/%s/v1", serverName),
		},
		{
			Old: "api.userExample.v1",
			New: fmt.Sprintf("api.%s.v1", strings.ReplaceAll(serverName, "-", "")), // proto package 不能存在"-"号
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
