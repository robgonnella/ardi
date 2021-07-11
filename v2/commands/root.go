package commands

import (
	"errors"
	"io/ioutil"
	"strings"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger *log.Logger
var verbose bool
var quiet bool
var ardiCore *core.ArdiCore
var cliInstance cli.Cli

type ardiLogFormatter struct {
	log.TextFormatter
}

func (a *ardiLogFormatter) Format(e *log.Entry) ([]byte, error) {
	b, err := a.TextFormatter.Format(e)
	if err != nil {
		return b, err
	}
	str := string(b)
	str = strings.Replace(str, strings.ToUpper(e.Level.String()), "ardi", 1)
	return []byte(str), nil
}

func setLogger() {
	logger.SetFormatter(&ardiLogFormatter{
		TextFormatter: log.TextFormatter{
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		},
	})

	if verbose {
		logger.SetLevel(log.DebugLevel)
		return
	}

	if quiet {
		logger.SetLevel(log.FatalLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	log.SetOutput(ioutil.Discard)
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

func cmdIsArdiAttach(cmd string) bool {
	return cmd == "ardi attach"
}

func shouldShowProjectError(cmd string) bool {
	return !util.IsProjectDirectory() &&
		!cmdIsProjectInit(cmd) &&
		!cmdIsArdiAttach(cmd) &&
		!cmdIsHelp(cmd) &&
		!cmdIsVersion(cmd)
}

func preRun(cmd *cobra.Command, args []string) error {
	setLogger()

	cmdPath := cmd.CommandPath()

	if shouldShowProjectError(cmdPath) {
		return errors.New("not an ardi project directory, run 'ardi project-init' first")
	}

	ardiConfig, svrSettings := util.GetAllSettings()
	cliSettingsPath := util.GetCliSettingsPath()

	writeOpts := util.WriteSettingsOpts{
		ArdiConfig:         ardiConfig,
		ArduinoCliSettings: svrSettings,
	}
	if util.IsProjectDirectory() {
		if err := util.WriteAllSettings(writeOpts); err != nil {
			return err
		}
	}

	ctx := cmd.Context()
	cliWrapper := cli.NewCli(ctx, cliSettingsPath, svrSettings, logger, cliInstance)

	coreOpts := core.NewArdiCoreOpts{
		Logger:             logger,
		Cli:                cliWrapper,
		ArdiConfig:         *ardiConfig,
		ArduinoCliSettings: *svrSettings,
	}
	ardiCore = core.NewArdiCore(coreOpts)

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
		DisableAutoGenTag: true,
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.SetHelpFunc(help)
	return rootCmd
}

// GetRootCmd adds all ardi commands to root and returns root command
func GetRootCmd(cmdLogger *log.Logger, instance cli.Cli) *cobra.Command {
	logger = cmdLogger
	cliInstance = instance
	rootCmd := getRootCommand()
	rootCmd.AddCommand(getAddCmd())
	rootCmd.AddCommand(getCleanCmd())
	rootCmd.AddCommand(getCompileCmd())
	rootCmd.AddCommand(getInstallCmd())
	rootCmd.AddCommand(getListCmd())
	rootCmd.AddCommand(getProjectInitCmd())
	rootCmd.AddCommand(getRemoveCmd())
	rootCmd.AddCommand(getSearchCmd())
	rootCmd.AddCommand(getUploadCmd())
	rootCmd.AddCommand(getVersionCmd())
	rootCmd.AddCommand(getWatchCmd())
	rootCmd.AddCommand(getAttachCmd())
	return rootCmd
}
