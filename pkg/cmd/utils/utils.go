package utils

import (
	"github.com/17xande/aoss/pkg/cmd/utils/traceMac"
	"github.com/spf13/cobra"
)

func NewCmdUtils() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "utils",
		Short:   "Utilities to automate switch.",
		Long:    "TODO:",
		Example: "TODO:",
		GroupID: "core",
	}

	cmd.AddCommand(traceMac.NewCmdTraceMac())
	return cmd
}
