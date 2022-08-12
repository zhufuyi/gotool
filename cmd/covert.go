package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/cmd/covert"
)

func convertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "covert <type>",
		Short: "resource type conversion",
		Long: `Resource type conversion.

Examples:
  # show covert types
  goctl covert list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		covert.ListTypesCommand(),
		covert.Sql2GormCommand(),
		covert.JSON2StructCommand(),
		covert.Yaml2StructCommand(),
	)

	return cmd
}
