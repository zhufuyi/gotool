package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/utils/sql2gorm"
)

const (
	covertTypeSql2gorm = "sql"
	//covertTypeJSON2Struct = "json"
)

func convertCommand() *cobra.Command {
	var (
		resourceArg string // 转换类型

		// sql to gorm 参数
		sqlArgs = sql2gorm.Args{}
	)

	cmd := &cobra.Command{
		Use:   "covert <type>",
		Short: "Covert sql to gorm",
		Long: `covert sql to gorm

Examples:
    # list covert types
    goctl covert list

    # covert sql to gorm from sql
    goctl covert sql --sql="sql text"

    # covert sql to gorm from file
    goctl covert sql --file=test.sql

    # covert sql to gorm from db
    goctl covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

    # covert sql to gorm, set package name and json tag
    goctl covert sql --file=test.sql --pkg-name=user --json-tag
    goctl covert sql --file=test.sql --pkg-name=user --json-tag --json-named-type=1
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must specify the type of resource to gen. use 'goctl resources' for a complete list of supported resources")
			}
			resourceArg = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			switch resourceArg {
			case "list":
				fmt.Println(string(ListTypeNames(covertTypeSql2gorm)))

			case covertTypeSql2gorm:
				out, err := sql2gorm.GetGormCode(&sqlArgs)
				if err != nil {
					return err
				}
				fmt.Println(out)

			default:
				return fmt.Errorf("unknown resource name '%s'. Use \"goctl resources\" for a complete list of supported resources.\n", resourceArg)
			}

			return nil
		},
	}

	// sql to gorm 参数
	cmd.Flags().StringVarP(&sqlArgs.Sql, "sql", "s", "", "sql")
	cmd.Flags().StringVarP(&sqlArgs.InputFile, "file", "f", "", "input sql file")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, E.g. user:password@(host:port)/database")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	cmd.Flags().StringVarP(&sqlArgs.Package, "pkg-name", "p", "", "package name")
	cmd.Flags().BoolVarP(&sqlArgs.JsonTag, "json-tag", "j", false, "whether to generate json tag")
	cmd.Flags().IntVarP(&sqlArgs.JsonNamedType, "json-named-type", "J", 0, "json named type, 0:snake_case, other:camelCase")

	return cmd
}
