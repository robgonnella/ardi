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
	client *rpc.Client
}

// New instance of core module for compile commands
func New(client *rpc.Client, logger *log.Logger) *Compile {
	return &Compile{
		logger: logger,
		client: client,
	}
}

// Compile a given project
func (c *Compile) Compile(sketchDir, fqbn string, buildProps []string, showProps bool) error {
	project, err := project.New(c.client, c.logger)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	if err := project.ProcessSketch(sketchDir); err != nil {
		c.logger.WithError(err).Error()
		return err
	}

	target, err := target.New(c.client, c.logger, fqbn, false)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	boardFQBN := target.Board.FQBN
	dir := project.Directory
	sketch := project.Sketch

	if err := c.client.Compile(boardFQBN, dir, sketch, "", buildProps, showProps); err != nil {
		c.logger.Error("Failed to compile sketch")
		return err
	}

	return nil
}
