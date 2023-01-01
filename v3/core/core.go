package core

import (
	"context"

	cli "github.com/robgonnella/ardi/v3/cli-wrapper"
	"github.com/robgonnella/ardi/v3/paths"
	"github.com/robgonnella/ardi/v3/types"
	log "github.com/sirupsen/logrus"
)

// ArdiCore represents the core package of ardi
type ArdiCore struct {
	Cli             *cli.Wrapper
	Config          *ArdiConfig
	CliConfig       *ArdiYAML
	Lib             *LibCore
	Platform        *PlatformCore
	Compiler        *CompileCore
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

		withLibCliWrapper := WithLibCliWrapper(c.Cli)
		c.Lib = NewLibCore(c.logger, withLibCliWrapper)

		withPlatformCliWrapper := WithPlatformCliWrapper(c.Cli)
		c.Platform = NewPlatformCore(c.logger, withPlatformCliWrapper)

		withCompileCliWrapper := WithCompileCoreCliWrapper(c.Cli)
		c.Compiler = NewCompileCore(c.logger, withCompileCliWrapper)
	}
}
