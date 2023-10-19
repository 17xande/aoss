package main

import (
	"fmt"
	"os"

	"github.com/17xande/aoss/pkg/cmd/root"
)

func main() {
	rootCmd, err := root.NewCmdRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %s\n", err)
		os.Exit(1)
	}

	rootCmd.Execute()
}
