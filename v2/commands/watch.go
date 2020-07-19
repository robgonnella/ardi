package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func process(fqbn, sketch string, buildProps []string) error {
	builds := ardiCore.Config.GetBuilds()
	compileOpts := &rpc.CompileOpts{}
	var err error

	connectedBoards := ardiCore.RPCClient.ConnectedBoards()
	allBoards := ardiCore.RPCClient.AllBoards()
	targetOpts := core.NewTargetOpts{
		ConnectedBoards: connectedBoards,
		AllBoards:       allBoards,
		OnlyConnected:   false,
		FQBN:            fqbn,
		Logger:          logger,
	}

	target, err := core.NewTarget(targetOpts)
	if err != nil {
		return err
	}

	project := &types.Project{}

	if _, ok := builds[sketch]; ok {
		compileOpts, _ = ardiCore.Config.GetCompileOpts(sketch)
	} else {
		project, err = util.ProcessSketch(sketch)
		if err != nil {
			return err
		}

		compileOpts.FQBN = target.Board.FQBN
		compileOpts.SketchDir = project.Directory
		compileOpts.SketchPath = project.Sketch
		compileOpts.ExportName = ""
		compileOpts.BuildProps = buildProps
		compileOpts.ShowProps = false
	}

	if err := ardiCore.Compiler.Compile(*compileOpts); err != nil {
		logger.WithError(err).Error("Failed to compile")
		return err
	}

	if err := ardiCore.Uploader.Upload(*target, compileOpts.SketchDir); err != nil {
		logger.WithError(err).Error("Failed to upload")
		return err
	}

	return ardiCore.Watcher.Watch(*compileOpts, *target, project.Baud)
}

func getWatchCmd() *cobra.Command {
	var fqbn string
	var buildProps []string
	var watchCmd = &cobra.Command{
		Use:   "attach-and-watch [sketch|build]",
		Short: "Compile, upload, watch board logs, and watch for sketch changes",
		Long: "\nCompile, upload, watch board logs, and watch for sketch " +
			"changes. Updates to .ino file will trigger automatic recompile, " +
			"reupload, and restarts the board log watcher. If the sketch argument " +
			"matches a user defined build in ardi.json, the build values will be " +
			"used for compilation, upload, and watch path",
		RunE: func(cmd *cobra.Command, args []string) error {
			sketchDir := "."
			if len(args) > 0 {
				sketchDir = args[0]
			}

			return process(fqbn, sketchDir, buildProps)
		},
	}

	watchCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	watchCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")

	return watchCmd
}
