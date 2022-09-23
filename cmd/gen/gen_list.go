package gen

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/utils"

	"github.com/spf13/cobra"
)

// ListModulesCommand show generate support module types
func ListModulesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List generate module types",
		Long: `list generate module types.

Examples:
  # show generate module types
  goctl gen list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(utils.ListTypeNames(allModules...))
			return nil
		},
	}
}
