package commands

import (
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getUploadCmd() *cobra.Command {
	var compileCmd = &cobra.Command{
		Use:   "upload [build-dir]",
		Long:  "\nUpload pre-compiled sketch build to a connected board",
		Short: "Upload pre-compiled sketch build to a connected board",
		RunE: func(cmd *cobra.Command, args []string) error {
			builds := ardiCore.Config.GetBuilds()

			build := "."
			if len(args) > 0 {
				build = args[0]
			}

			if ardiBuild, ok := builds[build]; ok {
				build = ardiBuild.Path
			}

			project, err := util.ProcessSketch(build)
			if err != nil {
				return err
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
			return nil
		},
	}
	return compileCmd
}
