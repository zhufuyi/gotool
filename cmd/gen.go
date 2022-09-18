package cmd

import (
	"github.com/zhufuyi/goctl/cmd/gen"

	"github.com/spf13/cobra"
)

func genCodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen <module>",
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
		gen.ListModulesCommand(),
		gen.ModelCommand(),
		gen.DaoCommand(),
		gen.HandlerCommand(),
		gen.HTTPCommand(),
		gen.ProtoCommand(),
	)

	return cmd
}
