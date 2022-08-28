package gen

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	// GenTypeDao dao 类型
	GenTypeDao = "dao"
	// GenTypeDao dao 类型
	GenTypeModel = "model"
	// GenTypeApi api 类型
	GenTypeApi = "api"
	// GenTypeWeb web 类型
	GenTypeWeb = "web"
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
				GenTypeDao,
				GenTypeApi,
				GenTypeWeb,
				GenTypeUser,
			))
			return nil
		},
	}

	return cmd
}
