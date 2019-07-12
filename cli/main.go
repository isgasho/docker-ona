package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/onaci/docker-ona/config"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
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

		// TODO: it'd be nice to be able to get the defaults from the plugin config file, but dockerCli.ConfigFile() isn't initilised until the cmdline is parsed..
		//       which also suggests the idea of contexts...
		flags.StringVar(&config.GitlabServer, "gitlab", "git.ona.im", "Show deployments managed by this gitlab server")
		flags.StringVar(&config.VaultServer, "vault", "vault.ona.im", "Use Secrets from vault server")

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
