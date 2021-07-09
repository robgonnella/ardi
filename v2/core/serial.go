package core

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// SerialPort represents a board port on which to stream logs
//go:generate mockgen -destination=../mocks/mock_serial.go -package=mocks github.com/robgonnella/ardi/v2/core SerialPort
type SerialPort interface {
	Watch() error
	Close()
	Streaming() bool
}

// ArdiSerialPort represents our serial port wrapper
type ArdiSerialPort struct {
	name          string
	baud          int
	stream        serial.Port
	stopChan      chan bool
	expectingStop bool
	logger        *log.Logger
}

// NewArdiSerialPort returns instance of serial port wrapper
func NewArdiSerialPort(name string, baud int, logger *log.Logger) SerialPort {
	return &ArdiSerialPort{
		name:          name,
		baud:          baud,
		stopChan:      make(chan bool),
		expectingStop: false,
		logger:        logger,
	}
}

// Watch connects to a serial port and prints any logs received.
func (p *ArdiSerialPort) Watch() error {
	logFields := log.Fields{"baud": p.baud, "name": p.name}

	if p.Streaming() {
		p.Close()
	}

	p.logger.WithField("port", p.name).Info("Attaching to port")

	mode := &serial.Mode{
		BaudRate: p.baud,
	}

	stream, err := serial.Open(p.name, mode)
	if err != nil {
		p.logger.WithError(err).WithFields(logFields).Warn("Failed to read from device")
		return err
	}
	p.stream = stream
	buf := make([]byte, 100)

	for {
		if !p.Streaming() {
			return nil
		}
		n, err := p.stream.Read(buf)
		if err != nil {
			p.logger.WithError(err).WithFields(logFields).Debug("Failed to read from serial port")
			p.stream = nil
			if p.expectingStop {
				p.stopChan <- true
			}
			return err
		}
		if n == 0 {
			err := errors.New("EOF")
			p.logger.WithError(err).WithField("port", p.name).Error("error reading from serial port")
			p.stream.Close()
			p.stream = nil
			return nil
		}
		fmt.Printf("%v", string(buf[:n]))
	}
}

// Close closees serial port logger
func (p *ArdiSerialPort) Close() {
	logWithField := p.logger.WithField("name", p.name)
	if p.Streaming() {
		logWithField.Info("Closing serial port connection")
		p.expectingStop = true
		p.stream.Close()
		<-p.stopChan
		p.expectingStop = false
	}
	logWithField.Info("Serial port closed")
}

// Streaming returns whether or not we are attached to the port and streaming logs
func (p *ArdiSerialPort) Streaming() bool {
	return p.stream != nil
}
