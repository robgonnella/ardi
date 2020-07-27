package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// SerialPort represents a board port on which to stream logs
//go:generate mockgen -destination=../mocks/mock_serial.go -package=mocks github.com/robgonnella/ardi/v2/core SerialPort
type SerialPort interface {
	Watch() error
	Stop()
	IsStreaming() bool
}

// ArdiSerialPort represents our serial port wrapper
type ArdiSerialPort struct {
	stream *serial.Port
	name   string
	baud   int
	logger *log.Logger
}

// NewArdiSerialPort returns instance of serial port wrapper
func NewArdiSerialPort(name string, baud int, logger *log.Logger) SerialPort {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serialPort := &ArdiSerialPort{
		name:   name,
		baud:   baud,
		logger: logger,
	}

	go func() {
		<-sigs
		logger.Debug("gracefully shutting down serial port stream")
		serialPort.Stop()
	}()

	return serialPort
}

// Watch connects to a serial port and prints any logs received.
func (p *ArdiSerialPort) Watch() error {
	defer p.Stop()

	logFields := log.Fields{"baud": p.baud, "name": p.name}

	p.Stop()
	p.logger.Info("Watching logs...")

	config := &serial.Config{Name: p.name, Baud: p.baud}

	stream, err := serial.OpenPort(config)
	if err != nil {
		p.logger.WithError(err).WithFields(logFields).Warn("Failed to read from device")
		return err
	}

	p.stream = stream

	for {
		if p.stream == nil {
			break
		}
		var buf = make([]byte, 128)
		n, err := stream.Read(buf)
		if err != nil {
			p.logger.WithError(err).WithFields(logFields).Warn("Failed to read from serial port")
			return err
		}
		fmt.Printf("%s", buf[:n])
	}

	return nil
}

// Stop printing serial port logs
func (p *ArdiSerialPort) Stop() {
	logWithField := p.logger.WithField("name", p.name)

	if p.stream != nil {
		logWithField.Debug("Closing serial port connection")

		if err := p.stream.Flush(); err != nil {
			logWithField.WithError(err).Debug("Failed to flush serial port connection")
		}

		if err := p.stream.Close(); err != nil {
			logWithField.WithError(err).Debug("Failed to close serial port connection")
		}

		p.stream = nil
		logWithField.Debug("Serial port closed")
	}
}

// IsStreaming returns whether or not we are currently printing logs from port
func (p *ArdiSerialPort) IsStreaming() bool {
	return p.stream != nil
}
