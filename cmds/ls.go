package cmds

import (
	"fmt"
	"log"

	"github.com/onaci/docker-ona/config"
	"github.com/onaci/docker-ona/secrets"

	"github.com/docker/cli/cli/streams"
	"github.com/xanzy/go-gitlab"
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
	fmt.Printf("hello %s: %v - %s\n", config.GitlabServer, gitlabUser, gitlabToken)
	// Step 3: ask the gitlab server what projects are configured
	git := gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(fmt.Sprintf("https://%s", config.GitlabServer))

	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		ps, resp, err := git.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			fmt.Printf("Found project: %s", p.Name)
		}

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return nil
}
