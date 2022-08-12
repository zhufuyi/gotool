package gen

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/utils"
)

const (
	genTypeApi  = "api"
	genTypeWeb  = "web"
	genTypeUser = "user"
)

// ListTypesCommand show generate support types
func ListTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List generate types",
		Long: `list generate types.

Examples:
  # show generate types
  goctl gen list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(utils.ListTypeNames(
				genTypeApi,
				genTypeWeb,
				genTypeUser,
			))
			return nil
		},
	}

	return cmd
}
