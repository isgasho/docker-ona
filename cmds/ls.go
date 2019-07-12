package cmds

import (
	"fmt"

	"github.com/onaci/docker-ona/config"
	"github.com/onaci/docker-ona/secrets"

	"github.com/docker/cli/cli/streams"
)

func LsCommand(output *streams.Out) error {
	fmt.Fprintf(output, "List all deployments on %s\n", config.GitlabServer)
	fmt.Fprintf(output, "using  %s\n", config.VaultServer)

	// Step 1: get the user to logged into vault and get a token
	client, err := secrets.Login(output)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Step 2: use vault to get the token for the gitlab server
	gitlabUser, err := secrets.GetSecret(client, config.VaultPath, config.GitlabServer, "USER")
	if err != nil {
		fmt.Printf("Error getting USER: %s\n", err)
		return err
	}
	gitlabToken, err := secrets.GetSecret(client, config.VaultPath, config.GitlabServer, "TOKEN")
	if err != nil {
		fmt.Printf("Error getting TOKEN: %s\n", err)
		return err
	}
	fmt.Printf("hello: %v - %s\n", gitlabUser, gitlabToken)
	// Step 3: ask the gitlab server what projects are configured

	return nil
}
