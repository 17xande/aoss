package ports

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type ports struct {
	Uri      string
	VlanID   int    `json:"vlan_id"`
	PortID   string `json:"port_id"`
	PortMode string `json:"port_mode"`
}

type vlanOptions struct {
	vlanID int
	portID string
}

func NewCmdPorts() *cobra.Command {
	opts := vlanOptions{}

	cmd := &cobra.Command{
		Use:     "ports",
		Short:   "Get or Set VLAN tagging for specified ports",
		Long:    "TODO:",
		Example: "TODO:",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			path := "vlans-ports"
			res := ""
			var err error

			if opts.vlanID == 0 && opts.portID == "" {
				if err != nil {
					return err
				}
			} else if opts.vlanID != 0 && opts.portID != "" {
				path = fmt.Sprintf("%s/%d-%s", path, opts.vlanID, opts.portID)
			}

			res, err = request.GetJson(host, path)
			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.vlanID, "vlanID", "i", 0, "Vlan ID")
	cmd.Flags().StringVarP(&opts.portID, "portID", "p", "", "Port ID")

	return cmd
}
