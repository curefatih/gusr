package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func MockRootCmd() *cobra.Command {
	var root = &cobra.Command{
		Use:   "gusr",
		Short: "Git User Switcher",
	}

	root.AddCommand(listCmd)

	return root
}

func TestListSuccessfully(t *testing.T) {
	rootCmd := MockRootCmd()

	rootCmd.Execute()

}
