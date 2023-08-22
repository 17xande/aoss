package remote

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type remoteOptions struct {
	port   string
	pretty bool
}

type response struct {
	request.CollectionResult `json:"collection_result"`
	LldpRemoteDeviceElement  []LldpRemoteDeviceElement `json:"lldp_remote_device_element"`
}

type LldpRemoteDeviceElement struct {
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
	CapabilitiesSupported   Capabilities `json:"capabilities_supported"`
	CapabilitiesEnabled     Capabilities `json:"capabilities_enabled"`
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
			path := "lldp/remote-device"

			if opts.port != "" {
				path += "/" + opts.port
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

	return cmd
}

func GetLldpRemote(host, port string, auth *request.Auth) (string, error) {
	path := "lldp/remote-device/" + port
	result := LldpRemoteDeviceElement{}
	if err := request.GetUnmarshalled(host, path, auth, &result); err != nil {
		return "", err
	}

	if !result.CapabilitiesEnabled.Bridge {
		// This is probably not a switch, return blank
		return "", nil
	}

	for _, ip := range result.RemoteManagementAddress {
		if ip.Type == "AFM_IP4" {
			return ip.Address, nil
		}
	}
	return "", fmt.Errorf("no lldp IPv4 address found on this port")
}
