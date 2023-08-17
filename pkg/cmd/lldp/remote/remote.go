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

			if opts.port != "" {
				r := request.New(host, "lldp/remote-device/"+opts.port, LldpRemoteDeviceElement{})
				if opts.pretty {
					res, err := r.GetPretty()
					if err != nil {
						return err
					}
					fmt.Println(res)
					return nil
				}
				return runE(r)
			} else {
				r := request.New(host, "lldp/remote-device", response{})
				if opts.pretty {
					res, err := r.GetPretty()
					if err != nil {
						return err
					}
					fmt.Println(res)
					return nil
				}
				return runE(r)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port ID")
	cmd.Flags().BoolVarP(&opts.pretty, "pretty", "P", false, "Pretty print")

	return cmd
}

func runE[T response | LldpRemoteDeviceElement](r *request.Request[T]) error {
	res, err := r.GetUnmarshalled()
	if err != nil {
		return fmt.Errorf("couldn't complete API request: %w", err)
	}

	fmt.Printf("%#v\n", res)
	return nil
}
