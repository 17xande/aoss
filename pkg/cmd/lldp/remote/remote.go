package remote

import "github.com/spf13/cobra"

type remoteOptions struct {
	PortID string
}

func NewCmdRemote() *cobra.Command {
	opts := remoteOptions{}
	cmd := &cobra.Command{
		Use:     "remote",
		Short:   "Query remote LLDP devices",
		Long:    "TODO:",
		Example: "TODO:",
		RunE:    runE,
	}

	cmd.Flags().StringVarP(&opts.PortID, "portID", "p", "", "Port ID")

	return cmd
}

func runE(c *cobra.Command, args []string) error {

	return nil
}
