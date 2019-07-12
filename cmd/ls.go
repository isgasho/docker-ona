package cmd

import (
	"fmt"

	"github.com/docker/cli/cli/streams"
)

func LsCommand(output *streams.Out, gitlabServer string) error {
	fmt.Fprintf(output, "List all deployments on %s\n", gitlabServer)

	return nil
}
