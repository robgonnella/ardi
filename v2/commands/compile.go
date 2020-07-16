package commands

import (
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
			if err := ardiCore.Compiler.Compile(sketchDir, fqbn, buildProps, showProps); err != nil {
				logger.WithError(err).Errorf("Failed to compile %s", sketchDir)
				return err
			}
			return nil
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return compileCmd
}
