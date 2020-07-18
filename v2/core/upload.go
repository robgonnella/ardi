package core

import (
	"github.com/robgonnella/ardi/v2/rpc"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
)

// UploadCore represents core module for ardi upload commands
type UploadCore struct {
	logger    *log.Logger
	client    rpc.Client
	uploading bool
}

// NewUploadCore returns new ardi upload core
func NewUploadCore(client rpc.Client, logger *log.Logger) *UploadCore {
	return &UploadCore{
		client:    client,
		logger:    logger,
		uploading: false,
	}
}

// Upload compiled sketches to the specified board
func (u *UploadCore) Upload(target Target, project types.Project) error {
	u.logger.Info("Uploading...")
	fqbn := target.Board.FQBN
	device := target.Board.Port
	sketchDir := project.Directory

	u.uploading = true
	if err := u.client.Upload(fqbn, sketchDir, device); err != nil {
		u.logger.WithError(err).Error("Failed to upload sketch")
		u.uploading = false
		return err
	}

	u.uploading = false
	return nil
}

// IsUploading returns whether or not core is currently uploading
func (u *UploadCore) IsUploading() bool {
	return u.uploading
}
