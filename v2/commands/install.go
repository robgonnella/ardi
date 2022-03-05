package commands

import (
	"fmt"

	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getInstallCmd(env *CommandEnv) *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install all project dependencies",
		Long:  "\nInstall all project dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			boardURLs := env.ArdiCore.CliConfig.Config.BoardManager.AdditionalUrls
			for _, url := range env.ArdiCore.Config.GetBoardURLS() {
				if !util.ArrayContains(boardURLs, url) {
					env.Logger.WithField("board-url", url).Info("Adding board url")
					if err := env.ArdiCore.CliConfig.AddBoardURL(url); err != nil {
						return err
					}
				}
			}
			for plat, vers := range env.ArdiCore.Config.GetPlatforms() {
				_, _, err := env.ArdiCore.Platform.Add(fmt.Sprintf("%s@%s", plat, vers))
				if err != nil {
					return err
				}
			}
			for lib, vers := range env.ArdiCore.Config.GetLibraries() {
				_, _, err := env.ArdiCore.Lib.Add(fmt.Sprintf("%s@%s", lib, vers))
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	return installCmd
}
