package lldp

import "github.com/spf13/cobra"

type lldpOptions struct {
	PortID string
	VlanID int
}

func NewCmdLldp() *cobra.Command {
	opts := &lldpOptions{}

	cmd := &cobra.Command{
		Use:     "lldp",
		Short:   "Work with the switch's LLDP service.",
		Long:    longExplainer(),
		Example: "TODO:",
		GroupID: "services",
		RunE:    runE,
	}

	cmd.Flags().IntVarP(&opts.VlanID, "vlan ID", "v", 0, "vlan ID")

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
