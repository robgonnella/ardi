package commands

import (
	"github.com/spf13/cobra"
)

func getLibSearchCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "search",
		Long:    "\nSearches for availables libraries with optional search filter",
		Short:   "Searches for availables libraries with optional search filter",
		Aliases: []string{"find"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ardiCore.Lib.Search(args[0])
		},
	}
	return initCmd
}

func getLibAddCommand() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Long:  "\nAdds specified libraries to either project or global library directory",
		Short: "Adds specified libraries to either project or global library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range args {
				ardiCore.Lib.Add(l)
			}
		},
	}
	return addCmd
}

func getLibRemoveCommand() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Long:  "\nRemoves specified libraries from project library directory",
		Short: "Removes specified libraries from project library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range args {
				ardiCore.Lib.Remove(l)
			}
		},
	}
	return removeCmd
}

func getLibCommand() *cobra.Command {
	var libCmd = &cobra.Command{
		Use:   "lib",
		Short: "Library manager",
		Long: "\nLibrary manager allowing you to add and remove libraries " +
			"either globally or at the project level. Each project can be " +
			"configured with its own list of dependencies for consistent " +
			"repeatable builds. See \"ardi help project\" form more " +
			"info on project level management with ardi",
	}
	libCmd.AddCommand(getLibAddCommand())
	libCmd.AddCommand(getLibRemoveCommand())
	libCmd.AddCommand(getLibSearchCommand())

	return libCmd
}
