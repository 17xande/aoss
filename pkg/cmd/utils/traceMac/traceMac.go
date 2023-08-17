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

func getPortWithMac(host, mac string) (string, error) {
	r := request.New(host, "mac-table/"+mac, macAddress.MacTableEntryElement{})
	res, err := r.GetUnmarshalled()

	if err != nil {
		return "", fmt.Errorf("couldn't complete API request: %w", err)
	}

	return res.PortID, nil
}

func getMacCountAtPort(host, port string) (int, error) {
	r := request.New(host, fmt.Sprintf("ports/%s/mac-table", port), macAddress.Response{})
	res, err := r.GetUnmarshalled()

	if err != nil {
		return 0, fmt.Errorf("couldn't complete API request: %w", err)
	}

	return res.CollectionResult.TotalElementsCount, nil
}

func getLldpRemote(host, port string) (string, error) {
	r := request.New(host, "lldp/remote-device/"+port, remote.LldpRemoteDeviceElement{})
	res, err := r.GetUnmarshalled()

	if err != nil {
		return "", fmt.Errorf("couldn't complete API request: %w", err)
	}

	if !res.CapabilitiesEnabled.Bridge {
		// This is probably not a switch, return blank
		return "", nil
	}

	for _, ip := range res.RemoteManagementAddress {
		if ip.Type == "AFM_IP4" {
			return ip.Address, nil
		}
	}
	return "", fmt.Errorf("no lldp IPv4 address found on this port")
}

func trace(host, mac string) (string, string, error) {
	for {
		// Get port that has this MAC address registered.
		port, err := getPortWithMac(host, mac)
		if err != nil {
			return "", "", err
		}

		fmt.Printf("Port %s on %s has %s\n", port, host, mac)

		// Check to see if this is the only mac address on that port.
		count, err := getMacCountAtPort(host, port)
		if err != nil {
			return "", "", err
		}

		if count == 1 {
			// Only one device at this port, this is not a switch,
			// this is the device we're looking for.
			return host, port, nil
		}

		newHost, err := getLldpRemote(host, port)
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
