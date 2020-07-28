package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getWatchCmd() *cobra.Command {
	var fqbn string
	var buildProps []string
	var watchCmd = &cobra.Command{
		Use:   "attach-and-watch [sketch|build]",
		Short: "Compile, upload, watch board logs, and watch for sketch changes",
		Long: "\nCompile, upload, watch board logs, and watch for sketch " +
			"changes. Updates to the sketch file will trigger automatic recompile, " +
			"reupload, and restarts the board log watcher. If the sketch argument " +
			"matches a user defined build in ardi.json, the build values will be " +
			"used for compilation, upload, and watch path",
		RunE: func(cmd *cobra.Command, args []string) error {
			var compileOpts *rpc.CompileOpts
			var board *rpc.Board
			var baud int
			var err error

			builds := ardiCore.Config.GetBuilds()
			sketch := "."

			if len(args) > 0 {
				sketch = args[0]
			}

			if build, ok := builds[sketch]; ok {
				baud = util.ParseSketchBaud(build.Sketch)
				buildOpts := core.CompileArdiBuildOpts{
					BuildName:           sketch,
					OnlyConnectedBoards: true,
				}
				if compileOpts, board, err = ardiCore.CompileArdiBuild(buildOpts); err != nil {
					return err
				}
			} else {
				baud = util.ParseSketchBaud(sketch)
				sketchOpts := core.CompileSketchOpts{
					Sketch:              sketch,
					FQBN:                fqbn,
					BuildPros:           buildProps,
					ShowProps:           false,
					OnlyConnectedBoards: true,
				}
				if compileOpts, board, err = ardiCore.CompileSketch(sketchOpts); err != nil {
					return err
				}
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
			return ardiCore.Watcher.Watch()
		},
	}

	watchCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	watchCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")

	return watchCmd
}