package commands

import (
	"errors"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/spf13/cobra"
)

func getCompileCmd() *cobra.Command {
	var all bool
	var fqbn string
	var buildProps []string
	var showProps bool
	var watch bool
	var compileCmd = &cobra.Command{
		Use: "compile [sketch|build(s)]",
		Long: "\nCompile sketches and builds for specified boards. When " +
			"compileing for a sketch, you must provide the board FQBN. If left " +
			"unspecified, a list of available choices will be be printed. If the " +
			"sketch argument matches a user defined build in ardi.json, the values " +
			"defined in build will be used to compile",
		Short: "Compile specified sketch or build(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ardiBuilds := ardiCore.Config.GetBuilds()

			if all {
				if watch {
					return errors.New("Cannot watch all builds. You can only watch one build at a time")
				}
				for name := range ardiBuilds {
					buildOpts := core.CompileArdiBuildOpts{
						BuildName:           name,
						OnlyConnectedBoards: false,
					}
					if _, _, err := ardiCore.CompileArdiBuild(buildOpts); err != nil {
						return err
					}
				}
				return nil
			}

			if len(args) == 0 {
				sketchOpts := core.CompileSketchOpts{
					Sketch:              ".",
					FQBN:                fqbn,
					BuildPros:           buildProps,
					ShowProps:           showProps,
					OnlyConnectedBoards: false,
				}
				opts, _, err := ardiCore.CompileSketch(sketchOpts)
				if err != nil {
					return err
				}
				if watch {
					return ardiCore.Compiler.WatchForChanges(*opts)
				}
				return nil
			}

			if len(args) == 1 {
				sketch := args[0]
				if _, ok := ardiBuilds[sketch]; ok {
					buildOpts := core.CompileArdiBuildOpts{
						BuildName:           sketch,
						OnlyConnectedBoards: false,
					}
					compileOpts, _, err := ardiCore.CompileArdiBuild(buildOpts)
					if err != nil {
						return err
					}
					if watch {
						return ardiCore.Compiler.WatchForChanges(*compileOpts)
					}
					return nil
				}

				sketchOpts := core.CompileSketchOpts{
					Sketch:              sketch,
					FQBN:                fqbn,
					BuildPros:           buildProps,
					ShowProps:           showProps,
					OnlyConnectedBoards: false,
				}
				compileOpts, _, err := ardiCore.CompileSketch(sketchOpts)
				if err != nil {
					return err
				}
				if watch {
					return ardiCore.Compiler.WatchForChanges(*compileOpts)
				}

				return nil
			}

			if watch {
				return errors.New("Cannot specifify watch with mutiple builds. You can only watch one build at a time")
			}

			for _, buildName := range args {
				buildOpts := core.CompileArdiBuildOpts{
					BuildName:           buildName,
					OnlyConnectedBoards: false,
				}
				if _, _, err := ardiCore.CompileArdiBuild(buildOpts); err != nil {
					return err
				}
			}

			return nil
		},
	}
	compileCmd.Flags().BoolVarP(&all, "all", "a", false, "Compile all builds specified in ardi.json")
	compileCmd.Flags().StringVarP(&fqbn, "fqbn", "f", "", "Specify fully qualified board name")
	compileCmd.Flags().StringArrayVarP(&buildProps, "build-prop", "p", []string{}, "Specify build property to compiler")
	compileCmd.Flags().BoolVarP(&showProps, "show-props", "s", false, "Show all build properties (does not compile)")
	compileCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch sketch file for changes and recompile")

	return withRPCConnectPreRun(compileCmd)
}
