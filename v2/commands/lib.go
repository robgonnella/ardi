package commands

import (
	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	"github.com/robgonnella/ardi/v2/core/lib"
	"github.com/robgonnella/ardi/v2/paths"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func getLibSearchCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "search",
		Short:   "Searches for availables libraries with optional search filter",
		Aliases: []string{"find"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(paths.ArdiGlobalDataConfig, logger)
			if err != nil {
				return
			}
			defer libCore.RPC.Connection.Close()
			libCore.Search(args[0])
		},
	}
	return initCmd
}

func getLibAddCommand() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds specified libraries to either project or global library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(paths.ArdiDataConfig, logger)
			if err != nil {
				return
			}
			defer libCore.RPC.Connection.Close()
			libCore.Add(args)
		},
	}
	return addCmd
}

func getLibRemoveCommand() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Removes specified libraries from project library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(paths.ArdiDataConfig, logger)
			if err != nil {
				return
			}
			defer libCore.RPC.Connection.Close()
			libCore.Remove(args)
		},
	}
	return removeCmd
}

func getLibInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Installs all project level libraries specified in ardi.json",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			libCore, err := lib.New(paths.ArdiDataConfig, logger)
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
		Short: "Lists installed libraries specified in ardi.json",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
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
		Short: "Library manager for ardi",
		Long: "Library manager for ardi allowing you to add and remove libraries\n" +
			"either globally or at the project level. Each project can be\n" +
			"configured with its own list of dependencies for consistent\n" +
			"repeatable builds every time.",
	}
	libCmd.AddCommand(getLibAddCommand())
	libCmd.AddCommand(getLibRemoveCommand())
	libCmd.AddCommand(getLibInstallCommand())
	libCmd.AddCommand(getLibSearchCommand())
	libCmd.AddCommand(getLibListCommand())

	return libCmd
}
