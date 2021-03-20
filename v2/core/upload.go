package core

import (
	"time"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	log "github.com/sirupsen/logrus"
)

// UploadCore represents core module for ardi upload commands
type UploadCore struct {
	logger    *log.Logger
	client    cli.Client
	uploading bool
}

// NewUploadCore returns new ardi upload core
func NewUploadCore(client cli.Client, logger *log.Logger) *UploadCore {
	return &UploadCore{
		client:    client,
		logger:    logger,
		uploading: false,
	}
}

// Upload compiled sketches to the specified board
func (u *UploadCore) Upload(board *cli.Board, buildDir string) error {
	fqbn := board.FQBN
	device := board.Port

	u.waitForUploadsToFinish()
	u.uploading = true
	if err := u.client.Upload(fqbn, buildDir, device); err != nil {
		u.logger.WithError(err).Error("Failed to upload sketch")
		u.uploading = false
		return err
	}

	u.uploading = false
	return nil
}

// Attach attaches to the associated board port and prints logs
func (u *UploadCore) Attach(device string, baud int, port SerialPort) {
	if port == nil {
		port = NewArdiSerialPort(device, baud, u.logger)
	} else {
		port.Stop()
	}
	port.Watch()
}

// IsUploading returns whether or not core is currently uploading
func (u *UploadCore) IsUploading() bool {
	return u.uploading
}

// private
func (u *UploadCore) waitForUploadsToFinish() {
	for {
		if !u.IsUploading() {
			break
		}
		u.logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}
