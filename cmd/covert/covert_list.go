package covert

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	covertTypeSQL2gorm    = "sql"
	covertTypeJSON2Struct = "json"
	covertTypeYaml2Struct = "yaml"
)

var covertTypes = []string{covertTypeSQL2gorm, covertTypeJSON2Struct, covertTypeYaml2Struct}

// ListTypesCommand show covert support types
func ListTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List covert types",
		Long: `list covert types.

Examples:
  # show covert types
  goctl covert list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(utils.ListTypeNames(covertTypes...))
			return nil
		},
	}

	return cmd
}
