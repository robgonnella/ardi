package commands

import (
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/core/platform"
	"github.com/spf13/cobra"
)

func getPlatformListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Long:    cyan("\nList all available platforms"),
		Short:   "List all available platforms",
		Aliases: []string{"search"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			platformCore, err := platform.New(logger)
			if err != nil {
				return
			}
			platformCore.List(query)
		},
	}
	return listCmd
}

func getPlatformCommand() *cobra.Command {
	platCmd := &cobra.Command{
		Use:   "platform",
		Long:  cyan("\nPlatform related commands"),
		Short: "Platform related commands",
	}
	platCmd.AddCommand(getPlatformListCmd())

	return platCmd
}
