package remote

import (
	"encoding/json"
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type remoteOptions struct {
	port string
}

type response struct {
	Uri                     string
	LocalPort               string `json:"local_port"`
	ChassisType             string `json:"chassis_type"`
	ChassisID               string `json:"chassis_id"`
	PortType                string `json:"port_type"`
	PortID                  string `json:"port_id"`
	PortDescription         string `json:"port_description"`
	SystemName              string `json:"system_name"`
	SystemDescription       string `json:"system_description"`
	PvID                    int
	CapabilitiesSupported   Capabilities
	CapabilitiesEnambled    Capabilities
	RemoteManagementAddress `json:"remote_management_address"`
	PoePlusInfo             `json:"poe_plus_info"`
}

type Capabilities struct {
	Repeater        bool `json:"repeater"`
	Bridge          bool `json:"bridge"`
	WlanAccessPoint bool `json:"wlan_access_point"`
	Router          bool `json:"router"`
	Telephone       bool `json:"telephone"`
	CableDevice     bool `json:"cable_device"`
	StationOnly     bool `json:"station_only"`
}

type RemoteManagementAddress []struct {
	Type    string
	Address string
}

type PoePlusInfo struct {
	PoeDeviceType         string `json:"poe_device_type"`
	PowerSource           string `json:"power_source"`
	PowerPriority         string `json:"power_priority"`
	RequestedPowerInWatts string `json:"requested_power_in_watts"`
	ActualPowerInWatts    string `json:"actual_power_in_watts"`
}

func NewCmdRemote() *cobra.Command {
	opts := remoteOptions{}

	cmd := &cobra.Command{
		Use:     "remote",
		Short:   "Query remote LLDP devices",
		Long:    "TODO:",
		Example: "TODO:",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			var r *request.Request

			if opts.port != "" {
				r = request.New(host, "lldp/remote-device/"+opts.port)
			} else {
				r = request.New(host, "lldp/remote-device")
			}

			return runE(r)
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port ID")

	return cmd
}

func runE(r *request.Request) error {
	body, err := r.Get()
	if err != nil {
		return fmt.Errorf("can't complete API request: %w", err)
	}

	// fmt.Println(body)

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("can't unmarshal JSON response: %w", err)
	}

	fmt.Printf("%#v\n", res)

	return nil
}
