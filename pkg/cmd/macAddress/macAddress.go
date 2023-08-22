package macAddress

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type Response struct {
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
			path := "mac-table"

			if opts.port != "" && opts.mac != "" {
				return fmt.Errorf("can't query both port and mac-address. Choose one or the other")
			}

			if opts.port == "" && opts.mac == "" {
				path = "mac-table"
			} else if opts.port != "" {
				path = fmt.Sprintf("ports/%s/mac-table", opts.port)
			} else if opts.mac != "" {
				path = "mac-table/" + opts.mac
			}

			res, err := request.GetJson(host, path)
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port ID")
	cmd.Flags().StringVarP(&opts.mac, "mac", "m", "", "MAC Address")

	return cmd
}

func GetPortWithMac(host, mac string, auth *request.Auth) (string, error) {
	path := "mac-table/" + mac
	result := MacTableEntryElement{}
	if err := request.GetUnmarshalled(host, path, auth, &result); err != nil {
		return "", err
	}

	return result.PortID, nil
}

func GetMacCountAtPort(host, port string, auth *request.Auth) (int, error) {
	path := fmt.Sprintf("ports/%s/mac-table", port)
	result := Response{}
	if err := request.GetUnmarshalled(host, path, auth, &result); err != nil {
		return 0, err
	}

	return result.CollectionResult.TotalElementsCount, nil
}
