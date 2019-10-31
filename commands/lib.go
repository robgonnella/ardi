package commands

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/robgonnella/ardi/v3/ardi"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func getLibInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes current directory with library dependency config file",
		Long: "Initializes current directory with library dependency config\n" +
			"file. Each project directory can then specify its own separate\n" +
			"library versions.",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Initializing library manager")
			dirConfig := ardi.LibraryDirConfig{
				ProxyType:      "auto",
				SketchbookPath: ".",
				ArduinoData:    ".",
				BoardManager:   make(map[string]interface{}),
			}
			yamlConfig, _ := yaml.Marshal(&dirConfig)
			ioutil.WriteFile(ardi.LibConfig, yamlConfig, 0644)

			depConfig := make(map[string]interface{})
			jsonConfig, _ := json.MarshalIndent(&depConfig, "\n", " ")
			ioutil.WriteFile(ardi.DepConfig, jsonConfig, 0644)
			logger.Info("Directory initialized")
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
			configFile := ardi.GlobalLibConfig
			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()
			ardi.LibSearch(client, instance, args[0])
		},
	}
	return initCmd
}

func getLibAddCommand() *cobra.Command {
	var global bool
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Adds specified libraries to either project or global library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !ardi.IsProjectDirectory() && !global {
				logger.Fatal("Directory not initialized. Run \"ardi lib init\" or run with --global")
			}
			libParts := strings.Split(args[0], "@")
			lib := libParts[0]
			version := ""
			if len(libParts) > 1 {
				version = libParts[1]
			}
			logger.Infof("Adding library: %s %s", lib, version)
			configFile := ardi.LibConfig
			if global {
				configFile = ardi.GlobalLibConfig
			}
			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()
			installedVersion := ardi.LibInstall(client, instance, lib, version)
			if !global {
				config := make(map[string]interface{})
				configData, _ := ioutil.ReadFile(ardi.DepConfig)
				json.Unmarshal(configData, &config)
				config[lib] = installedVersion
				newData, _ := json.MarshalIndent(config, "", " ")
				ioutil.WriteFile(ardi.DepConfig, newData, 0644)
			}
		},
	}
	addCmd.Flags().BoolVarP(&global, "global", "g", false, "Instructs ardi to install library globally")
	return addCmd
}

func getLibRemoveCommand() *cobra.Command {
	var global bool
	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Removes specified libraries from either project or global library directory",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if !ardi.IsProjectDirectory() && !global {
				logger.Fatal("Directory not initialized. Run \"ardi lib init\" or run with --global")
			}

			lib := args[0]
			logger.Infof("Removing library: %s", lib)
			configFile := ardi.LibConfig
			if global {
				configFile = ardi.GlobalLibConfig
			}
			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()
			ardi.LibUnInstall(client, instance, lib)
			if !global {
				config := make(map[string]interface{})
				newConfig := make(map[string]interface{})
				configData, _ := ioutil.ReadFile(ardi.DepConfig)
				json.Unmarshal(configData, &config)
				for name, version := range config {
					if name != lib {
						newConfig[name] = version
					}
				}
				newData, _ := json.MarshalIndent(newConfig, "\n", " ")
				ioutil.WriteFile(ardi.DepConfig, newData, 0644)
			}
		},
	}
	removeCmd.Flags().BoolVarP(&global, "global", "g", false, "Instructs ardi to uninstall library globally")
	return removeCmd
}

func getLibInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Installs all project level libraries specified in ardi.json",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !ardi.IsProjectDirectory() {
				logger.Fatal("Directory not initialized. Run \"ardi lib init\" or run with --global")
			}

			logger.Info("Installing libraries")
			configFile := ardi.LibConfig
			conn, client, instance := ardi.StartDaemonAndGetConnection(configFile)
			defer conn.Close()
			config := make(map[string]string)
			configData, err := ioutil.ReadFile(ardi.DepConfig)
			if err != nil {
				logger.WithError(err).Fatal("Could not install libraries. You may need to run \"ardi lib init\" in this directory")
			}
			err = json.Unmarshal(configData, &config)
			if err != nil {
				logger.WithError(err).Fatal("Could not install libraries. Potentially malformed json file")
			}
			for name, version := range config {
				ardi.LibInstall(client, instance, name, version)
			}
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
