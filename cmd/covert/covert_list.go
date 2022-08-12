package covert

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/utils"
)

const (
	covertTypeSql2gorm    = "sql"
	covertTypeJSON2Struct = "json"
	covertTypeYaml2Struct = "yaml"
)

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
			fmt.Println(utils.ListTypeNames(
				covertTypeSql2gorm,
				covertTypeJSON2Struct,
				covertTypeYaml2Struct,
			))
			return nil
		},
	}

	return cmd
}
