package covert

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/jy2struct"

	"github.com/spf13/cobra"
)

// Yaml2StructCommand covert yaml to struct command
func Yaml2StructCommand() *cobra.Command {
	var (
		// yaml to struct args
		ysArgs = jy2struct.Args{}
	)

	cmd := &cobra.Command{
		Use:   "yaml",
		Short: "Covert yaml to struct",
		Long: `covert yaml to struct.

Examples:
  # covert yaml to struct from data
  goctl covert yaml --data="yaml text"

  # covert yaml to struct from file
  goctl covert yaml --file=test.yaml

  # covert yaml to struct, set tag value and subStruct flag
  goctl covert yaml --file=test.sql --tags=gorm --sub-struct=false

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ysArgs.Format = covertTypeYaml2Struct
			out, err := jy2struct.Covert(&ysArgs)
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&ysArgs.Data, "data", "d", "", "yaml data")
	cmd.Flags().StringVarP(&ysArgs.InputFile, "file", "f", "", "yaml file")
	cmd.Flags().StringVarP(&ysArgs.Tags, "tags", "t", "", "specify tags in addition to the format, with multiple tags separated by commas")
	cmd.Flags().BoolVarP(&ysArgs.SubStruct, "sub-struct", "s", true, "create types for sub-structs (default is true)")

	return cmd
}
