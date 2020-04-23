package commands

import (
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/core/platform"
	"github.com/spf13/cobra"
)

func getPlatformCommand() *cobra.Command {
	platCmd := &cobra.Command{
		Use:   "platforms",
		Long:  cyan("\nList all available platforms"),
		Short: "List all available platforms",
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
			defer platformCore.Client.Connection.Close()
			platformCore.List(query)
		},
	}
	return platCmd
}
