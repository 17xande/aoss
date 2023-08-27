package root

import (
	"github.com/17xande/aoss/pkg/cmd/lldp"
	"github.com/17xande/aoss/pkg/cmd/macAddress"
	"github.com/17xande/aoss/pkg/cmd/system"
	"github.com/17xande/aoss/pkg/cmd/utils"
	"github.com/17xande/aoss/pkg/cmd/vlans"
	"github.com/spf13/cobra"
)

func NewCmdRoot() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "aoss <command> <subcommand> [flags]",
		Short: "ArubaOs-Switches CLI",
		Long:  `Work with Aruba Switches running the legacy AOS-S`,
		Example: `
$ aoss macAddress 000000-111111 --host=192.168.0.1
$ aoss lldp remoteDevices 22 --auth=creds.toml
$ aoss vlans 22`,
		Annotations: map[string]string{
			"versionInfo": "version?",
		},
		// Don't need this at the moment, but will keep it here so I don't
		// forget this option exists for now.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().BoolP("help", "H", false, "Show help for command")
	cmd.PersistentFlags().StringP("host", "h", "", "Switch's hostname or IP address (required)")
	cmd.MarkPersistentFlagRequired("host")

	// Groups:
	cmd.AddGroup(&cobra.Group{
		ID:    "services",
		Title: "Switch services",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core commands",
	})

	// Child commands:
	cmd.AddCommand(macAddress.NewCmdMacAddress())
	cmd.AddCommand(lldp.NewCmdLldp())
	cmd.AddCommand(utils.NewCmdUtils())
	cmd.AddCommand(system.NewCmdSystem())
	cmd.AddCommand(vlans.NewCmdVlans())

	return cmd, nil
}
