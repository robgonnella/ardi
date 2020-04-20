package compile

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/core/project"
	"github.com/robgonnella/ardi/v2/core/rpc"
	"github.com/robgonnella/ardi/v2/core/target"
	"github.com/robgonnella/ardi/v2/paths"
)

// Compile repsents core module for compile commands
type Compile struct {
	logger *log.Logger
	RPC    *rpc.RPC
}

// New instance of core module for compile commands
func New(logger *log.Logger) (*Compile, error) {
	rpc, err := rpc.New(paths.ArdiDataConfig, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize compiler")
		return nil, err
	}

	return &Compile{
		logger: logger,
		RPC:    rpc,
	}, nil
}

// Compile a given project
func (c *Compile) Compile(sketchDir, fqbn string, buildProps []string, showProps bool) error {
	if !isInitialized() {
		err := errors.New("Ardi has not been initialized. Please run \"ardi init\" first")
		c.logger.WithError(err).Error("Cannot compile")
		return err
	}

	project, err := project.New(sketchDir, c.logger)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	target, err := target.New(c.RPC, c.logger, fqbn, false)
	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
		return err
	}

	if err := c.RPC.Compile(target.Board.FQBN, project.Directory, buildProps, showProps); err != nil {
		c.logger.WithError(err).Error("Failed to compile sketch")
		return err
	}

	return nil
}

func isInitialized() bool {
	_, err := os.Stat(paths.ArdiGlobalDataDir)
	return !os.IsNotExist(err)
}
