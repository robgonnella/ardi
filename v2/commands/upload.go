package commands

import (
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
			if err := requireProjectInit(); err != nil {
				return err
			}
			defer env.ArdiCore.Uploader.Detach() // noop if not attached

			baud = env.ArdiCore.GetBaudFromArgs(baud, args)

			sketchDir, _, err := env.ArdiCore.GetSketchPathsFromArgs(args)
			if err != nil {
				return err
			}

			board, err := env.ArdiCore.Cli.GetTargetBoard(fqbn, port, true)
			if err != nil {
				return err
			}

			if err := env.ArdiCore.Uploader.Upload(board, sketchDir); err != nil {
				return err
			}

			if attach {
				env.ArdiCore.Uploader.SetPortTargets(board.Port, baud)
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
