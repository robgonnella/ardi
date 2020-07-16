package commands

import (
	"fmt"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
	"github.com/spf13/cobra"
)

func buildAll() error {
	if len(ardiCore.Config.GetBuilds()) == 0 {
		logger.Warn("No builds defined. Use 'ardi add build' to define a build.")
		return nil
	}
	for buildName, build := range ardiCore.Config.GetBuilds() {
		project, err := util.ProcessSketch(build.Path)
		if err != nil {
			return err
		}
		buildProps := []string{}
		for prop, instruction := range build.Props {
			buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
		}

		logger.Infof("Building %s", build.Path)
		opts := rpc.CompileOpts{
			FQBN:       build.FQBN,
			SketchDir:  project.Directory,
			SketchPath: project.Sketch,
			ExportName: buildName,
			BuildProps: buildProps,
			ShowProps:  false,
		}
		if err := ardiCore.RPCClient.Compile(opts); err != nil {
			return err
		}
	}
	return nil
}

func build(args []string) error {
	for _, buildName := range args {
		build, ok := ardiCore.Config.GetBuilds()[buildName]
		if !ok {
			return fmt.Errorf("No build specification for %s", buildName)
		}

		project, err := util.ProcessSketch(build.Path)
		if err != nil {
			return err
		}

		buildProps := []string{}
		for prop, instruction := range build.Props {
			buildProps = append(buildProps, fmt.Sprintf("%s=%s", prop, instruction))
		}

		logger.Infof("Building %s", buildName)
		opts := rpc.CompileOpts{
			FQBN:       build.FQBN,
			SketchDir:  project.Directory,
			SketchPath: project.Sketch,
			ExportName: buildName,
			BuildProps: buildProps,
			ShowProps:  false,
		}
		if err := ardiCore.RPCClient.Compile(opts); err != nil {
			return err
		}
	}

	return nil
}

func getBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Compile configured builds",
		Long:  "\nCompile configured builds",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return buildAll()
			}
			return build(args)
		},
	}
	return buildCmd
}
