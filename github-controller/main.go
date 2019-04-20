package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "github-controller",
	Short:         "Manage custom resource for GitHub",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	log.SetOutput(os.Stderr)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}