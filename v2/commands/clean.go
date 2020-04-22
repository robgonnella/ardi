package commands

import (
	"os"

	"github.com/robgonnella/ardi/v2/paths"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Delete all ardi global data",
		Long:  cyan("\nRemoves all installed platforms and libraries from ~/.ardi"),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			logger.Info("Cleaning ardi data directory...")
			if err := os.RemoveAll(paths.ArdiDataDir); err != nil {
				logger.WithError(err).Errorf("Failed to clean ardi directory. You can manually clean all data by removing %s", paths.ArdiDataDir)
			}
			logger.Infof("Successfully removed all data from %s", paths.ArdiDataDir)
		},
	}
}
