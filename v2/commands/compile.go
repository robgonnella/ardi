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
	var watch bool
	var compileCmd = &cobra.Command{
		Use: "compile [sketch|build]",
		Long: "\nCompile sketches for a specified board. You must provide the " +
			"board FQBN, if left unspecified, a list of available choices will be " +
			"be printed. If the sketch argument matches as user defined build in " +
			"ardi.json, the values defined in build will be used to compile",
		Short: "Compile specified sketch",
		RunE: func(cmd *cobra.Command, args []string) error {
			sketch := "."

			if len(args) > 0 {
				sketch = args[0]
			}

			compileOpts := &rpc.CompileOpts{}
			builds := ardiCore.Config.GetBuilds()

			if _, ok := builds[sketch]; ok {
				compileOpts, _ = ardiCore.Config.GetCompileOpts(sketch)
				compileOpts.ShowProps = showProps
			} else {
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

				project, err := util.ProcessSketch(sketch)
				if err != nil {
					return err
				}

				compileOpts.FQBN = target.Board.FQBN
				compileOpts.SketchDir = project.Directory
				compileOpts.SketchPath = project.Sketch
				compileOpts.ExportName = ""
				compileOpts.BuildProps = buildProps
				compileOpts.ShowProps = showProps
			}

			fields := logrus.Fields{
				"sketch": compileOpts.SketchPath,
				"fqbn":   compileOpts.FQBN,
			}
			logger.WithFields(fields).Info("Compiling...")

			if err := ardiCore.Compiler.Compile(*compileOpts); err != nil {
				logger.WithError(err).Errorf("Failed to compile %s", sketch)
				return err
			}

			logger.WithFields(fields).Info("Compilation successful")

			if watch {
				return ardiCore.Compiler.WatchForChanges(*compileOpts)
			}

			return nil
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")
	compileCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch sketch file for changes and recompile")

	return compileCmd
}
