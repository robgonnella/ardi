package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger *log.Logger
var port string
var client rpc.Client
var ardiCore *core.ArdiCore
var verbose bool
var quiet bool
var global bool
var dataDir = paths.ArdiProjectDataDir

func setLogger() {
	logger.Formatter = &log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	}
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

func preRun(cmd *cobra.Command, args []string) error {
	setLogger()
	cmdPath := cmd.CommandPath()

	if strings.HasPrefix(cmdPath, "ardi project") && global {
		logger.Error("Cannot specify --global with project command")
		return errors.New("Cannot specify --global with project command")
	}

	if shouldShowProjectError(cmdPath) {
		logger.Error("Not an ardi project directory")
		logger.Error("Try 'ardi project init', or run with '--global'")
		return errors.New("Not an ardi project directory")
	}

	if global || cmdPath == "ardi version" {
		dataDir = paths.ArdiGlobalDataDir
		confPath := paths.ArdiGlobalDataConfig
		util.InitDataDirectory(port, dataDir, confPath)
	}

	ctx := cmd.Context()
	client = rpc.NewClient(ctx, dataDir, port, logger)

	errChan := make(chan error, 1)
	successChan := make(chan string, 1)
	client.StartDaemon(verbose, successChan, errChan)
	select {
	case successMsg := <-successChan:
		logger.Debug(successMsg)
	case daemonErr := <-errChan:
		msg := fmt.Sprintf("arduino-cli daemon error: %s", daemonErr.Error())
		logger.Errorf(msg)
		return errors.New(msg)
	}

	if err := client.Connect(); err != nil {
		logger.WithError(err).Error("Failed to start ardi client")
		return err
	}

	if strings.Contains(cmdPath, "lib") || strings.Contains(cmdPath, "platform") {
		if err := client.UpdateIndexFiles(); err != nil {
			logger.WithError(err).Error("Failed to update index files")
		}
	}

	ardiCore = core.NewArdiCore(client, logger)

	if util.IsProjectDirectory() {
		if err := ardiCore.Project.SetConfigHelpers(); err != nil {
			logger.WithError(err).Error("Failed to initialize ardi project core")
			return err
		}
	}
	return nil
}

func getRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ardi",
		Short: "Ardi is a command line build manager for arduino projects.",
		Long: "\nArdi is a build tool that allows you to completely manage your arduino project from command line!\n\n" +
			"- Manage and store build configurations for projects with versioned dependencies\n- Run builds in CI Pipeline\n" +
			"- Compile & upload sketches to connected boards\n- Watch log output from connected boards in terminal\n" +
			"- Auto recompile / reupload on save",
		PersistentPreRunE: preRun,
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			client.Close()
			return nil
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Use global data directory")
	rootCmd.PersistentFlags().StringVar(&port, "port", "50051", "Set port for cli daemon")
	return rootCmd
}

// GetRootCmd adds all ardi commands to root and returns root command
func GetRootCmd(cmdLogger *log.Logger) *cobra.Command {
	logger = cmdLogger
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
