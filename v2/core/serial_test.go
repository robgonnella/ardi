package core_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getPort() string {
	if runtime.GOOS == "linux" {
		return "/dev/ptmx"
	}
	return "/dev/ttywf"
}

func TestSerialPort(t *testing.T) {
	t.Run("streams from serial port", func(st *testing.T) {

		device := getPort()
		baud := 9600
		logger := logrus.New()
		b := new(bytes.Buffer)

		logger.SetOutput(b)
		port := core.NewArdiSerialPort(logger)
		port.SetTargets(device, baud)

		st.Cleanup(func() {
			port.Close()
			port = nil
		})

		go port.Watch()

		time.Sleep(time.Second)
		assert.True(st, port.Streaming())

		r, w, _ := os.Pipe()
		st.Cleanup(func() {
			w.Close()
			r.Close()
		})

		msg := "this is a tty message\n"
		cmd := exec.Command("echo", "-e", msg, ">>", device)
		cmd.Stdout = w
		err := cmd.Run()

		assert.NoError(st, err)
		w.Close()

		var buf bytes.Buffer
		io.Copy(&buf, r)
		assert.Contains(st, buf.String(), msg)

		port.Close()
		assert.False(st, port.Streaming())
	})
}
