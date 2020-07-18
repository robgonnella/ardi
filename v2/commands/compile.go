package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCompileCmd() *cobra.Command {
	var fqbn string
	var buildProps []string
	var showProps bool
	var compileCmd = &cobra.Command{
		Use: "compile [sketch]",
		Long: "\nCompile sketches for a specified board. You must provide the " +
			"board FQBN, if left unspecified, a list of available choices will be " +
			"be printed.",
		Short: "Compile specified sketch",
		RunE: func(cmd *cobra.Command, args []string) error {
			sketchDir := "."
			if len(args) > 0 {
				sketchDir = args[0]
			}
			project, err := util.ProcessSketch(sketchDir)
			if err != nil {
				return err
			}

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

			fields := logrus.Fields{
				"sketch": project.Sketch,
				"baud":   project.Baud,
				"fqbn":   target.Board.FQBN,
				"device": target.Board.Port,
			}
			logger.WithFields(fields).Info("Compiling...")
			compileOpts := rpc.CompileOpts{
				FQBN:       target.Board.FQBN,
				SketchDir:  project.Directory,
				SketchPath: project.Sketch,
				ExportName: "",
				BuildProps: buildProps,
				ShowProps:  showProps,
			}
			if err := ardiCore.Compiler.Compile(compileOpts); err != nil {
				logger.WithError(err).Errorf("Failed to compile %s", sketchDir)
				return err
			}
			logger.WithFields(fields).Info("Compilation successful")
			return nil
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return compileCmd
}
