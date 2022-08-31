package cmd

import (
	"github.com/zhufuyi/goctl/cmd/gen"

	"github.com/spf13/cobra"
)

func genWebCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen <type>",
		Short: "Generate codes",
		Long: `generate codes.

Examples:
  # show generate code types
  goctl gen list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		gen.ListTypesCommand(),
		gen.ModelCommand(),
		gen.DaoCommand(),
		gen.HandlerCommand(),
		gen.HTTPCommand(),
	)

	return cmd
}
