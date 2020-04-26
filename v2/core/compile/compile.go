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
	Client *rpc.Client
}

// New instance of core module for compile commands
func New(logger *log.Logger) (*Compile, error) {
	client, err := rpc.NewClient(logger)
	if err != nil {
		return nil, err
	}

	return &Compile{
		logger: logger,
		Client: client,
	}, nil
}

// Compile a given project
func (c *Compile) Compile(sketchDir, fqbn string, buildProps []string, showProps bool) error {
	project, err := project.New(c.logger)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}
	if err := project.ProcessSketch(sketchDir); err != nil {
		c.logger.WithError(err).Error()
		return err
	}

	target, err := target.New(c.logger, fqbn, false)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	if err := c.Client.Compile(target.Board.FQBN, project.Directory, buildProps, showProps); err != nil {
		c.logger.WithError(err).Error("Failed to compile sketch")
		return err
	}

	return nil
}
