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

// ProtoCommand generate protobuf code
func ProtoCommand() *cobra.Command {
	var (
		projectName string // 项目名称，项目下面有多个子服务
		serverName  string // 服务名称
		outPath     string // 输出目录
		sqlArgs     = sql2code.Args{}
	)

	cmd := &cobra.Command{
		Use:   "proto",
		Short: "Generate protobuf code",
		Long: `generate protobuf code.

Examples:
  # generate protobuf code.
  goctl gen proto --project-name=yourProjectName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate protobuf code and specify the output directory, Note: if the file already exists in the path, it will replace the original file directly.
  goctl gen proto --project-name=yourProjectName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			err = runGenProtoCommand(projectName, serverName, ModuleProto, codes, outPath)
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

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")

	return cmd
}

func runGenProtoCommand(projectName string, serverName string, moduleName string, codes map[string]string, outPath string) error {
	r := templates.Replacers[moduleName]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{"api/userExample"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{}               // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"userExample.pb.go", "userExample.pb.validate.go",
		"userExample_grpc.pb.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addProtoFields(projectName, serverName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutDir(outPath, "gen_"+moduleName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, output = %s\n\n", moduleName, r.GetOutPath())
	return nil
}

func addProtoFields(projectName string, serverName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, addTheDeleteFields(r, protoFile)...)
	fields = append(fields, []replacer.Field{
		{ // 替换v1/userExample.proto文件内容
			Old: protoFileMark,
			New: codes[parser.CodeTypeProto],
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: projectName,
		},
		{
			Old: "api/userExample/v1",
			New: fmt.Sprintf("api/%s/v1", serverName),
		},
		{
			Old: "api.userExample.v1",
			New: fmt.Sprintf("api.%s.v1", serverName),
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
