package commands

import (
	"os"
	"strings"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = log.New()
var port string
var client rpc.Client
var ardiCore *core.ArdiCore
var verbose bool
var quiet bool
var global bool
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
	rootCmd := &cobra.Command{
		Use:   "ardi",
		Short: "Ardi is a command line build manager for arduino projects.",
		Long: "\nArdi is a build tool that allows you to completely manage your arduino project from command line!\n\n" +
			"- Manage and store build configurations for projects with versioned dependencies\n- Run builds in CI Pipeline\n" +
			"- Compile & upload sketches to connected boards\n- Watch log output from connected boards in terminal\n" +
			"- Auto recompile / reupload on save",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogger()
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
				util.InitDataDirectory(port, dataDir, confPath)
			}

			go rpc.StartDaemon(port, dataDir, verbose)

			var err error
			client, err = rpc.NewClient(port, logger)
			if err != nil {
				logger.WithError(err).Error("Failed to start ardi client")
				os.Exit(1)
			}

			if strings.Contains(cmdPath, "lib") || strings.Contains(cmdPath, "platform") {
				if err := client.UpdateIndexFiles(); err != nil {
					logger.WithError(err).Error("Failed to update index files")
				}
			}

			ardiCore = core.NewArdiCore(client, logger)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Use global data directory")
	rootCmd.PersistentFlags().StringVar(&port, "port", "50051", "Set port for cli daemon")
	return rootCmd
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

	return rootCmd
}
