package core

import (
	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v3/cli-wrapper"
)

// CompileCore represents core module for compile commands
type CompileCore struct {
	logger *log.Logger
	cli    *cli.Wrapper
}

// CompileCoreOption represents options for the CompileCore
type CompileCoreOption = func(c *CompileCore)

// NewCompileCore instance of core module for compile commands
func NewCompileCore(logger *log.Logger, options ...CompileCoreOption) *CompileCore {
	c := &CompileCore{
		logger: logger,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithCompileCoreCliWrapper allows an injectable cli wrapper
func WithCompileCoreCliWrapper(cliWrapper *cli.Wrapper) CompileCoreOption {
	return func(c *CompileCore) {
		c.cli = cliWrapper
	}
}

// Compile compiles a given project sketch
func (c *CompileCore) Compile(opts cli.CompileOpts) error {
	fields := log.Fields{
		"sketch": opts.SketchPath,
		"fqbn":   opts.FQBN,
	}
	fieldsLogger := c.logger.WithFields(fields)
	fieldsLogger.Info("Compiling...")
	if err := c.cli.Compile(opts); err != nil {
		fieldsLogger.WithError(err).Error("Compilation failed")
		return err
	}
	fieldsLogger.Info("Compilation successful")
	return nil
}
