package ports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	vlanID  int
	portID  string
	tagging string
}

func NewCmdPorts() *cobra.Command {
	opts := vlanOptions{}

	cmd := &cobra.Command{
		Use:     "ports",
		Short:   "Get or Set VLAN tagging for specified ports",
		Long:    "TODO:",
		Example: "TODO:",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if opts.tagging != "" && opts.tagging != "tagged" && opts.tagging != "untagged" {
				return fmt.Errorf("Invalid tagging option. Must be 'tagged' or 'untagged'")
			}

			host, _ := c.Flags().GetString("host")
			config, _ := c.Flags().GetBool("config")
			if config && opts.tagging == "" {
				return fmt.Errorf("no tagging option supplied. Include --tagging tagged OR untagged")
			}
			tagging := "POM_TAGGED"
			if opts.tagging == "untagged" {
				tagging = "POM_UNTAGGED"
			}

			path := "vlans-ports"
			res := ""
			var err error

			if opts.vlanID == 0 && opts.portID == "" {
				if err != nil {
					return err
				}
			} else if opts.vlanID != 0 && opts.portID != "" {
				if config {
					return putConfig(opts.vlanID, opts.portID, tagging, host, path, tagging)
				}
				path = fmt.Sprintf("%s/%d-%s", path, opts.vlanID, opts.portID)
			}

			res, err = request.GetJson(host, path)
			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.vlanID, "vlanID", "i", 0, "Vlan ID")
	cmd.Flags().StringVarP(&opts.portID, "portID", "p", "", "Port ID")
	cmd.Flags().StringVarP(&opts.tagging, "tagging", "t", "", "tagged|untagged")

	return cmd
}

func putConfig(vlanID int, portID, tagging, host, path, tag string) error {
	auth := request.Auth{
		Host: host,
	}

	if err := auth.Login(); err != nil {
		return fmt.Errorf("could not authenticate: %w", err)
	}
	if auth.Cookie.Raw != "" {
		defer auth.Logout()
	}

	p := ports{
		Uri:      fmt.Sprintf("vlans-ports/%d-%s", vlanID, portID),
		VlanID:   vlanID,
		PortID:   portID,
		PortMode: tag,
	}

	res, err := request.PostUnmarshalled(host, path, &auth, &p)

	if err != nil {
		return fmt.Errorf("could not PUT ports request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("can't read request body: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %d; body:%s\n", res.StatusCode, body)
	}

	var indented bytes.Buffer
	if err := json.Indent(&indented, body, "", "  "); err != nil {
		return err
	}

	fmt.Println(indented.String())

	return nil
}
