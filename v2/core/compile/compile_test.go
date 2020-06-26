package compile_test

import (
	"bytes"
	"path"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/core/compile"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/rpc"
)

type TestEnv struct {
	ctrl        *gomock.Controller
	logger      *log.Logger
	client      *mocks.MockClient
	compileCore *compile.Compile
	stdout      *bytes.Buffer
}

func runTest(name string, t *testing.T, f func(t *testing.T, env TestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctrl := gomock.NewController(t)
		client := mocks.NewMockClient(ctrl)
		logger := log.New()

		var b bytes.Buffer
		logger.SetLevel(log.DebugLevel)
		logger.SetOutput(&b)

		env := TestEnv{
			ctrl:        ctrl,
			logger:      logger,
			client:      client,
			compileCore: compile.New(client, logger),
			stdout:      &b,
		}

		f(st, env)
	})
}

func TestCompileCore(t *testing.T) {
	runTest("errors when compiling directory with no sketch", t, func(st *testing.T, env TestEnv) {
		defer env.ctrl.Finish()
		err := env.compileCore.Compile(".", "some-fqbn", []string{}, false)
		assert.Error(t, err)
	})

	runTest("succeeds when compiling directory with .ino file", t, func(st *testing.T, env TestEnv) {
		defer env.ctrl.Finish()
		here, _ := filepath.Abs(".")
		expectedFqbn := "some-fqbb"
		expectedDir := "../../test_data"
		expectedSketch := path.Join(here, "../../test_data/test.ino")
		expectedBuildProps := []string{}
		expectedShowProps := false
		compileOpts := rpc.CompileOpts{
			FQBN:       expectedFqbn,
			SketchDir:  expectedDir,
			SketchPath: expectedSketch,
			BuildProps: expectedBuildProps,
			ShowProps:  expectedShowProps,
		}
		env.client.EXPECT().ConnectedBoards().Times(1).Return([]*rpc.Board{})
		env.client.EXPECT().AllBoards().Times(1).Return([]*rpc.Board{})
		env.client.EXPECT().Compile(compileOpts).Times(1)
		err := env.compileCore.Compile(expectedDir, expectedFqbn, expectedBuildProps, expectedShowProps)
		assert.NoError(t, err)
	})
}
