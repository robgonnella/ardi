package core

import (
	"time"

	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
)

// CompileCore represents core module for compile commands
type CompileCore struct {
	logger    *log.Logger
	cli       *cli.Wrapper
	watcher   *FileWatcher
	compiling bool
}

// CompileCoreOption represents options for the CompileCore
type CompileCoreOption = func(c *CompileCore)

// NewCompileCore instance of core module for compile commands
func NewCompileCore(logger *log.Logger, options ...CompileCoreOption) *CompileCore {
	c := &CompileCore{
		logger:    logger,
		compiling: false,
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
	c.waitForCompilationsToFinish()
	c.compiling = true
	fields := log.Fields{
		"sketch": opts.SketchPath,
		"fqbn":   opts.FQBN,
	}
	fieldsLogger := c.logger.WithFields(fields)
	fieldsLogger.Info("Compiling...")
	if err := c.cli.Compile(opts); err != nil {
		fieldsLogger.WithError(err).Error("Compilation failed")
		c.compiling = false
		return err
	}
	fieldsLogger.Info("Compilation successful")
	c.compiling = false
	return nil
}

// WatchForChanges watches sketch for changes and recompiles
func (c *CompileCore) WatchForChanges(opts cli.CompileOpts) error {
	watcher, err := NewFileWatcher(opts.SketchPath, c.logger)
	if err != nil {
		return nil
	}
	c.watcher = watcher
	processingUpdate := false

	c.watcher.AddListener(func() {
		if processingUpdate {
			c.logger.Debug("ignoring update - still processing previous update")
			return
		}

		processingUpdate = true
		c.watcher.Stop()

		defer func() {
			processingUpdate = false
			c.watcher.Restart()
		}()

		c.Compile(opts)
	})

	c.watcher.Watch()
	return nil
}

// StopWatching stop watching for file changes
func (c *CompileCore) StopWatching() {
	if c.watcher != nil {
		c.watcher.Stop()
		c.watcher = nil
	}
}

// IsCompiling returns if core is currently compiling
func (c *CompileCore) IsCompiling() bool {
	return c.compiling
}

// private
func (c *CompileCore) waitForCompilationsToFinish() {
	for {
		if !c.IsCompiling() {
			break
		}
		c.logger.Info("Waiting for previous compilations to finish...")
		time.Sleep(time.Second)
	}
}
