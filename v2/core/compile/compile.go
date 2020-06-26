package compile

import (
	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/core/project"
	"github.com/robgonnella/ardi/v2/core/target"
	"github.com/robgonnella/ardi/v2/rpc"
)

// Compile represents core module for compile commands
type Compile struct {
	logger *log.Logger
	client rpc.Client
}

// New instance of core module for compile commands
func New(client rpc.Client, logger *log.Logger) *Compile {
	return &Compile{
		logger: logger,
		client: client,
	}
}

// Compile a given project
func (c *Compile) Compile(sketchDir, fqbn string, buildProps []string, showProps bool) error {
	sketchDir, sketchFile, _, err := project.ProcessSketch(sketchDir, c.logger)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	connectedBoards := c.client.ConnectedBoards()
	allBoards := c.client.AllBoards()

	target, err := target.New(connectedBoards, allBoards, c.logger, fqbn, false)
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
