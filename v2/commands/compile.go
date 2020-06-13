package commands

import (
	"github.com/robgonnella/ardi/v2/core/compile"
	"github.com/spf13/cobra"
)

func getCompileCommand() *cobra.Command {
	var fqbn string
	var buildProps []string
	var showProps bool
	var compileCmd = &cobra.Command{
		Use:   "compile [sketch]",
		Long:  "\nCompile specified sketch",
		Short: "Compile specified sketch",
		Run: func(cmd *cobra.Command, args []string) {
			sketchDir := "."
			if len(args) > 0 {
				sketchDir = args[0]
			}

			compileCore := compile.New(client, logger)
			compileCore.Compile(sketchDir, fqbn, buildProps, showProps)
		},
	}
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")

	return compileCmd
}
