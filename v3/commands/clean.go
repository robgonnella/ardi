package commands

import (
	"github.com/robgonnella/ardi/v3/paths"
	"github.com/robgonnella/ardi/v3/util"
	"github.com/spf13/cobra"
)

func newCleanCmd(env *CommandEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Delete project data directory",
		Long: "\nRemoves all installed platforms and libraries from project " +
			"data directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := paths.ArduinoCliProjectDataDir
			env.Logger.Infof("Cleaning ardi data directory: %s", dir)
			util.CleanDataDirectory(dir)
			env.Logger.Infof("Successfully removed all data from %s", dir)
			return nil
		},
	}
}
