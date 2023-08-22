package status

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/spf13/cobra"
)

type response struct {
	Uri                 string
	Name                string
	SerialNumber        string              `json:"serial_number"`
	FirmwareVersion     string              `json:"firmware_version"`
	HardwareRevision    string              `json:"hardware_revision"`
	ProductModel        string              `json:"product_model"`
	BaseEthernetAddress baseEthernetAddress `json:"base_ethernet_address"`
	TotalMemoryInBytes  int                 `json:"total_memory_in_bytes"`
	TotalPoeConsumption int                 `json:"total_poe_consumption"`
	SysFanStatus        bool                `json:"sys_fan_status"`
}

type baseEthernetAddress struct {
	Version string
	Octets  string
}

func NewCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Get device's system status",
		Long:    "TODO:",
		Example: "TODO:",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			path := "system/status"

			res, err := request.GetJson(host, path)
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	return cmd
}

func GetHardwareRevision(host string, auth *request.Auth) (string, error) {
	path := "system/status"
	result := response{}
	if err := request.GetUnmarshalled(host, path, auth, &result); err != nil {
		return "", err
	}

	return result.HardwareRevision, nil
}
