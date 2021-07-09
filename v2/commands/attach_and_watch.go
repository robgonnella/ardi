package commands

import (
	"os"
	"os/signal"
	"syscall"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getWatchCmd() *cobra.Command {
	var fqbn string
	var buildProps []string
	var port string
	var watchCmd = &cobra.Command{
		Use:   "attach-and-watch [sketch|build]",
		Short: "Compile, upload, watch board logs, and watch for sketch changes",
		Long: "\nCompile, upload, watch board logs, and watch for sketch " +
			"changes. Updates to the sketch file will trigger automatic recompile, " +
			"reupload, and restarts the board log watcher. If the sketch argument " +
			"matches a user defined build in ardi.json, the build values will be " +
			"used for compilation, upload, and watch path",
		RunE: func(cmd *cobra.Command, args []string) error {
			var compileOpts *cli.CompileOpts
			var baud int
			var err error

			builds := ardiCore.Config.GetBuilds()
			sketch := "."

			if len(args) > 0 {
				sketch = args[0]
			}

			if build, ok := builds[sketch]; ok {
				baud = util.ParseSketchBaud(build.Sketch)
				if compileOpts, err = ardiCore.CompileArdiBuild(sketch); err != nil {
					return err
				}
			} else {
				baud = util.ParseSketchBaud(sketch)
				sketchOpts := core.CompileSketchOpts{
					Sketch:    sketch,
					FQBN:      fqbn,
					BuildPros: buildProps,
					ShowProps: false,
				}
				if compileOpts, err = ardiCore.CompileSketch(sketchOpts); err != nil {
					return err
				}
			}

			board, err := ardiCore.Cli.GetTargetBoard(compileOpts.FQBN, port, true)
			if err != nil || board == nil {
				return err
			}

			if err := ardiCore.Uploader.Upload(board, compileOpts.SketchDir); err != nil {
				return err
			}

			targets := core.WatchCoreTargets{
				Board:       board,
				CompileOpts: compileOpts,
				Baud:        baud,
			}
			ardiCore.Watcher.SetTargets(targets)

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigs
				logger.Debug("gracefully shutting down attach-and-watch")
				ardiCore.Watcher.Stop()
			}()

			return ardiCore.Watcher.Watch()
		},
	}

	watchCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	watchCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	watchCmd.Flags().StringVar(&port, "port", "", "The port your arduino board is connected to")
	return watchCmd
}
