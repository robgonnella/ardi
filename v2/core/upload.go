package core

import (
	"time"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	log "github.com/sirupsen/logrus"
)

// UploadCore represents core module for ardi upload commands
type UploadCore struct {
	logger    *log.Logger
	cli       *cli.Wrapper
	uploading bool
	port      SerialPort
}

// NewUploadCore returns new ardi upload core
func NewUploadCore(cli *cli.Wrapper, logger *log.Logger) *UploadCore {
	return &UploadCore{
		cli:       cli,
		logger:    logger,
		uploading: false,
	}
}

// Upload compiled sketches to the specified board
func (c *UploadCore) Upload(board *cli.BoardWithPort, buildDir string) error {
	fqbn := board.FQBN
	device := board.Port

	c.waitForUploadsToFinish()
	c.uploading = true
	fields := log.Fields{
		"build":  buildDir,
		"fqbn":   board.FQBN,
		"device": board.Port,
	}
	fieldsLogger := c.logger.WithFields(fields)
	fieldsLogger.Info("Uploading...")
	if err := c.cli.Upload(fqbn, buildDir, device); err != nil {
		fieldsLogger.WithError(err).Error("Failed to upload sketch")
		c.uploading = false
		return err
	}
	fieldsLogger.Info("Upload successful")
	c.uploading = false
	return nil
}

// Attach attaches to the associated board port and prints logs
func (c *UploadCore) Attach(device string, baud int, port SerialPort) {
	if port == nil {
		c.port = NewArdiSerialPort(device, baud, c.logger)
	} else {
		c.port = port
		port.Close()
	}
	c.port.Watch()
}

// Detach detaches from the associated board port
func (c *UploadCore) Detach() {
	if c.port != nil {
		c.port.Close()
		c.port = nil
	}
}

// IsUploading returns whether or not core is currently uploading
func (c *UploadCore) IsUploading() bool {
	return c.uploading
}

// private
func (c *UploadCore) waitForUploadsToFinish() {
	for {
		if !c.IsUploading() {
			break
		}
		c.logger.Info("Waiting for previous upload to finish...")
		time.Sleep(time.Second)
	}
}
