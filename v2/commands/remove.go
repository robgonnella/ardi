package commands

import "github.com/spf13/cobra"

func getRemovePlatformCmd(env *CommandEnv) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nRemove platform(s) from project",
		Short:   "Remove platform(s) from project",
		Aliases: []string{"platform"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, p := range args {
				env.Logger.Infof("Removing platform: %s", p)
				removed, err := env.ArdiCore.Platform.Remove(p)
				if err != nil {
					return err
				}
				env.Logger.Infof("Removed %s", removed)
				if err := env.ArdiCore.Config.RemovePlatform(removed); err != nil {
					return err
				}
				env.Logger.Info("Udated config")
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveBuildCmd(env *CommandEnv) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "builds",
		Long:    "\nRemove build config from project",
		Short:   "Remove build config from project",
		Aliases: []string{"build"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, b := range args {
				if err := env.ArdiCore.Config.RemoveBuild(b); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveLibCmd(env *CommandEnv) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nRemove libraries from project",
		Short:   "Remove libraries from project",
		Aliases: []string{"libs", "lib", "library"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, l := range args {
				env.Logger.Infof("Removing library: %s", l)
				if err := env.ArdiCore.Lib.Remove(l); err != nil {
					return err
				}
				env.Logger.Infof("Removed %s", l)
				if err := env.ArdiCore.Config.RemoveLibrary(l); err != nil {
					return err
				}
				env.Logger.Info("Updated config")
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveBoardURLCmd(env *CommandEnv) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:     "board-urls",
		Long:    "\nRemove board urls from project",
		Short:   "Remove board urls from project",
		Aliases: []string{"board-url"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, url := range args {
				if err := env.ArdiCore.Config.RemoveBoardURL(url); err != nil {
					return err
				}
				if err := env.ArdiCore.CliConfig.RemoveBoardURL(url); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return removeCmd
}

func getRemoveCmd(env *CommandEnv) *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove project dependencies",
		Long:  "\nRemove project dependencies",
	}
	removeCmd.AddCommand(getRemovePlatformCmd(env))
	removeCmd.AddCommand(getRemoveBuildCmd(env))
	removeCmd.AddCommand(getRemoveLibCmd(env))
	removeCmd.AddCommand(getRemoveBoardURLCmd(env))
	return removeCmd
}
