package commands

import (
	"errors"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func getUploadCmd(env *CommandEnv) *cobra.Command {
	var baud int
	var attach bool
	var fqbn string
	var port string
	var uploadCmd = &cobra.Command{
		Use:   "upload [sketch-dir|build]",
		Short: "Upload pre-compiled sketch build to a connected board",
		Long: "\nUpload pre-compiled sketch build to a connected board. If " +
			"the sketch argument matches a user defined build in ardi.json, the " +
			"build values will be used to find the appropraite build to upload",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			defer env.ArdiCore.Uploader.Detach() // noop if not attached

			project := &types.Project{}
			var err error

			builds := env.ArdiCore.Config.GetBuilds()
			defaultBuild, defaultExits := builds["default"]

			if len(args) == 0 {
				if defaultExits {
					env.Logger.Info("Using build definition: default")
					project.Baud = defaultBuild.Baud
					project.Directory = defaultBuild.Directory
					project.Sketch = defaultBuild.Sketch
					if fqbn == "" {
						fqbn = defaultBuild.FQBN
					}
				} else if len(builds) == 1 {
					for name, b := range builds {
						env.Logger.Infof("Using build definition: %s", name)
						project.Baud = b.Baud
						project.Directory = b.Directory
						project.Sketch = b.Sketch
						if fqbn == "" {
							fqbn = b.FQBN
						}
					}
				} else {
					env.Logger.Info("Using ino in current directory")
					project, err = util.ProcessSketch(".")
					if err != nil {
						return err
					}
				}
			} else {
				if b, ok := builds[args[0]]; ok {
					env.Logger.Infof("Using build definition: %s", args[0])
					project.Baud = b.Baud
					project.Directory = b.Directory
					project.Sketch = b.Sketch
					if fqbn == "" {
						fqbn = b.FQBN
					}
				} else {
					env.Logger.Info("Using ino in current directory")
					project, err = util.ProcessSketch(args[0])
					if err != nil {
						return err
					}
				}
			}

			if baud != 0 {
				project.Baud = baud
			}

			// Ignore errors here as user may have provided fqbn via build to mitigate
			// custom boards that don't show up via auto detect for some reason
			board, _ := env.ArdiCore.Cli.GetTargetBoard(fqbn, port, true)

			if board == nil && fqbn != "" && port != "" {
				board = &cli.BoardWithPort{FQBN: fqbn, Port: port}
			}

			if board == nil {
				return errors.New("no connected boards detected")
			}

			if err := env.ArdiCore.Uploader.Upload(board, project.Directory); err != nil {
				return err
			}

			if attach {
				env.ArdiCore.Uploader.SetPortTargets(board.Port, project.Baud)
				env.ArdiCore.Uploader.Attach()
			}

			return nil
		},
	}
	uploadCmd.Flags().BoolVarP(&attach, "attach", "a", false, "Attach to board port and print logs")
	uploadCmd.Flags().IntVarP(&baud, "baud", "b", 0, "Specify baud rate when using \"attach\" flag")
	uploadCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "The FQBN of the board you want to upload to")
	uploadCmd.Flags().StringVarP(&port, "port", "p", "", "The port your arduino board is connected to")

	return uploadCmd
}
