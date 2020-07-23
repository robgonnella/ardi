package core_test

import (
	"bytes"
	"os/exec"
	"testing"
	"time"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSerialPort(t *testing.T) {
	t.Run("streams from serial port", func(st *testing.T) {
		device := "/dev/ttyp0"
		baud := 9600
		logger := logrus.New()
		b := new(bytes.Buffer)
		logger.SetOutput(b)
		port := core.NewArdiSerialPort(device, baud, logger)

		st.Cleanup(func() {
			if port.IsStreaming() {
				port.Stop()
			}
			port = nil
		})

		go port.Watch()

		time.Sleep(time.Second)
		assert.True(st, port.IsStreaming())

		msg := "this is a tty message\n"
		cmd := exec.Command("echo", "-e", msg, ">>", device)
		err := cmd.Run()
		assert.NoError(st, err)

		port.Stop()
		assert.False(st, port.IsStreaming())
	})
}
