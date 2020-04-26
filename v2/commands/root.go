package commands

import (
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = log.New()

var noDaemon = []string{
	"ardi version",
	"ardi clean",
	"ardi project init",
}

var noProjectCheck = []string{
	"ardi project init",
	"ardi version",
}

var verbose bool
var quiet bool
var dataDir = paths.ArdiProjectDataDir

func setLogger() {
	if verbose {
		logger.SetLevel(log.DebugLevel)
	} else if quiet {
		logger.SetLevel(log.FatalLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
}

func getRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ardi",
		Short: "Ardi uploads sketches and prints logs for a variety of arduino boards.",
		Long: cyan("\nA light wrapper around arduino-cli that offers a quick way to upload\n" +
			"sketches and watch logs from command line for a variety of arduino boards."),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogger()
			cmdPath := cmd.CommandPath()
			dataDir := paths.ArdiProjectDataDir

			if !util.IsProjectDirectory() {
				dataDir = paths.ArdiGlobalDataDir
				confPath := paths.ArdiGlobalDataConfig
				util.InitDataDirectory(dataDir, confPath)
			}

			if !util.ArrayContains(noDaemon, cmdPath) {
				go rpc.StartDaemon(dataDir)
			}
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")

	return rootCmd
}
