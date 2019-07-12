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
	domain := params[0]
	fmt.Fprintf(output, "Create a new deployment %s on %s:\n", domain, config.GitlabServer)
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

	group := createGroup(git, domain)
	createProject(git, group, "swarm-infra", "https://github.com/onaci/swarm-infra")

	return nil
}

func createProject(git *gitlab.Client, group *gitlab.Group, name, url string) *gitlab.Project {
	// see if the group already exists
	psearch, _, err := git.Projects.GetProject((interface{}) (group.Name +"/"+name), &gitlab.GetProjectOptions{})
	if err == nil && psearch != nil {
		return psearch
	}

	// Create new project
	p := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String(name),
		NamespaceID:		&group.ID,
		Description:          gitlab.String("What every project needs :)"),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
		ImportURL:			gitlab.String(url),
		Mirror: gitlab.Bool(true),
		MirrorTriggerBuilds: gitlab.Bool(true),
	}
	project, _, err := git.Projects.CreateProject(p)
	if err != nil {
		log.Fatalf("Creating project %s", err)
	}
	fmt.Printf("Something %#v", project)
	return project
}

func createGroup(git *gitlab.Client, domain string) *gitlab.Group {
	// see if the group already exists
	gsearch, _, err := git.Groups.SearchGroup(domain)
	if err != nil {
		log.Fatalf("Searching for group %s", err)
	}
	if len(gsearch) > 0 {
		return gsearch[0]
	}

	// make a group
	g := &gitlab.CreateGroupOptions{
		Name:                 gitlab.String(domain),
		Path:                 gitlab.String(domain),
		Description:          gitlab.String("auto-created by docker ona create"),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
	}
	group, _, err := git.Groups.CreateGroup(g)
	if err != nil {
		log.Fatalf("Creating group %s", err)
	}
	fmt.Printf("Something %#v", group)
	return group
}