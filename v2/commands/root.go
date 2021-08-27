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

// CommandEnv environment for all commands
type CommandEnv struct {
	Logger   *log.Logger
	Verbose  bool
	Quiet    bool
	ArdiCore *core.ArdiCore
	MockCli  cli.Cli
}

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

func setLogger(env *CommandEnv) {
	env.Logger.SetFormatter(&ardiLogFormatter{
		TextFormatter: log.TextFormatter{
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		},
	})

	if env.Verbose {
		env.Logger.SetLevel(log.DebugLevel)
		return
	}

	if env.Quiet {
		env.Logger.SetLevel(log.FatalLevel)
	} else {
		env.Logger.SetLevel(log.InfoLevel)
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

func getPreRun(env *CommandEnv) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger(env)

		cmdPath := cmd.CommandPath()

		if shouldShowProjectError(cmdPath) {
			return errors.New("not an ardi project directory, run 'ardi project-init' first")
		}

		return nil
	}
}

func getRootCommand(env *CommandEnv) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ardi",
		Short: "Ardi is a command line build manager for arduino projects.",
		Long: "\nArdi is a build tool that allows you to completely manage your arduino project from command line!\n\n" +
			"- Manage and store build configurations for projects with versioned dependencies\n- Run builds in CI Pipeline\n" +
			"- Compile & upload sketches to connected boards\n- Watch log output from connected boards in terminal\n" +
			"- Auto recompile / reupload on save",
		PersistentPreRunE: getPreRun(env),
		DisableAutoGenTag: true,
	}

	rootCmd.PersistentFlags().BoolVarP(&env.Verbose, "verbose", "v", false, "Print all logs")
	rootCmd.PersistentFlags().BoolVarP(&env.Quiet, "quiet", "q", false, "Silence all logs")
	rootCmd.SetHelpFunc(Help)
	return rootCmd
}

// GetRootCmd adds all ardi commands to root and returns root command
func GetRootCmd(env *CommandEnv) *cobra.Command {
	rootCmd := getRootCommand(env)
	rootCmd.AddCommand(getAddCmd(env))
	rootCmd.AddCommand(getCleanCmd(env))
	rootCmd.AddCommand(getCompileCmd(env))
	rootCmd.AddCommand(getCompileAndUploadCmd(env))
	rootCmd.AddCommand(getInstallCmd(env))
	rootCmd.AddCommand(getListCmd(env))
	rootCmd.AddCommand(getProjectInitCmd(env))
	rootCmd.AddCommand(getRemoveCmd(env))
	rootCmd.AddCommand(getSearchCmd(env))
	rootCmd.AddCommand(getUploadCmd(env))
	rootCmd.AddCommand(getVersionCmd(env))
	rootCmd.AddCommand(getWatchCmd(env))
	rootCmd.AddCommand(getAttachCmd(env))
	return rootCmd
}
