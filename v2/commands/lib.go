package commands

import (
	"strings"

	"github.com/robgonnella/ardi/v2/core/lib"
	"github.com/robgonnella/ardi/v2/paths"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func getLibInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes current directory with library dependency config file",
		Long: "Initializes current directory with library dependency config\n" +
			"file. Each project directory can then specify its own separate\n" +
			"library versions.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.New()
			logger.Info("Initializing library manager")
			libCore, err := lib.New(paths.ArdiDataConfig, logger)
			if err != nil {
				return
			}
			defer libCore.RPC.Connection.Close()
			libCore.Init()
		},
	}
	return initCmd
}

func getLibSearchCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:     "search",
		Short:   "Searches for availables libraries with optional search filter",
		Aliases: []string{"find", "list"},
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
			libParts := strings.Split(args[0], "@")
			library := libParts[0]
			version := ""
			if len(libParts) > 1 {
				version = libParts[1]
			}
			libCore, err := lib.New(paths.ArdiDataConfig, logger)
			if err != nil {
				return
			}
			defer libCore.RPC.Connection.Close()
			libCore.Add(library, version)
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
			libCore.Remove(args[0])
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

func getLibCommand() *cobra.Command {
	var libCmd = &cobra.Command{
		Use:   "lib",
		Short: "Library manager for ardi",
		Long: "Library manager for ardi allowing you to add and remove libraries\n" +
			"either globally or at the project level. Each project can be\n" +
			"configured with its own list of dependencies for consistent\n" +
			"repeatable builds every time.",
	}
	libCmd.AddCommand(getLibInitCommand())
	libCmd.AddCommand(getLibAddCommand())
	libCmd.AddCommand(getLibRemoveCommand())
	libCmd.AddCommand(getLibInstallCommand())
	libCmd.AddCommand(getLibSearchCommand())

	return libCmd
}
