package core

import (
	"context"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/paths"
	"github.com/robgonnella/ardi/v2/types"
	"github.com/robgonnella/ardi/v2/util"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	Cli             *cli.Wrapper
	Config          *ArdiConfig
	CliConfig       *ArdiYAML
	Watcher         *WatchCore
	Board           *BoardCore
	Compiler        *CompileCore
	Uploader        *UploadCore
	Lib             *LibCore
	Platform        *PlatformCore
	ctx             context.Context
	cliSettingsPath string
	logger          *log.Logger
}

// ArdiCoreOption represents options for ArdiCore
type ArdiCoreOption = func(c *ArdiCore)

// NewArdiCoreOpts options fore creating new ardi core
type NewArdiCoreOpts struct {
	ArdiConfig         types.ArdiConfig
	ArduinoCliSettings types.ArduinoCliSettings
	CliSettingsPath    string
	Logger             *log.Logger
	Ctx                context.Context
}

// NewArdiCore returns a new ardi core
func NewArdiCore(opts NewArdiCoreOpts, options ...ArdiCoreOption) *ArdiCore {
	ardiConf := paths.ArdiProjectConfig
	cliConf := paths.ArduinoCliProjectConfig

	core := &ArdiCore{
		ctx:             opts.Ctx,
		cliSettingsPath: opts.CliSettingsPath,
		Config:          NewArdiConfig(ardiConf, opts.ArdiConfig, opts.Logger),
		CliConfig:       NewArdiYAML(cliConf, opts.ArduinoCliSettings),
		logger:          opts.Logger,
	}

	for _, o := range options {
		o(core)
	}

	return core
}

// WithArduinoCli allows an injectable arduino cli interface
func WithArduinoCli(arduinoCli cli.Cli) func(c *ArdiCore) {
	return func(c *ArdiCore) {
		withArduinoCli := cli.WithArduinoCli(arduinoCli)
		c.Cli = cli.NewCli(c.ctx, c.cliSettingsPath, c.logger, withArduinoCli)

		serialPortManager := NewArdiSerialPort(c.logger)

		withCompileWrapper := WithCompileCoreCliWrapper(c.Cli)
		withUploadWrapper := WithUploadCoreCliWrapper(c.Cli)
		withUploadPortManager := WithUploaderSerialPortManager(serialPortManager)

		c.Compiler = NewCompileCore(c.logger, withCompileWrapper)
		c.Uploader = NewUploadCore(c.logger, withUploadWrapper, withUploadPortManager)

		withWathCompiler := WithWatchCoreCompiler(c.Compiler)
		withWatchUploader := WithWatchCoreUploader(c.Uploader)

		c.Watcher = NewWatchCore(c.logger, withWatchUploader, withWathCompiler)

		withBoardCliWrapper := WithBoardCliWrapper(c.Cli)
		c.Board = NewBoardCore(c.logger, withBoardCliWrapper)

		withLibCliWrapper := WithLibCliWrapper(c.Cli)
		c.Lib = NewLibCore(c.logger, withLibCliWrapper)

		withPlatformCliWrapper := WithPlatformCliWrapper(c.Cli)
		c.Platform = NewPlatformCore(c.logger, withPlatformCliWrapper)
	}
}

// WithCoreSerialPortManager allows an injectable serial port interface
func WithCoreSerialPortManager(port SerialPort) ArdiCoreOption {
	return func(c *ArdiCore) {
		withCli := WithUploadCoreCliWrapper(c.Cli)
		withPort := WithUploaderSerialPortManager(port)

		c.Uploader = NewUploadCore(c.logger, withCli, withPort)

		withWathCompiler := WithWatchCoreCompiler(c.Compiler)
		withWatchUploader := WithWatchCoreUploader(c.Uploader)

		c.Watcher = NewWatchCore(c.logger, withWatchUploader, withWathCompiler)
	}
}

