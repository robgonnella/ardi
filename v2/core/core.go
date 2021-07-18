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
				c.logger.Infof("Using build definition: %s", k)
				compileOpts, _ := c.Config.GetCompileOpts(k)
				opts = append(opts, compileOpts)
			}
		} else if defaultExists {
			c.logger.Info("Using build definition: default")
			compileOpts, _ := c.Config.GetCompileOpts("default")
			opts = append(opts, compileOpts)
		} else {
			c.logger.Info("Using ino file in current directory")
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
			c.logger.Infof("Using build definition: %s", sketch)
			compileOpts, _ := c.Config.GetCompileOpts(sketch)
			opts = append(opts, compileOpts)
			return opts, nil
		}

		c.logger.Info("Using ino file in current directory")
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
		c.logger.Infof("Using build definition: %s", buildName)
		compileOpts, _ := c.Config.GetCompileOpts(buildName)
		opts = append(opts, compileOpts)
	}

	return opts, nil
}
