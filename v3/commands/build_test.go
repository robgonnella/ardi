package commands_test

import (
	"errors"
	"path"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/golang/mock/gomock"
	"github.com/robgonnella/ardi/v3/testutil"
	"github.com/stretchr/testify/assert"
)

type compileReqMatcher struct {
	expectedReq *rpc.CompileRequest
}

func (m *compileReqMatcher) String() string {
	return "Matches CompileRequests"
}

func (m *compileReqMatcher) Matches(x interface{}) bool {
	req, ok := x.(*rpc.CompileRequest)

	if !ok {
		return false
	}

	if req.Instance != m.expectedReq.Instance {
		return false
	}

	if req.Fqbn != m.expectedReq.Fqbn {
		return false
	}

	if req.SketchPath != m.expectedReq.SketchPath {
		return false
	}

	if req.ShowProperties != m.expectedReq.ShowProperties {
		return false
	}

	if req.ExportDir != m.expectedReq.ExportDir {
		return false
	}

	anyOrder := gomock.InAnyOrder(m.expectedReq.BuildProperties)
	return anyOrder.Matches(req.BuildProperties)
}

type oneOfCompileReqMatcher struct {
	expectedReqs []*rpc.CompileRequest
}

func (m *oneOfCompileReqMatcher) String() string {
	return "Matches one of a list of expected CompileRequests"
}

func (m *oneOfCompileReqMatcher) Matches(x interface{}) bool {
	req, ok := x.(*rpc.CompileRequest)

	if !ok {
		return false
	}

	for _, r := range m.expectedReqs {
		matcher := &compileReqMatcher{expectedReq: r}
		if matcher.Matches(req) {
			return true
		}
	}
	return false
}

func TestBuildCommand(t *testing.T) {
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

	expectUsual := func(env *testutil.MockIntegrationTestEnv) {
		env.ArduinoCli.EXPECT().CreateInstance().Return(instance)
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
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", buildName1}
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

		oneOfReqMatcher := &oneOfCompileReqMatcher{
			expectedReqs: []*rpc.CompileRequest{req1, req2},
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), oneOfReqMatcher, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), oneOfReqMatcher, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(1)

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", buildName1, buildName2}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("shows build props", t, func(env *testutil.MockIntegrationTestEnv) {
		env.RunProjectInit()

		buildProps := []string{"some.buildProp=true", "test.anotherProps=1"}

		req := &rpc.CompileRequest{
			Instance:        instance,
			Fqbn:            fqbn1,
			SketchPath:      sketchPath1,
			ShowProperties:  true,
			BuildProperties: buildProps,
			ExportDir:       buildDir1,
		}

		var reqMatcher gomock.Matcher = &compileReqMatcher{
			expectedReq: req,
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), reqMatcher, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1, "--build-prop", buildProps[0], "--build-prop", buildProps[1]}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", buildName1, "--show-props"}
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
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), req1, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dummyErr)

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", buildName1, buildName2}
		err = env.Execute(args)
		assert.Error(env.T, err)
		assert.ErrorIs(env.T, err, dummyErr)
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

		oneOfReqMatcher := &oneOfCompileReqMatcher{
			expectedReqs: []*rpc.CompileRequest{req1, req2},
		}

		expectUsual(env)
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), oneOfReqMatcher, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		env.ArduinoCli.EXPECT().Compile(gomock.Any(), oneOfReqMatcher, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		args := []string{"add", "build", "-n", buildName1, "-f", fqbn1, "-s", sketchDir1}
		err := env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"add", "build", "-n", buildName2, "-f", fqbn2, "-s", sketchDir2}
		err = env.Execute(args)
		assert.NoError(env.T, err)

		args = []string{"build", "--all"}
		err = env.Execute(args)
		assert.NoError(env.T, err)
	})
}
