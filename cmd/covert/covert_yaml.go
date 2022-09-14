package covert

import (
	"fmt"
	"github.com/zhufuyi/pkg/gofile"
	"os"
	"path/filepath"

	"github.com/zhufuyi/goctl/pkg/jy2struct"

	"github.com/spf13/cobra"
)

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
  goctl covert yaml --data="yaml text"

  # covert yaml to struct from file
  goctl covert yaml --file=test.yaml

  # covert yaml to struct, set tag value and subStruct flag
  goctl covert yaml --file=test.yaml --tags=gorm --sub-struct=false

  # covert yaml to struct, set tag value, save to specified directory, file name is config.go
  goctl covert yaml --file=test.yaml --tags=json --out=/tmp
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

	cmd.Flags().StringVarP(&ysArgs.Data, "data", "d", "", "yaml data")
	cmd.Flags().StringVarP(&ysArgs.InputFile, "file", "f", "", "yaml file")
	cmd.Flags().StringVarP(&ysArgs.Tags, "tags", "t", "", "specify tags in addition to the format, with multiple tags separated by commas")
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
	err = os.WriteFile(outFile, []byte(goParseCode+code), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("covert '%s' to go struct successfully, output = %s\n\n", inputFile, outFile)
	return nil
}

const goParseCode = `// nolint
// code generated from config file.

package config

import "github.com/zhufuyi/pkg/conf"

type Config = GenerateName

var config *Config

// Init parsing configuration files to struct, including yaml, toml, json, etc.
func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show() {
	conf.Show(config)
}

func Get() *Config {
	if config == nil {
		panic("config is nil")
	}
	return config
}
`
