package system

import (
	"fmt"

	"github.com/17xande/aoss/internal/api/request"
	"github.com/17xande/aoss/pkg/cmd/system/status"
	"github.com/spf13/cobra"
)

func NewCmdSystem() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "system",
		Short:   "Get the switch's system information",
		Long:    "TODO:",
		Example: "TODO:",
		GroupID: "services",
		RunE: func(c *cobra.Command, args []string) error {
			host, _ := c.Flags().GetString("host")
			path := "system"

			res, err := request.GetJson(host, path)
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	cmd.AddCommand(status.NewCmdStatus())

	return cmd
}
