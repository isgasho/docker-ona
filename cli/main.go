package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/onaci/docker-ona/config"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/spf13/cobra"
)

// dockerCli has these IO hooks:
// In() *streams.In
// Out() *streams.Out
// Err() io.Writer

var dockerPluginCommand = "ona"

func RegisterCommands() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		apiversion := &cobra.Command{
			Use:   "apiversion",
			Short: "Print the API version of the server",
			RunE: func(_ *cobra.Command, _ []string) error {
				cli := dockerCli.Client()
				ping, err := cli.Ping(context.Background())
				if err != nil {
					return err
				}
				fmt.Println(ping.APIVersion)
				return nil
			},
		}

		exitStatus2 := &cobra.Command{
			Use:   "exitstatus2",
			Short: "Exit with status 2",
			RunE: func(_ *cobra.Command, _ []string) error {
				fmt.Fprintln(dockerCli.Err(), "Exiting with error status 2")
				os.Exit(2)
				return nil
			},
		}

		var who string
		cmd := &cobra.Command{
			Use:   dockerPluginCommand,
			Short: "A Docker commandline plugin that provisions complete system deployments using the infrastructure we have at CSIRO O&A.",
			// This is redundant but included to exercise
			// the path where a plugin overrides this
			// hook.
			PersistentPreRunE: plugin.PersistentPreRunE,
			RunE: func(cmd *cobra.Command, args []string) error {
				if who == "" {
					who, _ = dockerCli.ConfigFile().PluginConfig("helloworld", "who")
				}
				if who == "" {
					who = "World"
				}

				fmt.Fprintf(dockerCli.Out(), "Hello %s!\n", who)
				fmt.Fprintf(dockerCli.Out(), "gitlab: %s\n", config.GitlabServer)
				fmt.Fprintf(dockerCli.Out(), "vault: %s\n", config.VaultServer)

				dockerCli.ConfigFile().SetPluginConfig("helloworld", "lastwho", who)
				return dockerCli.ConfigFile().Save()
			},
		}
		flags := cmd.Flags()

		// Get the defaults from .docker/config.json file
		//       which also suggests the idea of contexts...
		defaultGitlabServer, err := getConfigValue(dockerCli, "gitlab", "git.ona.im")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		defaultVaultServer, err := getConfigValue(dockerCli, "vault", "vault.ona.im")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		flags.StringVar(&config.GitlabServer, "gitlab", defaultGitlabServer, "Show deployments managed by this gitlab server")
		flags.StringVar(&config.VaultServer, "vault", defaultVaultServer, "Use Secrets from vault server")

		cmd.AddCommand(lsFunc(dockerCli))
		cmd.AddCommand(apiversion, exitStatus2)
		return cmd
	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "CSIRO Oceans & Atmosphere.",
			Version:       "ONA deployment v0.1.0",
		})
}

var configFile *configfile.ConfigFile
func getConfigValue(dockerCli command.Cli, name, defaultValue string) (string, error) {
	if configFile == nil {
		configFile = cliconfig.LoadDefaultConfigFile(dockerCli.Err())
		if configFile == nil {
			return "", errors.New("Failed to load Docker config.json")
		}
	}
	value, ok := configFile.PluginConfig(dockerPluginCommand, name)
	if !ok {
		value = defaultValue
		configFile.SetPluginConfig(dockerPluginCommand, name, value)
		err := configFile.Save()
		return value, err
	}
	return value, nil
}