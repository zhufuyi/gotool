package covert

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhufuyi/gotool/pkg/gofile"
	"github.com/zhufuyi/gotool/pkg/jy2struct"

	"github.com/spf13/cobra"
)

const covertTypeYaml2Struct = "yaml"

// Yaml2StructCommand covert yaml to struct command
func Yaml2StructCommand() *cobra.Command {
	var (
		// yaml to struct args
		ysArgs  = jy2struct.Args{}
		outPath = ""
	)

	cmd := &cobra.Command{
		Use:   "yaml",
		Short: "Covert yaml to struct",
		Long: `covert yaml to struct.

Examples:
  # covert yaml to struct from data
  gotool covert yaml --data="yaml text"

  # covert yaml to struct from file
  gotool covert yaml --file=test.yaml

  # covert yaml to struct, set tag value, save to specified directory, file name is config.go
  gotool covert yaml --file=test.yaml --tags=json --out=/tmp

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ysArgs.Format = covertTypeYaml2Struct
			out, err := jy2struct.Covert(&ysArgs)
			if err != nil {
				return err
			}

			if outPath != "" {
				return saveFile(ysArgs.InputFile, outPath, out)
			}

			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&ysArgs.Data, "data", "d", "", "yaml content")
	cmd.Flags().StringVarP(&ysArgs.InputFile, "file", "f", "", "yaml file")
	cmd.Flags().StringVarP(&ysArgs.Tags, "tags", "t", "", "struct tags, multiple tags separated by commas")
	cmd.Flags().BoolVarP(&ysArgs.SubStruct, "sub-struct", "s", true, "create types for sub-structs (default is true)")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "export the code path")
	return cmd
}

func saveFile(inputFile string, outPath string, code string) error {
	abs, err := filepath.Abs(outPath)
	if err != nil {
		return err
	}
	outFile := abs + gofile.GetPathDelimiter() + "config.go"
	err = os.WriteFile(outFile, []byte(code), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("covert '%s' to go struct successfully, output = %s\n\n", inputFile, outFile)
	return nil
}
