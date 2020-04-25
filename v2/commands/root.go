package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var logger = log.New()

func getRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ardi",
		Short: "Ardi uploads sketches and prints logs for a variety of arduino boards.",
		Long: cyan("\nA light wrapper around arduino-cli that offers a quick way to upload\n" +
			"sketches and watch logs from command line for a variety of arduino boards."),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "init" {
				if err := initializeDataDirectory(); err != nil {
					logger.WithError(err).Error("Failed to initialize data directory")
					return
				}
				if err := initializeArdiJSON(); err != nil {
					logger.WithError(err).Error("Failed to initialize ardi.json")
					return
				}
			}
			dataConfig := paths.ArdiDataDir
			go rpc.StartDaemon(dataConfig)
		},
	}
}

// Initialize adds all ardi commands to root and returns root command
func Initialize(version string) *cobra.Command {
	rootCmd := getRootCommand()

	rootCmd.AddCommand(getVersionCommand(version))
	rootCmd.AddCommand(getCleanCommand())
	rootCmd.AddCommand(getGoCommand())
	rootCmd.AddCommand(getCompileCommand())
	rootCmd.AddCommand(getLibCommand())
	rootCmd.AddCommand(getPlatformCommand())
	rootCmd.AddCommand(getBoardCommand())
	rootCmd.AddCommand(getProjectCommand())

	return rootCmd
}

// private helpers
func initializeDataDirectory() error {
	if _, err := os.Stat(paths.ArdiDataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(paths.ArdiDataDir, 0777); err != nil {
			return err
		}
	}

	if _, err := os.Stat(paths.ArdiDataConfig); os.IsNotExist(err) {
		dataConfig := types.DataConfig{
			BoardManager: types.BoardManager{AdditionalUrls: []string{}},
			Directories: types.Directories{
				Data:      paths.ArdiDataDir,
				Downloads: path.Join(paths.ArdiDataDir, "staging"),
				User:      path.Join(paths.ArdiDataDir, "Arduino"),
			},
			Telemetry: types.Telemetry{Enabled: false},
		}
		yamlConfig, _ := yaml.Marshal(&dataConfig)
		if err := ioutil.WriteFile(paths.ArdiDataConfig, yamlConfig, 0644); err != nil {
			return err
		}
	}

	return nil
}

func initializeArdiJSON() error {
	if _, err := os.Stat(paths.ArdiBuildConfig); os.IsNotExist(err) {
		buildConfig := types.ArdiConfig{
			Libraries: make(map[string]string),
			Builds:    make(map[string]types.ArdiBuildJSON),
		}
		jsonConfig, _ := json.MarshalIndent(&buildConfig, "\n", " ")
		if err := ioutil.WriteFile(paths.ArdiBuildConfig, jsonConfig, 0644); err != nil {
			return err
		}
		logger.Info("ardi.json initialized")
	}
	return nil
}
