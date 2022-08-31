package cmd

import (
	"github.com/zhufuyi/goctl/cmd/covert"

	"github.com/spf13/cobra"
)

func convertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "covert <type>",
		Short: "covert resources",
		Long: `Covert resources.

Examples:
  # show covert resource types
  goctl covert list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		covert.ListTypesCommand(),
		covert.SQL2GormCommand(),
		covert.JSON2StructCommand(),
		covert.Yaml2StructCommand(),
	)

	return cmd
}
