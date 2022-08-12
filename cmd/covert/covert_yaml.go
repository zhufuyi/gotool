package covert

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/utils/toStruct"
)

// Yaml2StructCommand covert yaml to struct command
func Yaml2StructCommand() *cobra.Command {
	var (
		// yaml to struct args
		toStructArgs = toStruct.Args{}
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
			toStructArgs.Format = covertTypeYaml2Struct
			out, err := toStruct.GetStructCode(&toStructArgs)
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&toStructArgs.Data, "data", "d", "", "yaml data")
	cmd.Flags().StringVarP(&toStructArgs.InputFile, "file", "f", "", "yaml file")
	cmd.Flags().StringVarP(&toStructArgs.Tags, "tags", "t", "", "specify tags in addition to the format, with multiple tags separated by commas")
	cmd.Flags().BoolVarP(&toStructArgs.SubStruct, "sub-struct", "s", true, "create types for sub-structs (default is true)")

	return cmd
}
