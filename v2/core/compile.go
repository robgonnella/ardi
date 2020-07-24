package core

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// CompileCore represents core module for compile commands
type CompileCore struct {
	logger    *log.Logger
	client    rpc.Client
	compiling bool
}

// NewCompileCore instance of core module for compile commands
func NewCompileCore(client rpc.Client, logger *log.Logger) *CompileCore {
	return &CompileCore{
		logger:    logger,
		client:    client,
		compiling: false,
	}
}

// Compile compiles a given project sketch
func (c *CompileCore) Compile(opts rpc.CompileOpts) error {
	c.waitForCompilationsToFinish()
	c.compiling = true
	if err := c.client.Compile(opts); err != nil {
		c.compiling = false
		return err
	}
	c.compiling = false
	return nil
}

// WatchForChanges watches sketch for changes and recompiles
func (c *CompileCore) WatchForChanges(opts rpc.CompileOpts) error {
	watcher, err := NewFileWatcher(opts.SketchPath, c.logger)
	if err != nil {
		return nil
	}

	watcher.AddListener(func() {
		c.waitForCompilationsToFinish()
		c.logger.Infof("Recompiling %s", opts.SketchPath)
		if err := c.Compile(opts); err != nil {
			c.logger.WithError(err).Error("Compilation failed")
		}
		c.logger.Info("Compilation successful")
	})

	watcher.Watch()

	return nil
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
