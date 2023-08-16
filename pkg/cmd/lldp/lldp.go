package lldp

import (
	"github.com/17xande/aoss/pkg/cmd/lldp/remote"
	"github.com/spf13/cobra"
)

func NewCmdLldp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lldp",
		Short:   "Work with the switch's LLDP service.",
		Long:    longExplainer(),
		Example: "TODO:",
		GroupID: "services",
		RunE:    runE,
	}

	cmd.AddCommand(remote.NewCmdRemote())

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	return nil
}

func longExplainer() string {
	return `
	LLDP (Link Layer Discovery Protocol)
	Used to discover neighbours and their details on the network.
	`
}
