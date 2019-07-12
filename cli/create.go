package cli

import (
	"os"

	"github.com/docker/cli/cli/command"
	"github.com/onaci/docker-ona/cmds"
	"github.com/spf13/cobra"
)

func createFunc(dockerCli command.Cli) *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "create a deployment",
		RunE: func(cc *cobra.Command, params []string) error {
			err := cmds.CreateCommand(dockerCli.Out(), params)
			if err != nil {
				os.Exit(-1)
			}
			saveVaultToken()
			return nil
		},
	}

	return c
}
