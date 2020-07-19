package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getUploadCmd() *cobra.Command {
	var watchBoardLogs bool
	var uploadCmd = &cobra.Command{
		Use:  "upload [sketch-dir|build]",
		Long: "\nUpload pre-compiled sketch build to a connected board",
		Short: "Upload pre-compiled sketch build to a connected board. If " +
			"the sketch argument matches a user defined build in ardi.json, the " +
			"build values will be used to find the appropraite build to upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			builds := ardiCore.Config.GetBuilds()

			build := "."
			if len(args) > 0 {
				build = args[0]
			}

			project := &types.Project{}
			var err error

			if ardiBuild, ok := builds[build]; ok {
				project.Directory = ardiBuild.Directory
				project.Sketch = ardiBuild.Sketch
				project.Baud = ardiBuild.Baud
			} else {
				project, err = util.ProcessSketch(build)
				if err != nil {
					return err
				}
			}

			connectedBoards := ardiCore.RPCClient.ConnectedBoards()
			allBoards := []*rpc.Board{}
			targetOpts := core.NewTargetOpts{
				ConnectedBoards: connectedBoards,
				AllBoards:       allBoards,
				OnlyConnected:   true,
				FQBN:            "",
				Logger:          logger,
			}
			target, err := core.NewTarget(targetOpts)
			if err != nil {
				return err
			}

			fields := logrus.Fields{
				"build":  project.Directory,
				"fqbn":   target.Board.FQBN,
				"device": target.Board.Port,
			}

			logger.WithFields(fields).Info("Uploading...")

			if err := ardiCore.Uploader.Upload(*target, project.Directory); err != nil {
				logger.WithError(err).Errorf("Failed to upload %s", project.Directory)
				return err
			}

			logger.Info("Upload successful")

			if watchBoardLogs {
				port := core.NewArdiSerialPort(target.Board.Port, project.Baud, logger)
				port.Watch()
			}

			return nil
		},
	}
	uploadCmd.Flags().BoolVarP(&watchBoardLogs, "log", "l", false, "Watch board logs after uploading")

	return uploadCmd
}
