package macAddress

import (
	"encoding/json"
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type response struct {
	CollectionResult     `json:"collection_result"`
	MacTableEntryElement []MacTableEntryElement `json:"mac_table_entry_element"`
}

type CollectionResult struct {
	TotalElementsCount    int `json:"total_elements_count"`
	FilteredElementsCount int `json:"filtered_elements_count"`
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
		Short:   "Query mac addresses",
		Long:    "TODO:",
		Example: "TODO:",
		GroupID: "services",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
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

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port ID")
	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "MAC Address")

	return cmd
}

func runE(r *request.Request) error {
	body, err := r.Get()
	if err != nil {
		return fmt.Errorf("could not complete API request: %w", err)
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("can't unmarshal JSON response: %w", err)
	}

	fmt.Printf("%#v\n", res)

	// fmt.Println(res)

	return nil
}
