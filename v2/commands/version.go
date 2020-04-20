package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints current version of ardi",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ardi: v%s\n", version)
		},
	}
}
