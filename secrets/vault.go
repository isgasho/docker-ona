package secrets

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/onaci/docker-ona/config"

	"github.com/docker/cli/cli/streams"
	"github.com/hashicorp/vault/api"
)

func Login(output *streams.Out) (client *api.Client, err error) {
	vConfig := &api.Config{
		Address: fmt.Sprintf("https://%s", config.VaultServer),
	}
	client, err = api.NewClient(vConfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client.SetToken(config.VaultToken)
	// TODO: if a token exists, try it, and if that fails, ask the user to auth
	if config.VaultToken != "" {
		_, err := client.Logical().Read(
			fmt.Sprintf("%s/data/",
				config.VaultPath,
				"asdf",	// it doesn't matter what the secret is to verify that token auth is ok
			))
		//fmt.Printf("mmmm: %#v\n", secret)
		if err == nil {
			return client, nil
		}
	}

	if config.VaultPassword == "" {
		fmt.Fprintf(output, "enter password for: %s as ldap user: %s\n", config.VaultServer, config.VaultUser)

		passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		config.VaultPassword = string(passphrase)
	}
	auth, err := client.Logical().Write(
		fmt.Sprintf("auth/ldap/login/%s", config.VaultUser),
		map[string]interface{}{
			"password": config.VaultPassword,
		})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	config.VaultToken = auth.Auth.ClientToken
	fmt.Fprintf(output, "Auth good: token: %s\n", config.VaultToken)
	return client, err
}

func GetSecret(client *api.Client, path, key, value string) (val string, err error) {
	// Step 2: use vault to get the token for the gitlab server
	client.SetToken(config.VaultToken)
	secret, err := client.Logical().Read(
		fmt.Sprintf("%s/data/%s",
			config.VaultPath,
			key,
		))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("nope: %#v\n", secret)
		return "", err
	}
	//fmt.Printf("hello: %v - %s\n", m["USER"], m["TOKEN"])

	return fmt.Sprintf("%s", m[value]), nil
}