// GetCompileOptsFromArgs return a list of compile opts from the given args
func (c *ArdiCore) GetCompileOptsFromArgs(fqbn string, buildProps []string, showProps bool, args []string) ([]*cli.CompileOpts, error) {
	ardiBuilds := c.Config.GetBuilds()
	_, defaultExists := ardiBuilds["default"]

	opts := []*cli.CompileOpts{}

	if len(args) == 0 {
		if len(ardiBuilds) == 1 {
			for k := range ardiBuilds {
				c.logger.Infof("Using compile opts from build definition: %s", k)
				compileOpts, _ := c.Config.GetCompileOpts(k)
				opts = append(opts, compileOpts)
			}
		} else if defaultExists {
			c.logger.Info("Using compile opts from build definition: default")
			compileOpts, _ := c.Config.GetCompileOpts("default")
			opts = append(opts, compileOpts)
		} else {
			c.logger.Info("Creating compile opts from provided values and current directory sketch")
			project, err := util.ProcessSketch(".")
			if err != nil {
				return nil, err
			}
			compileOpts := &cli.CompileOpts{
				FQBN:       fqbn,
				BuildProps: buildProps,
				ShowProps:  showProps,
				SketchDir:  project.Directory,
				SketchPath: project.Sketch,
			}
			opts = append(opts, compileOpts)
		}

		return opts, nil
	}

	if len(args) == 1 {
		sketch := args[0]
		if _, ok := ardiBuilds[sketch]; ok {
			c.logger.Infof("Using compile opts from build definition: %s", sketch)
			compileOpts, _ := c.Config.GetCompileOpts(sketch)
			opts = append(opts, compileOpts)
			return opts, nil
		}

		c.logger.WithField("sketch", sketch).Info("Creating compile opts from provided values")
		project, err := util.ProcessSketch(sketch)
		if err != nil {
			return nil, err
		}

		compileOpts := &cli.CompileOpts{
			FQBN:       fqbn,
			BuildProps: buildProps,
			ShowProps:  showProps,
			SketchDir:  project.Directory,
			SketchPath: project.Sketch,
		}

		opts = append(opts, compileOpts)

		return opts, nil
	}

	for _, buildName := range args {
		c.logger.Infof("Using compile opts from build definition: %s", buildName)
		compileOpts, _ := c.Config.GetCompileOpts(buildName)
		opts = append(opts, compileOpts)
	}

	return opts, nil
}

// GetBaudFromArgs return baud rate from the given args
func (c *ArdiCore) GetBaudFromArgs(baudArg int, args []string) int {
	if baudArg != 0 {
		c.logger.WithField("baud", baudArg).Info("Using user provided baud")
		return baudArg
	}

	defaultBaud := 9600

	ardiBuilds := c.Config.GetBuilds()
	defaultBuild, defaultExists := ardiBuilds["default"]

	if len(args) == 0 {
		if len(ardiBuilds) == 1 {
			for k, v := range ardiBuilds {
				c.logger.WithField("baud", v.Baud).Infof("Using baud from build definition: %s", k)
				return v.Baud
			}
		}
		if defaultExists {
			c.logger.WithField("baud", defaultBuild.Baud).Info("Using baud from default build definition")
			return defaultBuild.Baud
		}
		project, err := util.ProcessSketch(".")
		if err != nil {
			c.logger.WithField("defaultBaud", defaultBaud).Info("Unable to parse sketch file in current directory, using defualt baud rate")
			return defaultBaud
		}

		fields := log.Fields{
			"baud":   project.Baud,
			"sketch": project.Sketch,
		}
		c.logger.WithFields(fields).Info("Using baud parsed from sketch file")
		return project.Baud

	}

	sketch := args[0]
	if b, ok := ardiBuilds[sketch]; ok {
		c.logger.WithField("baud", b.Baud).Infof("Using baud from build definition: %s", sketch)
		return b.Baud
	}

	project, err := util.ProcessSketch(sketch)
	if err != nil {
		fields := log.Fields{
			"defaultBaud": defaultBaud,
			"sketch":      sketch,
		}
		c.logger.WithFields(fields).Info("Unable to parse sketch file, using defualt baud rate")
		return defaultBaud
	}

	fields := log.Fields{
		"baud":   project.Baud,
		"sketch": project.Sketch,
	}
	c.logger.WithFields(fields).Info("Using baud parsed from sketch file")
	return project.Baud
}

// GetSketchPathsFromArgs returns sketchDir and sketchPath from given args
func (c *ArdiCore) GetSketchPathsFromArgs(args []string) (string, string, error) {
	builds := c.Config.GetBuilds()
	defaultBuild, defaultExits := builds["default"]

	if len(args) == 0 {
		if defaultExits {
			fields := log.Fields{
				"directory": defaultBuild.Directory,
				"sketch":    defaultBuild.Sketch,
			}
			c.logger.WithFields(fields).Info("Using build definition: default")
			return defaultBuild.Directory, defaultBuild.Sketch, nil
		}

		if len(builds) == 1 {
			for name, b := range builds {
				fields := log.Fields{
					"directory": b.Directory,
					"sketch":    b.Sketch,
				}
				c.logger.WithFields(fields).Infof("Using build definition: %s", name)
				return b.Directory, b.Sketch, nil
			}
		}

		project, err := util.ProcessSketch(".")
		if err != nil {
			return "", "", err
		}

		fields := log.Fields{
			"directory": project.Directory,
			"sketch":    project.Sketch,
		}
		c.logger.WithFields(fields).Info("Using ino in current directory")
		return project.Directory, project.Sketch, nil
	}

	if b, ok := builds[args[0]]; ok {
		fields := log.Fields{
			"directory": b.Directory,
			"sketch":    b.Sketch,
		}
		c.logger.WithFields(fields).Infof("Using build definition: %s", args[0])
		return b.Directory, b.Sketch, nil
	}

	project, err := util.ProcessSketch(args[0])
	if err != nil {
		return "", "", err
	}
	fields := log.Fields{
		"directory": project.Directory,
		"sketch":    project.Sketch,
	}
	c.logger.WithFields(fields).Infof("Using sketch: %s", args[0])
	return project.Directory, project.Sketch, nil
}
