package vlans

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type VlanElement struct {
	Uri              string
	VlanID           int
	Name             string
	Status           string
	Type             string
	IsVoiceEnabled   bool
	IsJumboEnabled   bool
	IsManagementVlan bool
	IsPrimaryVlan    bool
}

type vlanOptions struct {
	vlanID int
}

func NewCmdVlans() *cobra.Command {
	opts := vlanOptions{}

	cmd := &cobra.Command{
		Use:     "vlans",
		Short:   "Get or Set VLAN tagging on ports",
		Long:    "TODO:",
		Example: "TODO",
		GroupID: "services",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			path := "vlans"

			if opts.vlanID != 0 {
				path = fmt.Sprintf("%s/%d", path, opts.vlanID)
			}

			res, err := request.GetJson(host, path)
			if err != nil {
				return err
			}

			fmt.Println(res)

			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.vlanID, "vlanID", "i", 0, "Vlan ID")

	return cmd
}
