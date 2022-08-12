package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/cmd/gen"
)

func genGinCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen <type>",
		Short: "Generate code",
		Long: `generate code.

Examples:
  # show generate types
  goctl gen list

`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		gen.ListTypesCommand(),
		gen.ApiCommand(),
		gen.WebCommand(),
		gen.UserCommand(),
	)

	return cmd
}
