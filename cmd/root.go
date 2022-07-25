package cmd

import (
	"github.com/spf13/cobra"
)

// Version 命令版本号
const Version = "0.0.1"

// NewRootCMD 命令入口
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "goctl",
		Long:          "go language development tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       Version,
	}

	cmd.AddCommand(
		genGinCommand(),
		resourcesCommand(),
		//addCommand(),
		//deleteCommand(),
	)
	return cmd
}
