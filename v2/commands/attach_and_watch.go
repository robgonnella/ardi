package commands

import (
	"errors"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getWatchCmd() *cobra.Command {
	var fqbn string
	var buildProps []string
	var port string
	var baud int
	var watchCmd = &cobra.Command{
		Use:   "attach-and-watch [sketch|build]",
		Short: "Compile, upload, watch board logs, and watch for sketch changes",
		Long: "\nCompile, upload, watch board logs, and watch for sketch " +
			"changes. Updates to the sketch file will trigger automatic recompile, " +
			"reupload, and restarts the board log watcher. If the sketch argument " +
			"matches a user defined build in ardi.json, the build values will be " +
			"used for compilation, upload, and watch path",
		RunE: func(cmd *cobra.Command, args []string) error {
			optsList, err := ardiCore.GetCompileOptsFromArgs(fqbn, buildProps, false, args)
			if err != nil {
				return err
			}

			opts := optsList[0]

			if baud == 0 {
				baud = util.ParseSketchBaud(opts.SketchPath)
			}

			// Ignore errors here as user may have provided fqbn via build to mitigate
			// custom boards that don't show up via auto detect for some reason
			board, _ := ardiCore.Cli.GetTargetBoard(fqbn, port, true)

			if board == nil && opts.FQBN != "" && port != "" {
				board = &cli.BoardWithPort{FQBN: opts.FQBN, Port: port}
			}

			if board == nil {
				return errors.New("no connected boards detected")
			}

			opts.FQBN = board.FQBN

			if err := ardiCore.Compiler.Compile(*opts); err != nil {
				return err
			}

			if err := ardiCore.Uploader.Upload(board, opts.SketchDir); err != nil {
				return err
			}

			targets := core.WatchCoreTargets{
				Board:       board,
				CompileOpts: opts,
				Baud:        baud,
			}

			ardiCore.Watcher.SetTargets(targets)

			defer ardiCore.Watcher.Stop()
			return ardiCore.Watcher.Watch()
		},
	}

	watchCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	watchCmd.Flags().IntVarP(&baud, "baud", "b", 0, "Specify baud rate")
	watchCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	watchCmd.Flags().StringVar(&port, "port", "", "The port your arduino board is connected to")
	return watchCmd
}
