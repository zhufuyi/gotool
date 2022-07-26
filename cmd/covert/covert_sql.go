package covert

import (
	"fmt"

	"github.com/zhufuyi/gotool/pkg/sql2code"

	"github.com/spf13/cobra"
)

// SQL2GormCommand sql to gorm
func SQL2GormCommand() *cobra.Command {
	var (
		// sql to gorm args
		sqlArgs = sql2code.Args{}
	)

	cmd := &cobra.Command{
		Use:   "sql",
		Short: "Covert sql to gorm",
		Long: `covert sql to gorm.

Examples:
  # covert sql to gorm model code
  gotool covert sql --sql="sql text"

  # covert sql file to gorm model code
  gotool covert sql --file=test.sql

  # covert mysql table gorm model code
  gotool covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # covert mysql table, structure fields correspond to the column names of the table.
  gotool covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # covert mysql table to handler request and respond struct code,  other type json or dao
  gotool covert sql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --code-type=handler

  # covert sql file to gorm model code and add json tag
  gotool covert sql --file=test.sql --pkg-name=user --json-tag
  gotool covert sql --file=test.sql --pkg-name=user --json-tag --json-named-type=1
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := sql2code.GenerateOne(&sqlArgs)
			if err != nil {
				return err
			}
			fmt.Println(out)

			return nil
		},
	}

	// sql to gorm 参数
	cmd.Flags().StringVarP(&sqlArgs.SQL, "sql", "s", "", "sql data")
	cmd.Flags().StringVarP(&sqlArgs.DDLFile, "file", "f", "", "input DDL sql file")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, E.g. user:password@(host:port)/database")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	cmd.Flags().StringVarP(&sqlArgs.Package, "pkg-name", "p", "", "package name")
	cmd.Flags().StringVarP(&sqlArgs.CodeType, "code-type", "c", "model", "specify the use of the generated code, support 4 types, model(default), json, dao, handler")
	cmd.Flags().BoolVarP(&sqlArgs.JSONTag, "json-tag", "j", false, "whether to generate json tag")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-named-type", "J", 0, "json named type, 0:snake_case, other:camelCase")

	return cmd
}
