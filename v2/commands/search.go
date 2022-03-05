package commands

import "github.com/spf13/cobra"

func getSearchPlatformCmd(env *CommandEnv) *cobra.Command {
	searchCmd := &cobra.Command{
		Use:     "platforms",
		Long:    "\nSearch all available platforms",
		Short:   "Search all available platforms",
		Aliases: []string{"platform"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			env.Logger.Info("Available platforms")
			if err := env.ArdiCore.Platform.ListAll(); err != nil {
				return err
			}
			return nil
		},
	}
	return searchCmd
}

func getSearchLibCmd(env *CommandEnv) *cobra.Command {
	searchCmd := &cobra.Command{
		Use:     "libraries",
		Long:    "\nSearches for availables libraries with optional search filter",
		Short:   "Searches for availables libraries with optional search filter",
		Aliases: []string{"lib", "libs", "library"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireProjectInit(); err != nil {
				return err
			}
			searchArg := ""
			if len(args) > 0 {
				searchArg = args[0]
			}
			if err := env.ArdiCore.Lib.Search(searchArg); err != nil {
				return err
			}
			return nil
		},
	}
	return searchCmd
}

func getSearchCmd(env *CommandEnv) *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for arduino platforms, libraries, and boards",
		Long:  "\nSearch for arduino platforms, libraries, and boards",
	}
	searchCmd.AddCommand(getSearchPlatformCmd(env))
	searchCmd.AddCommand(getSearchLibCmd(env))
	return searchCmd
}
