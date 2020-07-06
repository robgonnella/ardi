package commands

import (
	"os"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getCleanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Delete project, or global, data directory",
		Long: "\nRemoves all installed platforms and libraries from project " +
			"data directory. If run with \"--global\" all data will be removed " +
			"from ~/.ardi",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Infof("Cleaning ardi data directory: %s", dataDir)
			if err := util.CleanDataDirectory(dataDir); err != nil {
				logger.WithError(err).Errorf("Failed to clean ardi directory. You can manually clean all data by removing %s", dataDir)
				return err
			}

			if !global {
				logger.Info("Cleaning ardi build config")
				if err := os.RemoveAll(paths.ArdiProjectBuildConfig); err != nil {
					logger.WithError(err).Errorf("Failed to remove %s", paths.ArdiProjectBuildConfig)
					return err
				}
			}

			logger.Infof("Successfully removed all data from %s", dataDir)
			return nil
		},
	}
}
