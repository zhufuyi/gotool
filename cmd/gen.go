package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zhufuyi/goctl/cmd/gen"
)

func genWebCommand() *cobra.Command {
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
		gen.ModelCommand(),
		gen.DaoCommand(),
		gen.ApiCommand(),
		gen.WebCommand(),
		gen.UserCommand(),
	)

	return cmd
}
