package covert

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/utils/toStruct"
)

// JSON2StructCommand covert json to struct command
func JSON2StructCommand() *cobra.Command {
	var (
		// json to struct args
		toStructArgs = toStruct.Args{}
	)

	cmd := &cobra.Command{
		Use:   "json",
		Short: "Covert json to struct",
		Long: `covert json to struct.

Examples:
  # covert json to struct from data
  goctl covert json --data="json text"

  # covert json to struct from file
  goctl covert json --file=test.json

  # covert json to struct, set tag value and subStruct flag
  goctl covert json --file=test.sql --tags=gorm --sub-struct=false

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			toStructArgs.Format = covertTypeJSON2Struct
			out, err := toStruct.GetStructCode(&toStructArgs)
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&toStructArgs.Data, "data", "d", "", "json data")
	cmd.Flags().StringVarP(&toStructArgs.InputFile, "file", "f", "", "json file")
	cmd.Flags().StringVarP(&toStructArgs.Tags, "tags", "t", "", "specify tags in addition to the format, with multiple tags separated by commas")
	cmd.Flags().BoolVarP(&toStructArgs.SubStruct, "sub-struct", "s", true, "create types for sub-structs (default is true)")

	return cmd
}
