package cmds

import (
	"fmt"
	"log"

	"github.com/onaci/docker-ona/config"
	"github.com/onaci/docker-ona/secrets"

	"github.com/docker/cli/cli/streams"
	gitlab "github.com/xanzy/go-gitlab"
)

func CreateCommand(output *streams.Out, params []string) error {
	fmt.Fprintf(output, "Create a new deployment %s on %s:\n", params[0], config.GitlabServer)
	//fmt.Fprintf(output, "using  %s\n", config.VaultServer)

	// Step 1: get the user to logged into vault and get a token
	client, err := secrets.Login(output)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Step 2: use vault to get the token for the gitlab server
	gitlabToken, err := secrets.GetSecret(client, config.VaultPath, config.GitlabServer, "TOKEN")
	if err != nil {
		fmt.Printf("Error getting TOKEN: %s\n", err)
		return err
	}
	//fmt.Printf("hello %s: %v - %s\n", config.GitlabServer, gitlabUser, gitlabToken)
	// Step 3: ask the gitlab server what projects are configured
	git := gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(fmt.Sprintf("https://%s", config.GitlabServer))

	// make a group

	// Create new project
	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String("My Project"),
		Description:          gitlab.String("Just a test project to play with"),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
		ImportURL:			gitlab.String("https://github.com/onaci/swarm-infra"),
	}
	project, _, err := git.Projects.CreateProject(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Something %#v", project)

	return nil
}
