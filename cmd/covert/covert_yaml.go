package covert

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhufuyi/goctl/pkg/jy2struct"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/pkg/gofile"
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
  # covert yaml to struct from file
  goctl covert yaml --file=test.yaml

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

	cmd.Flags().StringVarP(&ysArgs.InputFile, "file", "f", "", "yaml file")
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

func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show() string {
	return conf.Show(config)
}

func Get() *Config {
	if config == nil {
		panic("config is nil")
	}
	return config
}`
