package core

import (
	"time"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	log "github.com/sirupsen/logrus"
)

// UploadCore represents core module for ardi upload commands
type UploadCore struct {
	logger      *log.Logger
	cli         *cli.Wrapper
	portManager SerialPort
	uploading   bool
}

// UploadCoreOption reprents options for UploadCore
type UploadCoreOption = func(c *UploadCore)

// NewUploadCore returns new ardi upload core
func NewUploadCore(logger *log.Logger, options ...UploadCoreOption) *UploadCore {
	c := &UploadCore{
		logger:    logger,
		uploading: false,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithUploadCoreCliWrapper allows an injectable cli wrapper
func WithUploadCoreCliWrapper(wrapper *cli.Wrapper) UploadCoreOption {
	return func(c *UploadCore) {
		c.cli = wrapper
	}
}

// WithUploaderSerialPortManager allows and injectable serial port manager
func WithUploaderSerialPortManager(portManager SerialPort) UploadCoreOption {
	return func(c *UploadCore) {
		c.portManager = portManager
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

// SetPortTargets sets the device and baud for the port manager
func (c *UploadCore) SetPortTargets(device string, baud int) {
	c.portManager.SetTargets(device, baud)
}

// Attach attaches to the associated board port and prints logs
func (c *UploadCore) Attach() {
	c.portManager.Watch()
}

// Detach detaches from the associated board port
func (c *UploadCore) Detach() {
	c.portManager.Close()
	c.portManager.SetTargets("", 0)
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
