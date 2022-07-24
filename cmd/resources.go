package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func resourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resources",
		Short: "List of supported resources",
		Long: `list of supported resources. 

Examples:
    goctl resources
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(string(ListResourceNames()))
			return nil
		},
	}

	return cmd
}

// --------------------------------------------------------------------------------------

const (
	// api资源
	api = "api"
	// web资源
	web = "web"
)

// 支持的资源名称列表
var resourceNames = []string{
	api,
	web,
}

// ListResourceNames 资源名称列表
func ListResourceNames() []byte {
	content := []string{"resources list:\n\n"}
	for _, name := range resourceNames {
		content = append(content, name+"\n\n")
	}

	return []byte(strings.Join(content, ""))
}
