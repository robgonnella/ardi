package core_test

import (
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestArdiCore(t *testing.T) {
	testutil.RunUnitTest("compiles ardi build", t, func(env *testutil.UnitTestEnv) {
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"
		exportDir := path.Join(testutil.BlinkProjectDir(), "build")

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, 0, []string{})
		assert.NoError(env.T, err)

		expectedCompileOpts := cli.CompileOpts{
			FQBN:       fqbn,
			SketchDir:  testutil.BlinkProjectDir(),
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn,
			SketchPath:      sketch,
			ExportDir:       exportDir,
			BuildProperties: []string{},
			ShowProperties:  false,
			Verbose:         true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		compileOpts, err := env.ArdiCore.CompileArdiBuild(buildName)
		assert.NoError(env.T, err)
		assert.Equal(env.T, &expectedCompileOpts, compileOpts)
	})

	testutil.RunUnitTest("compiles sketch", t, func(env *testutil.UnitTestEnv) {
		buildName := "somebuild"
		sketch := path.Join(testutil.BlinkProjectDir(), "blink.ino")
		fqbn := "someboardfqbn"
		exportDir := path.Join(testutil.BlinkProjectDir(), "build")

		err := env.ArdiCore.Config.AddBuild(buildName, sketch, fqbn, 0, []string{})
		assert.NoError(env.T, err)

		sketchOpts := core.CompileSketchOpts{
			Sketch:    sketch,
			FQBN:      fqbn,
			BuildPros: []string{},
			ShowProps: false,
		}

		expectedCompileOpts := cli.CompileOpts{
			FQBN:       fqbn,
			SketchDir:  testutil.BlinkProjectDir(),
			SketchPath: sketch,
			BuildProps: []string{},
			ShowProps:  false,
		}

		instance := &rpc.Instance{Id: int32(1)}
		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn,
			SketchPath:      sketch,
			ExportDir:       exportDir,
			BuildProperties: []string{},
			ShowProperties:  false,
			Verbose:         true,
		}

		env.Cli.EXPECT().CreateInstance().Return(instance).AnyTimes()
		env.Cli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		compileOpts, err := env.ArdiCore.CompileSketch(sketchOpts)
		assert.NoError(env.T, err)
		assert.Equal(env.T, &expectedCompileOpts, compileOpts)
	})
}
