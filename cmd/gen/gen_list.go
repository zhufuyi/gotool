package gen

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	// GenTypeDao model 类型
	GenTypeModel = "model"

	// GenTypeDao dao 类型
	GenTypeDao = "dao"

	// GenTypeHandler handler 类型
	GenTypeHandler = "handler"

	// GenTypeUser user 类型
	GenTypeUser = "user"
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
				GenTypeModel,
				GenTypeDao,
				GenTypeHandler,
				GenTypeUser,
			))
			return nil
		},
	}

	return cmd
}
