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

// NewCompileCore instance of core module for compile commands
func NewCompileCore(cli *cli.Wrapper, logger *log.Logger) *CompileCore {
	return &CompileCore{
		logger:    logger,
		cli:       cli,
		compiling: false,
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
		fieldsLogger.WithError(err).Error("failed to compile")
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

	c.watcher.AddListener(func() {
		c.logger.Infof("Recompiling %s", opts.SketchPath)
		if err := c.Compile(opts); err != nil {
			c.logger.WithError(err).Error("Compilation failed")
		}
		c.logger.Info("Compilation successful")
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
