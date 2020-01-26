package commands

import (
	"os"

	"github.com/robgonnella/ardi/ardi"
	"github.com/spf13/cobra"
)

func getCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Delete all ardi global data",
		Long:  "Removes all installed platforms and libraries from ~/.ardi",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Cleaning ardi data directory...")
			if err := os.RemoveAll(ardi.ArdiDir); err != nil {
				logger.WithError(err).Fatalf("Failed to clean ardi directory. You can manually clean all data by removing %s", ardi.ArdiDir)
			}
			logger.Infof("Successfully removed all data from %s", ardi.ArdiDir)
		},
	}
}
