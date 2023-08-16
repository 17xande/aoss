package macAddress

import (
	"fmt"
	"net/url"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type Response struct {
	Uri        url.URL
	MacAddress string `json:"mac_address"`
	PortID     string `json:"port_id"`
	VlanID     int    `json:"vlan_id"`
}

type macAddressOptions struct {
	port string
	mac  string
}

func NewCmdMacAddress() *cobra.Command {
	opts := &macAddressOptions{}

	cmd := &cobra.Command{
		Use:     "macAddress",
		Short:   "Query mac addresses",
		Long:    "TODO:",
		Example: "TODO:",
		GroupID: "services",
		RunE: func(c *cobra.Command, args []string) error {
			host, err := c.Flags().GetString("host")
			if err != nil {
				return fmt.Errorf("can't get host flag: %w", err)
			}

			var r *request.Request

			if opts.port != "" && opts.mac != "" {
				return fmt.Errorf("can't query both port and mac-address. Choose one or the other")
			}

			if opts.port != "" {
				r = request.New(host, fmt.Sprintf("ports/%s/mac-table", opts.port))
			}

			if opts.mac != "" {
				r = request.New(host, "mac-table/"+opts.mac)
			}

			if opts.port == "" && opts.mac == "" {
				r = request.New(host, "mac-table")
			}

			return runE(r)
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "port")
	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "mac address")

	return cmd
}

func runE(r *request.Request) error {
	res, err := r.Get()
	if err != nil {
		return fmt.Errorf("could not complete API request: %w", err)
	}

	fmt.Println(res)

	return nil
}
