package cli

import (
	"os"

	"github.com/docker/cli/cli/command"
	"github.com/onaci/docker-ona/cmd"
	"github.com/spf13/cobra"
)

func lsFunc(dockerCli command.Cli) *cobra.Command {
	c := &cobra.Command{
		Use:   "ls",
		Short: "list all deployments",
		RunE: func(cc *cobra.Command, _ []string) error {
			err := cmd.LsCommand(dockerCli.Out())
			if err != nil {
				os.Exit(-1)
			}
			return nil
		},
	}

	return c
}
