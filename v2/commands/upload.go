package commands

import (
	"path/filepath"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/spf13/cobra"
)

func getUploadCmd() *cobra.Command {
	var compileCmd = &cobra.Command{
		Use:   "upload [build-dir]",
		Long:  "\nUpload pre-compiled sketch build to a connected board",
		Short: "Upload pre-compiled sketch build to a connected board",
		RunE: func(cmd *cobra.Command, args []string) error {
			buildDir := "."
			if len(args) > 0 {
				buildDir = args[0]
			}
			buildDir, err := filepath.Abs(buildDir)
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

			if err := ardiCore.Uploader.Upload(*target, buildDir); err != nil {
				logger.WithError(err).Errorf("Failed to upload %s", buildDir)
				return err
			}
			return nil
		},
	}
	return compileCmd
}
