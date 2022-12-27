package commands

import "github.com/spf13/cobra"

func newListPlatformCmd(env *CommandEnv) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nList project platforms",
		Short:   "List project platforms",
		Aliases: []string{"platform"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			env.Logger.Info("Platforms specified in ardi.json")
			env.ArdiCore.Config.ListPlatforms()

			env.Logger.Info("Installed platforms")
			if err := env.ArdiCore.Platform.ListInstalled(); err != nil {
				return err
			}

			return nil
		},
	}
	return listCmd
}

func newListLibrariesCmd(env *CommandEnv) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nList project libraries",
		Short:   "List project libraries",
		Aliases: []string{"libs"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			env.Logger.Info("Libraries specified in ardi.json")
			env.ArdiCore.Config.ListLibraries()
			env.Logger.Info("Installed libraries")
			if err := env.ArdiCore.Lib.ListInstalled(); err != nil {
				return err
			}
			return nil
		},
	}
	return listCmd
}

func newListBuildsCmd(env *CommandEnv) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "builds",
		Long:    "\nList project builds",
		Short:   "List project builds",
		Aliases: []string{"build"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			env.ArdiCore.Config.ListBuilds(args)
			return nil
		},
	}
	return listCmd
}

func newListBoardURLSCmd(env *CommandEnv) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "board-urls",
		Long:  "\nList project board urls",
		Short: "List project board urls",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			env.ArdiCore.Config.ListBoardURLS()
			return nil
		},
	}
	return listCmd
}

func newListCmd(env *CommandEnv) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nList platforms, libraries, board urls, and builds",
		Short: "List platforms, libraries, board urls, and builds",
	}
	listCmd.AddCommand(newListPlatformCmd(env))
	listCmd.AddCommand(newListLibrariesCmd(env))
	listCmd.AddCommand(newListBuildsCmd(env))
	listCmd.AddCommand(newListBoardURLSCmd(env))
	return listCmd
}
