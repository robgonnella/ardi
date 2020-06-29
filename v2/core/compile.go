package core

import (
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/util"
)

// CompileCore represents core module for compile commands
type CompileCore struct {
	logger *log.Logger
	client rpc.Client
}

// NewCompileCore instance of core module for compile commands
func NewCompileCore(client rpc.Client, logger *log.Logger) *CompileCore {
	return &CompileCore{
		logger: logger,
		client: client,
	}
}

// Compile a given project
func (c *CompileCore) Compile(sketchDir, fqbn string, buildProps []string, showProps bool) error {
	sketchDir, sketchFile, _, err := util.ProcessSketch(sketchDir, c.logger)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	connectedBoards := c.client.ConnectedBoards()
	allBoards := c.client.AllBoards()

	target, err := NewTarget(connectedBoards, allBoards, c.logger, fqbn, false)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	opts := rpc.CompileOpts{
		FQBN:       target.Board.FQBN,
		SketchDir:  sketchDir,
		SketchPath: sketchFile,
		ExportName: "",
		BuildProps: buildProps,
		ShowProps:  showProps,
	}

	if err := c.client.Compile(opts); err != nil {
		c.logger.Error("Failed to compile sketch")
		return err
	}

	return nil
}
