package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// SerialCore represents our serial port wrapper
type SerialCore struct {
	stream *serial.Port
	name   string
	baud   int
	logger *log.Logger
}

// NewSerialCore returns instance of serial port wrapper
func NewSerialCore(name string, baud int, logger *log.Logger) *SerialCore {
	return &SerialCore{
		name:   name,
		baud:   baud,
		logger: logger,
	}
}

// Watch connects to a serial port and prints any logs received.
func (p SerialCore) Watch() {
	logFields := log.Fields{"baud": p.baud, "name": p.name}

	p.Stop()
	p.logger.Info("Watching logs...")

	config := &serial.Config{Name: p.name, Baud: p.baud}
	stream, err := serial.OpenPort(config)
	if err != nil {
		p.logger.WithError(err).WithFields(logFields).Error("Failed to read from device")
		return
	}

	p.stream = stream

	for {
		if p.stream == nil {
			break
		}
		var buf = make([]byte, 128)
		n, err := stream.Read(buf)
		if err != nil {
			p.logger.WithError(err).WithFields(logFields).Error("Failed to read from serial port")
			return
		}
		fmt.Printf("%s", buf[:n])
	}
}

// Stop printing serial port logs
func (p SerialCore) Stop() {
	if p.stream != nil {
		logWithField := p.logger.WithField("name", p.name)
		logWithField.Info("Closing serial port connection")

		if err := p.stream.Close(); err != nil {
			logWithField.WithError(err).Error("Failed to close serial port connection")
		}

		if err := p.stream.Flush(); err != nil {
			logWithField.WithError(err).Error("Failed to flush serial port connection")
		}

		p.stream = nil
	}
}

// IsStreaming returns whether or not we are currently printing logs from port
func (p SerialCore) IsStreaming() bool {
	return p.stream != nil
}
