package commands_test

import (
	"errors"
	"os"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	fqbn1 := testutil.Esp8266WifiduinoFQBN()
	buildName1 := "blink"
	sketchDir1 := testutil.BlinkProjectDir()
	sketchPath1 := path.Join(sketchDir1, "blink.ino")
	buildDir1 := path.Join(sketchDir1, "build")

	fqbn2 := testutil.ArduinoMegaFQBN()
	buildName2 := "pixie"
	sketchDir2 := testutil.PixieProjectDir()
	sketchPath2 := path.Join(sketchDir2, "pixie.ino")
	buildDir2 := path.Join(sketchDir2, "build")

	platformReq := &rpc.PlatformListRequest{
		Instance: instance,
		All:      true,
	}

	expectUsual := func(env *testutil.MockIntegrationTestEnv) {
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId())
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
	}

	testutil.RunMockIntegrationTest("compiles ardi.json build", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir1,
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"compile", buildName1}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles multiple ardi.json builds", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		buildName1 := "blink"

		req1 := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir1,
		}

		req2 := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn2,
			SketchPath:      sketchPath2,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir2,
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req1, gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req2, gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"compile", buildName1, buildName2}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns error if one build fails", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		dummyErr := errors.New("dummy error")

		req1 := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir1,
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req1, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"compile", buildName1, buildName2}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})

	testutil.RunMockIntegrationTest("errors if attempt to watch multiple builds", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		expectUsual(env)

		args = []string{"compile", buildName1, buildName2, "--watch"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles all ardi.json builds", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		req1 := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir1,
		}

		req2 := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn2,
			SketchPath:      sketchPath2,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir2,
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req1, gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req2, gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"compile", "--all"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if attempting to watch all builds", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		expectUsual(env)

		args = []string{"compile", "--all", "--watch"}
		err = env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if .ino file not found in current directory", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		expectUsual(env)
		args := []string{"compile"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if fqbn is missing", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		expectUsual(env)
		args := []string{"compile", sketchDir1}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles directory if sketch arg missing", t, func(env *testutil.MockIntegrationTestEnv) {
		currentDir, _ := os.Getwd()
		blinkDir := testutil.BlinkProjectDir()
		os.Chdir(blinkDir)
		defer os.Chdir(currentDir)

		env.RunProjectInit()

		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  false,
			BuildProperties: []string{},
			ExportDir:       buildDir1,
		}

		boardItem := &rpc.BoardListItem{
			Name: "Some fancy board",
			Fqbn: fqbn1,
		}

		port := &rpc.DetectedPort{
			Address: "/dev/null",
			Boards:  []*rpc.BoardListItem{boardItem},
		}

		detectedPorts := []*rpc.DetectedPort{port}

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any())

		args := []string{"compile", "--fqbn", fqbn1}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if not a valid project directory", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"compile", ".", "--fqbn", fqbn1}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})
}
