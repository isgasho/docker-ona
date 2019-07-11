package cmd

import (
	"fmt"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func lsFunc(dockerCli command.Cli) *cobra.Command {
	var gitlabServer string
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list all deployments",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(dockerCli.Out(), "List all deployments on %s\n", gitlabServer)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&gitlabServer, "gitlab", "git.ona.im", "Show deployments managed by this gitlab server")

	return cmd
}
