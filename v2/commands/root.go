package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger *log.Logger
var port string
var verbose bool
var quiet bool
var global bool
var ardiCore *core.ArdiCore

func setLogger() {
	logger.Formatter = &log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	}

	if verbose {
		logger.SetLevel(log.DebugLevel)
		return
	}

	if quiet {
		logger.SetLevel(log.FatalLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	logrus.SetOutput(ioutil.Discard)
}

func cmdIsProjectInit(cmd string) bool {
	return cmd == "ardi project-init"
}

func cmdIsHelp(cmd string) bool {
	return strings.Contains(cmd, "help")
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

func useGlobalData(cmd string) bool {
	return global || cmdIsVersion(cmd) || cmdIsHelp(cmd)
}

func preRun(cmd *cobra.Command, args []string) error {
	setLogger()

	cmdPath := cmd.CommandPath()
	useGlobal := useGlobalData(cmdPath)
	daemonLogLevel := util.GetDaemonLogLevel(logger)

	if shouldShowProjectError(cmdPath) {
		logger.Error("Not an ardi project directory")
		logger.Error("Try 'ardi project-init', or run with '--global'")
		return errors.New("Not an ardi project directory")
	}

	getOpts := util.GetAllSettingsOpts{
		Global:   useGlobal,
		LogLevel: daemonLogLevel,
		Port:     port,
	}
	ardiConfig, svrSettings := util.GetAllSettings(getOpts)

	writeOpts := util.WriteSettingsOpts{
		Global:             useGlobal,
		ArdiConfig:         ardiConfig,
		ArduinoCliSettings: svrSettings,
	}
	if useGlobal || util.IsProjectDirectory() {
		if err := util.WriteAllSettings(writeOpts); err != nil {
			return err
		}
	}

	ctx := cmd.Context()
	client := rpc.NewClient(ctx, svrSettings, logger)
	coreOpts := core.NewArdiCoreOpts{
		Global:             useGlobal,
		Logger:             logger,
		Client:             client,
		ArdiConfig:         *ardiConfig,
		ArduinoCliSettings: *svrSettings,
	}
	ardiCore = core.NewArdiCore(coreOpts)

	errChan := make(chan error, 1)
	successChan := make(chan string, 1)
	ardiCore.RPCClient.StartDaemon(successChan, errChan)
	select {
	case successMsg := <-successChan:
		logger.Debug(successMsg)
	case daemonErr := <-errChan:
		msg := fmt.Sprintf("arduino-cli daemon error: %s", daemonErr.Error())
		logger.Errorf(msg)
		return errors.New(msg)
	}

	if err := client.Connect(); err != nil {
		logger.WithError(err).Error("Failed to connect ardi client")
		return err
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
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			ardiCore.RPCClient.Close()
		},
		DisableAutoGenTag: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Use global data directory")
	rootCmd.PersistentFlags().StringVar(&port, "port", "", "Set port for cli daemon")
	return rootCmd
}

// GetRootCmd adds all ardi commands to root and returns root command
func GetRootCmd(cmdLogger *log.Logger) *cobra.Command {
	logger = cmdLogger
	rootCmd := getRootCommand()
	rootCmd.AddCommand(getAddCmd())
	rootCmd.AddCommand(getBuildCmd())
	rootCmd.AddCommand(getCleanCmd())
	rootCmd.AddCommand(getCompileCmd())
	rootCmd.AddCommand(getInstallCmd())
	rootCmd.AddCommand(getListCmd())
	rootCmd.AddCommand(getProjectInitCmd())
	rootCmd.AddCommand(getRemoveCmd())
	rootCmd.AddCommand(getSearchCmd())
	rootCmd.AddCommand(getVersionCmd())
	rootCmd.AddCommand(getWatchCmd())

	return rootCmd
}
