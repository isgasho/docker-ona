package cmd

import (
	"fmt"

	"github.com/docker/cli/cli/streams"
)

// TODO: move to config.go
var GitlabServer, VaultServer string

func LsCommand(output *streams.Out) error {
	fmt.Fprintf(output, "List all deployments on %s\n", GitlabServer)
	fmt.Fprintf(output, "using  %s\n", VaultServer)

	return nil
}
