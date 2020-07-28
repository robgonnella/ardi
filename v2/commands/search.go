package commands

import "github.com/spf13/cobra"

func getSearchPlatformCmd() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nSearch all available platforms",
		Short:   "Search all available platforms",
		Aliases: []string{"platform"},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Available platforms")
			if err := ardiCore.Platform.ListAll(); err != nil {
				return err
			}
			return nil
		},
	}
	return searchCmd
}

func getSearchLibCmd() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nSearches for availables libraries with optional search filter",
		Short:   "Searches for availables libraries with optional search filter",
		Aliases: []string{"lib", "libs", "library"},
		RunE: func(cmd *cobra.Command, args []string) error {
			searchArg := ""
			if len(args) > 0 {
				searchArg = args[0]
			}
			if err := ardiCore.Lib.Search(searchArg); err != nil {
				return err
			}
			return nil
		},
	}
	return searchCmd
}

func getSearchCmd() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for arduino platforms, libraries, and boards",
		Long:  "\nSearch for arduino platforms, libraries, and boards",
	}
	searchCmd.AddCommand(getSearchPlatformCmd())
	searchCmd.AddCommand(getSearchLibCmd())
	return searchCmd
}
