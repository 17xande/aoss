package macAddress

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type response struct {
	request.CollectionResult `json:"collection_result"`
	MacTableEntryElement     []MacTableEntryElement `json:"mac_table_entry_element"`
}

type MacTableEntryElement struct {
	Uri        string
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
		Short:   "Query mac addresses.",
		Long:    "TODO:",
		Example: "TODO:",
		GroupID: "services",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")

			if opts.port != "" && opts.mac != "" {
				return fmt.Errorf("can't query both port and mac-address. Choose one or the other")
			}

			if opts.port != "" {
				r := request.New(host, fmt.Sprintf("ports/%s/mac-table", opts.port), response{})
				return runE(r)
			}

			if opts.mac != "" {
				r := request.New(host, "mac-table/"+opts.mac, MacTableEntryElement{})
				return runE(r)
			}

			if opts.port == "" && opts.mac == "" {
				r := request.New(host, fmt.Sprintf("ports/%s/mac-table", opts.port), response{})
				return runE(r)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port ID")
	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "MAC Address")

	return cmd
}

func runE[T response | MacTableEntryElement](r *request.Request[T]) error {
	res, err := r.GetUnmarshalled()
	if err != nil {
		return fmt.Errorf("couldn't complete API request: %w", err)
	}

	fmt.Printf("%#v\n", res)
	return nil
}
