package gen

import (
	"fmt"

	"github.com/zhufuyi/goctl/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	// GenTypeModel model 类型
	GenTypeModel = "model"

	// GenTypeDao dao 类型
	GenTypeDao = "dao"

	// GenTypeHandler handler 类型
	GenTypeHandler = "handler"

	// GenTypeHTTP http 类型
	GenTypeHTTP = "http"
)

var genTypes = []string{GenTypeModel, GenTypeDao, GenTypeHandler, GenTypeHTTP}

// ListTypesCommand show generate support types
func ListTypesCommand() *cobra.Command {
	return &cobra.Command{
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
			fmt.Println(utils.ListTypeNames(genTypes...))
			return nil
		},
	}
}
