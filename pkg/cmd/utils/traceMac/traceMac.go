package traceMac

import (
	"github.com/spf13/cobra"
)

type traceMacOptions struct {
	mac string
}

type response struct {
}

func NewCmdTraceMac() *cobra.Command {
	opts := traceMacOptions{}

	cmd := &cobra.Command{
		Use:     "traceMac",
		Short:   "Trace a MAC address through the switches' MAC tables.",
		Long:    "TODO:",
		Example: "TODO:",
		RunE: func(c *cobra.Command, args []string) error {

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "MAC Address")
	cmd.MarkFlagRequired("mac")

	return cmd
}

func getPortWithMac(host, mac string) (int, error) {
	// r := request.New(host, "mac-table/"+mac)
	return 0, nil
}
