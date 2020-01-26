package commands

import (
	"os"

	"github.com/robgonnella/ardi/ardi"
	"github.com/robgonnella/ardi/arguments"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func process(sketchDir string, baud int, buildProps []string, watchSketch, verbose bool) {
	if verbose {
		logger.SetLevel(log.DebugLevel)
		ardi.SetLogLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
		ardi.SetLogLevel(log.InfoLevel)
	}

	sketchDir, sketchFile, baud := arguments.ProcessSketch(sketchDir, baud)

	logFields := log.Fields{"baud": baud, "sketch": sketchDir}
	logWithFields := logger.WithFields(logFields)

	configFile := ardi.GlobalLibConfig
	if ardi.IsProjectDirectory() {
		configFile = ardi.LibConfig
	}

	conn, client, rpcInstance := ardi.StartDaemonAndGetConnection(configFile)
	defer conn.Close()

	list := ardi.GetTargetList(client, rpcInstance, sketchDir, sketchFile, baud)

	logWithFields.Debug("Parsing target")
	target := ardi.GetTargetInfo(list)
	target.BuildProps = buildProps

	logWithFields.WithField("target", target).Debug("Found target")

	ardi.Compile(client, rpcInstance, &target)
	if target.CompileError {
		os.Exit(1)
	}
	ardi.Upload(client, rpcInstance, &target)

	if watchSketch {
		ardi.WatchSketch(client, rpcInstance, &target)
	} else {
		ardi.WatchLogs(&target)
	}

}

func getGoCommand() *cobra.Command {

	var baud int
	var watchSketch bool
	var verbose bool
	var buildProps []string

	var goCmd = &cobra.Command{
		Use:   "go [sketch]",
		Short: "Compile and upload code to a connected arduino board",
		Long: "Compile and upload code to an arduino board. Simply pass the\n" +
			"directory containing the .ino file as the first argument. To watch\n" +
			"your sketch file for changes and auto re-compile & re-upload, use\n" +
			"the --watch flag. You can also specify the baud rate with --baud\n" +
			"(default is 9600).",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sketchDir := args[0]
			process(sketchDir, baud, buildProps, watchSketch, verbose)
		},
	}
	goCmd.Flags().IntVarP(&baud, "baud", "b", 9600, "specify sketch baud rate")
	goCmd.Flags().BoolVarP(&watchSketch, "watch", "w", false, "watch for changes, recompile and reupload")
	goCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print all logs")
	goCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")

	return goCmd
}
