package core

import (
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

// Compile a given project
func (c *CompileCore) Compile(opts rpc.CompileOpts) error {
	c.compiling = true
	if err := c.client.Compile(opts); err != nil {
		c.compiling = false
		return err
	}

	c.compiling = false
	return nil
}

// IsCompiling returns if core is currently compiling
func (c *CompileCore) IsCompiling() bool {
	return c.compiling
}
