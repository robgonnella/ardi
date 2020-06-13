package commands

import (
	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	"github.com/robgonnella/ardi/v2/core/lib"

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
			libCore, err := lib.New(client, logger)
			if err != nil {
				return
			}

			libCore.Search(args[0])
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
			libCore, err := lib.New(client, logger)
			if err != nil {
				return
			}

			libCore.Add(args)
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
			libCore, err := lib.New(client, logger)
			if err != nil {
				return
			}

			libCore.Remove(args)
		},
	}
	return removeCmd
}

func getLibInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Long:  "\nInstalls all project level libraries specified in ardi.json",
		Short: "Installs all project level libraries specified in ardi.json",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			libCore, err := lib.New(client, logger)
			if err != nil {
				return
			}

			libCore.Install()
		},
	}
	return installCmd
}

func getLibListCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "list",
		Long:  "\nLists installed libraries specified in ardi.json",
		Short: "Lists installed libraries specified in ardi.json",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ardiJSON, err := ardijson.New(logger)
			if err != nil {
				return
			}
			ardiJSON.ListLibraries()
		},
	}
	return installCmd
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
	libCmd.AddCommand(getLibInstallCommand())
	libCmd.AddCommand(getLibSearchCommand())
	libCmd.AddCommand(getLibListCommand())

	return libCmd
}
