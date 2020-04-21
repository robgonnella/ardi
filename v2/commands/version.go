package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Long:  cyan("\nPrints current version of ardi"),
		Short: "Prints current version of ardi",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ardi: v%s\n", version)
		},
	}
}
