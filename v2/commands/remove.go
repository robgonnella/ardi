package commands

import "github.com/spf13/cobra"

func getRemovePlatformCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nRemove platform(s) from project",
		Short:   "Remove platform(s) from project",
		Aliases: []string{"platform"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, p := range args {
				logger.Infof("Removing platform: %s", p)
				removed, err := ardiCore.Platform.Remove(p)
				if err != nil {
					return err
				}
				logger.Infof("Removed %s", removed)
				if err := ardiCore.Config.RemovePlatform(removed); err != nil {
					return err
				}
				logger.Info("Udated config")
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveBuildCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "builds",
		Long:    "\nRemove build config from project",
		Short:   "Remove build config from project",
		Aliases: []string{"build"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, b := range args {
				if err := ardiCore.Config.RemoveBuild(b); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveLibCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nRemove libraries from project",
		Short:   "Remove libraries from project",
		Aliases: []string{"libs", "lib", "library"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, l := range args {
				logger.Infof("Removing library: %s", l)
				if err := ardiCore.Lib.Remove(l); err != nil {
					return err
				}
				logger.Infof("Removed %s", l)
				if err := ardiCore.Config.RemoveLibrary(l); err != nil {
					return err
				}
				logger.Info("Updated config")
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveBoardURLCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "board-urls",
		Long:    "\nRemove board urls from project",
		Short:   "Remove board urls from project",
		Aliases: []string{"board-url"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, url := range args {
				if err := ardiCore.Config.RemoveBoardURL(url); err != nil {
					return err
				}
				if err := ardiCore.CliConfig.RemoveBoardURL(url); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveCmd() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove project dependencies",
		Long:  "\nRemove project dependencies",
	}
	removeCmd.AddCommand(getRemovePlatformCmd())
	removeCmd.AddCommand(getRemoveBuildCmd())
	removeCmd.AddCommand(getRemoveLibCmd())
	removeCmd.AddCommand(getRemoveBoardURLCmd())
	return removeCmd
}
