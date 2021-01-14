package commands

import "github.com/spf13/cobra"

func getListPlatformCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nList project platforms",
		Short:   "List project platforms",
		Aliases: []string{"platform"},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Platforms specified in ardi.json")
			ardiCore.Config.ListPlatforms()

			logger.Info("Installed platforms")
			if err := ardiCore.Platform.ListInstalled(); err != nil {
				return err
			}

			return nil
		},
	}
	return withRPCConnectPreRun(listCmd)
}

func getListLibrariesCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nList project libraries",
		Short:   "List project libraries",
		Aliases: []string{"libs"},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Libraries specified in ardi.json")
			ardiCore.Config.ListLibraries()
			logger.Info("Installed libraries")
			if err := ardiCore.Lib.ListInstalled(); err != nil {
				return err
			}
			return nil
		},
	}
	return withRPCConnectPreRun(listCmd)
}

func getListBuildsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "builds",
		Long:    "\nList project builds",
		Short:   "List project builds",
		Aliases: []string{"build"},
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Config.ListBuilds(args)
		},
	}
	return listCmd
}

func getListBoardURLSCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "board-urls",
		Long:  "\nList project board urls",
		Short: "List project board urls",
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Config.ListBoardURLS()
		},
	}
	return listCmd
}

func getListBoardFQBNSCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "board-fqbns",
		Long:  "\nList boards with associated fqbns",
		Short: "List boards with associated fqbns",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			return ardiCore.Board.FQBNS(query)
		},
	}
	return withRPCConnectPreRun(listCmd)
}

func getListBoardPlatformsCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "board-platforms",
		Long:  "\nList boards with their associated platform",
		Short: "List boards with their associated platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			return ardiCore.Board.Platforms(query)
		},
	}
	return withRPCConnectPreRun(listCmd)
}

func getListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nList platforms, libraries, board urls, and builds",
		Short: "List platforms, libraries, board urls, and builds",
	}
	listCmd.AddCommand(getListPlatformCmd())
	listCmd.AddCommand(getListLibrariesCmd())
	listCmd.AddCommand(getListBuildsCmd())
	listCmd.AddCommand(getListBoardURLSCmd())
	listCmd.AddCommand(getListBoardFQBNSCmd())
	listCmd.AddCommand(getListBoardPlatformsCmd())
	return listCmd
}
