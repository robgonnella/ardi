package commands

import (
	"os"
	"strings"

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
	"ardi help",
}

var verbose bool
var quiet bool
var global bool
var dataDir = paths.ArdiProjectDataDir
var client *rpc.Client

func setLogger() {
	if verbose {
		logger.SetLevel(log.DebugLevel)
	} else if quiet {
		logger.SetLevel(log.FatalLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
}

func cmdIsProjectInit(cmd string) bool {
	return cmd == "ardi project init"
}

func cmdIsHelp(cmd string) bool {
	return strings.HasPrefix(cmd, "ardi help")
}

func cmdIsVersion(cmd string) bool {
	return cmd == "ardi version"
}

func shouldShowProjectError(cmd string) bool {
	return !global &&
		!util.IsProjectDirectory() &&
		!cmdIsProjectInit(cmd) &&
		!cmdIsHelp(cmd) &&
		!cmdIsVersion(cmd)
}

func getRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ardi",
		Short: "Ardi manages builds, uploads sketches and prints logs for a variety of arduino boards.",
		Long: "\nA light wrapper around arduino-cli that offers a quick way to manage builds, " +
			"upload sketches, and watch logs from command line for a variety of arduino boards.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogger()
			var err error
			cmdPath := cmd.CommandPath()

			if strings.HasPrefix(cmdPath, "ardi project") && global {
				logger.Error("Cannot specify --global with project command")
				os.Exit(1)
			}

			if shouldShowProjectError(cmdPath) {
				logger.Error("Not an ardi project directory")
				logger.Error("Try \"ardi project init\", or run with \"--global\"")
				os.Exit(1)
			}

			if global {
				dataDir = paths.ArdiGlobalDataDir
				confPath := paths.ArdiGlobalDataConfig
				util.InitDataDirectory(dataDir, confPath)
			}

			if !util.ArrayContains(noDaemon, cmdPath) {
				go rpc.StartDaemon(dataDir, verbose)
				if client, err = rpc.NewClient(logger); err != nil {
					os.Exit(1)
				}
			}

		},
	}
}

// GetRootCmd adds all ardi commands to root and returns root command
func GetRootCmd() *cobra.Command {
	rootCmd := getRootCommand()

	rootCmd.AddCommand(getVersionCommand())
	rootCmd.AddCommand(getCleanCommand())
	rootCmd.AddCommand(getGoCommand())
	rootCmd.AddCommand(getCompileCommand())
	rootCmd.AddCommand(getLibCommand())
	rootCmd.AddCommand(getPlatformCommand())
	rootCmd.AddCommand(getBoardCommand())
	rootCmd.AddCommand(getProjectCommand())
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Use global data directory")

	return rootCmd
}
