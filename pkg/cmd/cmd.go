package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/cmd/create"
)

// NewKubectlCommand creates the `kubectl` command and its nested children.
func NewKubectlCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "kubectl",
		Short: "kubectl controls the Kubernetes cluster manager",
	}
	rootCmd.AddCommand(create.NewCmdCreate())
	return rootCmd
}
