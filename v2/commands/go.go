package commands

import (
	ardiGoCore "github.com/robgonnella/ardi/v2/core/ardi-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func process(sketchDir string, buildProps []string, watchSketch, verbose bool) {

	logger := log.New()
	ardiGo, err := ardiGoCore.New(sketchDir, buildProps, logger)
	if err != nil {
		return
	}
	defer ardiGo.RPC.Connection.Close()

	if err := ardiGo.Compile(); err != nil {
		return
	}

	if err := ardiGo.Upload(); err != nil {
		return
	}

	if watchSketch {
		ardiGo.WatchSketch()
	} else {
		ardiGo.WatchLogs()
	}

}

func getGoCommand() *cobra.Command {

	var watchSketch bool
	var verbose bool
	var buildProps []string

	var goCmd = &cobra.Command{
		Use:   "go [sketch]",
		Short: "Compile and upload code to a connected arduino board",
		Long: "Compile and upload code to an arduino board. Simply pass the\n" +
			"directory containing the .ino file as the first argument. To watch\n" +
			"your sketch file for changes and auto re-compile & re-upload, use\n" +
			"the --watch flag. Baud will automatically be detected from sketch file.",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sketchDir := args[0]
			process(sketchDir, buildProps, watchSketch, verbose)
		},
	}
	goCmd.Flags().BoolVarP(&watchSketch, "watch", "w", false, "watch for changes, recompile and reupload")
	goCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print all logs")
	goCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")

	return goCmd
}
