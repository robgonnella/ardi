package commands

import (
	"fmt"

	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getInstallCmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install all project dependencies",
		Long:  "\nInstall all project dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := util.InitProjectDirectory(port); err != nil {
				return err
			}
			for _, url := range ardiCore.Config.GetBoardURLS() {
				if err := ardiCore.Config.AddBoardURL(url); err != nil {
					logger.WithError(err).Errorf("Failed add board url %s", url)
					return err
				}
			}
			for plat, vers := range ardiCore.Config.GetPlatforms() {
				_, _, err := ardiCore.Platform.Add(fmt.Sprintf("%s@%s", plat, vers))
				if err != nil {
					logger.WithError(err).Errorf("Failed to install %s", plat)
					return err
				}
			}
			for lib, vers := range ardiCore.Config.GetLibraries() {
				_, _, err := ardiCore.Lib.Add(fmt.Sprintf("%s@%s", lib, vers))
				if err != nil {
					logger.WithError(err).Errorf("Failed to install %s", lib)
					return err
				}
			}
			return nil
		},
	}
	return installCmd
}