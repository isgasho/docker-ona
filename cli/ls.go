package cli

import (
	"os"

	"github.com/docker/cli/cli/command"
	"github.com/onaci/docker-ona/cmd"
	"github.com/spf13/cobra"
)

func lsFunc(dockerCli command.Cli) *cobra.Command {
	var gitlabServer string
	c := &cobra.Command{
		Use:   "ls",
		Short: "list all deployments",
		RunE: func(cc *cobra.Command, _ []string) error {
			err := cmd.LsCommand(dockerCli.Out(), gitlabServer)
			if err != nil {
				os.Exit(-1)
			}
			return nil
		},
	}
	flags := c.Flags()
	flags.StringVar(&gitlabServer, "gitlab", "git.ona.im", "Show deployments managed by this gitlab server")

	return c
}
