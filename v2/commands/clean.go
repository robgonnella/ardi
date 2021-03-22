package commands

import (
	"fmt"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Delete project, or global, data directory",
		Long: "\nRemoves all installed platforms and libraries from project " +
			"data directory. If run with \"--global\" all data will be removed " +
			"from ~/.ardi",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := paths.ArdiProjectDataDir
			if global {
				dir = paths.ArdiGlobalDataDir
			}
			logger.Infof("Cleaning ardi data directory: %s", dir)
			if err := util.CleanDataDirectory(dir); err != nil {
				errMsg := err.Error()
				fullErr := fmt.Errorf("%s: You can manually clean all data by removing %s", errMsg, dir)
				return fullErr
			}

			logger.Infof("Successfully removed all data from %s", dir)
			return nil
		},
	}
}
