package commands

import (
	"errors"
	"fmt"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/spf13/cobra"
)

func getCompileAndUploadCmd(env *CommandEnv) *cobra.Command {
	var fqbn string
	var buildProps []string
	var showProps bool
	var baud int
	var port string
	compileAndUploadCmd := &cobra.Command{
		Use:     "compile-and-upload [build|sketch]",
		Aliases: []string{"build-and-upload", "deploy"},
		Short:   "Compiles then uploads to connected arduino board",
		Long: "\nCompiles and uploads sketches for connected boards. If a " +
			"connected board cannot be detected, you can provide the fqbn and port " +
			"via command flags. If the sketch argument matches a user defined " +
			"build in ardi.json, the values defined in build will be used to " +
			"compile and upload.",
		RunE: func(cmd *cobra.Command, args []string) error {
			optsList, err := env.ArdiCore.GetCompileOptsFromArgs(fqbn, buildProps, showProps, args)
			if err != nil {
				return err
			}
			if len(optsList) == 0 {
				return fmt.Errorf("unable to generate compile options from provided args: %s", args)
			}

			compileOpts := optsList[0]

			// Ignore errors here as user may have provided fqbn via build to mitigate
			// custom boards that don't show up via auto detect for some reason
			board, _ := env.ArdiCore.Cli.GetTargetBoard(fqbn, port, true)

			if board == nil && compileOpts.FQBN != "" && port != "" {
				board = &cli.BoardWithPort{FQBN: compileOpts.FQBN, Port: port}
			}

			if board == nil {
				return errors.New("no connected boards detected")
			}

			compileOpts.FQBN = board.FQBN
			baud = env.ArdiCore.GetBaudFromArgs(baud, args)

			sketchDir, _, err := env.ArdiCore.GetSketchPathsFromArgs(args)
			if err != nil {
				return err
			}

			if err := env.ArdiCore.Compiler.Compile(*compileOpts); err != nil {
				return err
			}

			return env.ArdiCore.Uploader.Upload(board, sketchDir)
		},
	}

	compileAndUploadCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileAndUploadCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileAndUploadCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")
	compileAndUploadCmd.Flags().IntVarP(&baud, "baud", "b", 0, "Specify baud rate when using \"attach\" flag")
	compileAndUploadCmd.Flags().StringVar(&port, "port", "", "The port your arduino board is connected to")

	return compileAndUploadCmd
}
