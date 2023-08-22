package traceMac

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/17xande/aoss/pkg/cmd/lldp/remote"
	"github.com/17xande/aoss/pkg/cmd/macAddress"
	"github.com/spf13/cobra"
)

type traceMacOptions struct {
	mac string
}

func NewCmdTraceMac() *cobra.Command {
	opts := traceMacOptions{}

	cmd := &cobra.Command{
		Use:     "traceMac",
		Short:   "Trace a MAC address through the switches' MAC tables.",
		Long:    "TODO:",
		Example: "TODO:",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			host, port, err := trace(host, opts.mac)
			if err != nil {
				return err
			}

			fmt.Printf("Device with MAC %s found at host %s port %s\n", opts.mac, host, port)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "MAC Address")
	cmd.MarkFlagRequired("mac")

	return cmd
}

func trace(host, mac string) (string, string, error) {
	for {
		auth := request.Auth{
			Host: host,
		}

		if err := auth.Login(); err != nil {
			return "", "", fmt.Errorf("could not authenticate: %w", err)
		}
		defer auth.Logout()

		// Get port that has this MAC address registered.
		port, err := macAddress.GetPortWithMac(host, mac, &auth)
		if err != nil {
			return "", "", err
		}

		fmt.Printf("Port %s on %s has %s\n", port, host, mac)

		// Check to see if this is the only mac address on that port.
		count, err := macAddress.GetMacCountAtPort(host, port, &auth)
		if err != nil {
			return "", "", err
		}

		if count == 1 {
			// Only one device at this port, this is not a switch,
			// this is the device we're looking for.
			return host, port, nil
		}

		newHost, err := remote.GetLldpRemote(host, port, &auth)
		if err != nil {
			return "", "", err
		}

		if newHost == "" {
			// No new lldp host at this port, we've found the device.
			return host, port, nil
		}

		host = newHost
		fmt.Printf("Trying new host %s at port %s\n", host, port)
	}
}
