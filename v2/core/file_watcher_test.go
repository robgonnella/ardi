package core_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFileWatcher(t *testing.T) {
	testutil.RunUnitTest("runs listener function on file changes", t, func(env *testutil.UnitTestEnv) {
		fileName := "test_file"
		data := "some test data\n"
		moreData := "some more test data\n"
		successMsg := "successfully called listener"

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		assert.NoError(env.T, err)
		defer func() {
			file.Close()
			os.RemoveAll(fileName)
		}()

		_, err = file.WriteString(data)
		assert.NoError(env.T, err)

		watcher, err := core.NewFileWatcher(fileName, env.Logger)
		assert.NoError(env.T, err)

		listener := func() {
			env.Logger.Info(successMsg)
			watcher.Stop()
		}
		watcher.AddListener(listener)

		env.ClearStdout()
		go watcher.Watch()

		_, err = file.WriteString(moreData)
		time.Sleep(time.Second)
		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), successMsg)
	})

	testutil.RunUnitTest("runs listener function on atomic updates", t, func(env *testutil.UnitTestEnv) {
		fileName := "test_file"
		tmpFileName := "tmpfile"
		data := "some test data\n"
		newData := "some different data\n"
		successMsg := "successfully called listener"

		tmpFile, err := os.OpenFile(tmpFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		assert.NoError(env.T, err)
		tmpFile.WriteString(newData)
		tmpFile.Close()

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		assert.NoError(env.T, err)
		file.WriteString(data)
		file.Close()

		defer func() {
			os.RemoveAll(fileName)
			os.RemoveAll(tmpFileName)
		}()

		watcher, err := core.NewFileWatcher(fileName, env.Logger)
		assert.NoError(env.T, err)

		listener := func() {
			env.Logger.Info(successMsg)
			watcher.Stop()
		}
		watcher.AddListener(listener)

		env.ClearStdout()
		go watcher.Watch()

		err = exec.Command("mv", tmpFileName, fileName).Run()
		assert.NoError(env.T, err)

		time.Sleep(time.Second)

		assert.NoError(env.T, err)
		assert.Contains(env.T, env.Stdout.String(), successMsg)
	})
}
