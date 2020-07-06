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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ardiCore.Lib.Search(args[0]); err != nil {
				logger.WithError(err).Error("Failed to find arduino libraries")
				return err
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, l := range args {
				if _, _, err := ardiCore.Lib.Add(l); err != nil {
					logger.WithError(err).Errorf("Failed to add library %s", l)
					return err
				}
			}
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, l := range args {
				if err := ardiCore.Lib.Remove(l); err != nil {
					logger.WithError(err).Errorf("Failed to remove library %s", l)
					return err
				}
			}
			return nil
		},
	}
	return removeCmd
}

func getLibListCommand() *cobra.Command {
	listCmd := &cobra.Command{
		Use: "list",
		Long: "\nLists all installed libraries. Use \"--global\" to list " +
			"globally installed libraries",
		Short: "Lists all installed libraries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ardiCore.Lib.ListInstalled(); err != nil {
				logger.WithError(err).Error("Failed to list installed libraries")
				return err
			}
			return nil
		},
	}
	return listCmd
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
	libCmd.AddCommand(getLibListCommand())

	return libCmd
}
