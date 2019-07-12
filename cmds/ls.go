package cmds

import (
	"fmt"
	"log"

	"github.com/onaci/docker-ona/config"
	"github.com/onaci/docker-ona/secrets"

	"github.com/docker/cli/cli/streams"
	gitlab "github.com/xanzy/go-gitlab"
)

func LsCommand(output *streams.Out) error {
	fmt.Fprintf(output, "List all deployments on %s:\n", config.GitlabServer)
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

	opt := &gitlab.ListGroupsOptions{	}

	// Get the first page with projects.
	gs, _, err := git.Groups.ListGroups(opt)
	if err != nil {
		log.Fatal(err)
	}

	// List all the projects we've found so far.
	for _, g := range gs {
		fmt.Printf("%s\n", g.Name)

		opt := &gitlab.ListProjectsOptions{	}

		// Get the first page with projects.
		ps, _, err := git.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			if p.Namespace.Name == g.Name {
				fmt.Printf("\t%s\n", p.Name)
			}
		}

	}


	return nil
}
