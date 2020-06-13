package commands

import (
	"fmt"

	"github.com/robgonnella/ardi/v2/version"

	"github.com/spf13/cobra"
)

func getVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Long:  "\nPrints current version of ardi",
		Short: "Prints current version of ardi",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ardi: v%s\n", version.VERSION)
		},
	}
}
