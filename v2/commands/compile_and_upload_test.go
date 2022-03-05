package commands_test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCompileAndUploadCommand(t *testing.T) {
	instance := &rpc.Instance{Id: 1}

	fqbn1 := testutil.Esp8266WifiduinoFQBN()
	board1 := testutil.GenerateRPCBoard("Esp8266 Wifiduino", fqbn1)
	rpcPort1 := &rpc.Port{
		Address: board1.Port,
	}
	buildName1 := "blink"
	sketchDir1 := testutil.BlinkProjectDir()
	sketchPath1 := path.Join(sketchDir1, "blink.ino")
	buildDir1 := path.Join(sketchDir1, "build")

	fqbn2 := testutil.ArduinoMegaFQBN()
	board2 := testutil.GenerateRPCBoard("Arduino Mega", fqbn2)
	rpcPort2 := &rpc.Port{
		Address: board2.Port,
	}
	buildName2 := "pixie"
	sketchDir2 := testutil.PixieProjectDir()
	sketchPath2 := path.Join(sketchDir2, "pixie.ino")
	buildDir2 := path.Join(sketchDir2, "build")

	buildName_default := "default"

	compileReq1 := &rpc.CompileRequest{
		Instance:        instance,
		Fqbn:            fqbn1,
		SketchPath:      sketchPath1,
		ShowProperties:  false,
		BuildProperties: []string{},
		ExportDir:       buildDir1,
	}

	uploadReq1 := &rpc.UploadRequest{
		Instance:   instance,
		SketchPath: sketchDir1,
		Fqbn:       fqbn1,
		Port:       rpcPort1,
	}

	compileReq2 := &rpc.CompileRequest{
		Instance:        instance,
		Fqbn:            fqbn2,
		SketchPath:      sketchPath2,
		ShowProperties:  false,
		BuildProperties: []string{},
		ExportDir:       buildDir2,
	}

	uploadReq2 := &rpc.UploadRequest{
		Instance:   instance,
		SketchPath: sketchDir2,
		Fqbn:       fqbn2,
		Port:       rpcPort2,
	}

	platformReq := &rpc.PlatformListRequest{
		Instance: instance,
		All:      true,
	}

	boardReq := &rpc.BoardListRequest{
		Instance: instance,
	}

	boardItem1 := &rpc.BoardListItem{
		Name: board1.Name,
		Fqbn: board1.FQBN,
	}

	boardItem2 := &rpc.BoardListItem{
		Name: board2.Name,
		Fqbn: board2.FQBN,
	}

	port1 := &rpc.DetectedPort{
		Port:           rpcPort1,
		MatchingBoards: []*rpc.BoardListItem{boardItem1},
	}

	port2 := &rpc.DetectedPort{
		Port:           rpcPort2,
		MatchingBoards: []*rpc.BoardListItem{boardItem2},
	}

	detectedPorts1 := []*rpc.DetectedPort{port1}
	detectedPorts2 := []*rpc.DetectedPort{port2}

	expectUsual := func(env *testutil.MockIntegrationTestEnv) {
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq)
	}

	addBuild1 := func(e *testutil.MockIntegrationTestEnv) {
		err := e.RunProjectInit()
		assert.NoError(e.T, err)

		args := []string{"add", "build", "--name", buildName1, "--fqbn", fqbn1, "--sketch", sketchDir1}
		err = e.Execute(args)
		assert.NoError(e.T, err)
	}

	addBuild2 := func(e *testutil.MockIntegrationTestEnv) {
		err := e.RunProjectInit()
		assert.NoError(e.T, err)

		args := []string{"add", "build", "--name", buildName2, "--fqbn", fqbn2, "--sketch", sketchDir2}
		err = e.Execute(args)
		assert.NoError(e.T, err)
	}

	addBuild_default := func(e *testutil.MockIntegrationTestEnv) {
		err := e.RunProjectInit()
		assert.NoError(e.T, err)

		args := []string{"add", "build", "--name", buildName_default, "--fqbn", fqbn2, "--sketch", sketchDir2}
		err = e.Execute(args)
		assert.NoError(e.T, err)
	}

	testutil.RunMockIntegrationTest("compiles and uploads ardi.json build", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts1, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq1, gomock.Any(), gomock.Any())

		addBuild1(env)
		addBuild2(env)

		args := []string{"compile-and-upload", buildName1}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles and uploads single defined build", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq2, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts2, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq2, gomock.Any(), gomock.Any())

		addBuild2(env)

		fmt.Printf("builds: %+v\n", env.ArdiCore.Config.GetBuilds())

		args := []string{"compile-and-upload"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles and uploads default build", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq2, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts2, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq2, gomock.Any(), gomock.Any())

		addBuild1(env)
		addBuild_default(env)

		args := []string{"compile-and-upload"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if .ino file not found in current directory", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"compile-and-upload"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns error if build doesn't exist", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		args := []string{"compile-and-upload", "noop"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if no board is connected and fqbn is missing", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()
		expectUsual(env)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return([]*rpc.DetectedPort{}, nil)
		args := []string{"compile-and-upload", sketchDir1}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles and uploads sketch from directory", t, func(env *testutil.MockIntegrationTestEnv) {
		currentDir, _ := os.Getwd()
		os.Chdir(sketchDir1)
		defer os.Chdir(currentDir)

		env.RunProjectInit()

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts1, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq1, gomock.Any(), gomock.Any())

		args := []string{"compile-and-upload"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles and uploads sketch path", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts1, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq1, gomock.Any(), gomock.Any())

		args := []string{"compile-and-upload", sketchPath1}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("compiles and uploads using provided fqbn and port", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq1, gomock.Any(), gomock.Any())

		args := []string{"compile-and-upload", sketchPath1, "--fqbn", fqbn1, "--port", board1.Port}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("errors if not a valid project directory", t, func(env *testutil.MockIntegrationTestEnv) {
		args := []string{"compile-and-upload"}
		err := env.Execute(args)
		assert.Error(env.T, err)
	})

	testutil.RunMockIntegrationTest("returns compilation error", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		dummyErr := errors.New("dummyErr")

		expectUsual(env)
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts1, nil)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)
		addBuild1(env)

		args := []string{"compile-and-upload", buildName1}
		err := env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})

	testutil.RunMockIntegrationTest("returns upload error", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		dummyErr := errors.New("dummyErr")

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), compileReq1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().ConnectedBoards(boardReq).Return(detectedPorts1, nil)
		env.ArduinoCli.EXPECT().Upload(gomock.Any(), uploadReq1, gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		addBuild1(env)
		args := []string{"compile-and-upload", buildName1}
		err := env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
	})
}
