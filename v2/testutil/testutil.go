package testutil

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
	"github.com/robgonnella/ardi/v2/commands"
	"github.com/robgonnella/ardi/v2/core"
	"github.com/robgonnella/ardi/v2/mocks"
	"github.com/robgonnella/ardi/v2/util"
)

// UnitTestEnv represents our unit test environment
type UnitTestEnv struct {
	T            *testing.T
	Ctx          context.Context
	Ctrl         *gomock.Controller
	Logger       *log.Logger
	Cli          *mocks.MockCli
	ArdiCore     *core.ArdiCore
	Stdout       *bytes.Buffer
	PixieProjDir string
	EmptyProjDIr string
}

// ClearStdout clears stdout for unit test
func (e *UnitTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.Logger.SetOutput(&b)
	e.Stdout = &b
}

// RunUnitTest runs an ardi unit test
func RunUnitTest(name string, t *testing.T, f func(env *UnitTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ctrl := gomock.NewController(st)
		defer cancel()
		defer CleanAll()

		cliInstance := mocks.NewMockCli(ctrl)
		logger := log.New()

		CleanAll()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		ardiConfig, svrSettings := util.GetAllSettings()
		settingsPath := util.GetCliSettingsPath()

		cliInstance.EXPECT().InitSettings(settingsPath).AnyTimes()
		cliWrapper := cli.NewCli(ctx, settingsPath, svrSettings, logger, cliInstance)

		coreOpts := core.NewArdiCoreOpts{
			Logger:             logger,
			Cli:                cliWrapper,
			ArdiConfig:         *ardiConfig,
			ArduinoCliSettings: *svrSettings,
		}
		ardiCore := core.NewArdiCore(coreOpts)

		env := UnitTestEnv{
			T:        st,
			Ctx:      ctx,
			Ctrl:     ctrl,
			Logger:   logger,
			Cli:      cliInstance,
			ArdiCore: ardiCore,
			Stdout:   &b,
		}

		f(&env)
	})
}

// IntegrationTestEnv represents our integration test environment
type IntegrationTestEnv struct {
	T      *testing.T
	Stdout *bytes.Buffer
	ctx    context.Context
	logger *log.Logger
}

// RunIntegrationTest runs an ardi integration test
func RunIntegrationTest(name string, t *testing.T, f func(env *IntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer CleanAll()

		CleanAll()

		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.InfoLevel)

		env := IntegrationTestEnv{
			T:      st,
			Stdout: &b,
			ctx:    ctx,
			logger: logger,
		}

		f(&env)
	})
}

// RunProjectInit initializes and ardi project directory
func (e *IntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project-init"}
	return e.Execute(projectInitArgs)
}

// Execute executes the root command with given arguments
func (e *IntegrationTestEnv) Execute(args []string) error {
	cmdEnv := &commands.CommandEnv{Logger: e.logger}
	rootCmd := commands.GetRootCmd(cmdEnv)
	rootCmd.SetOut(e.logger.Out)
	rootCmd.SetArgs(args)

	return rootCmd.ExecuteContext(e.ctx)
}

// ClearStdout clears integration test env stdout
func (e *IntegrationTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.logger.SetOutput(&b)
	e.Stdout = &b
}

// MockCliIntegrationTestEnv represents our integration test environment with a mocked arduino cli
type MockCliIntegrationTestEnv struct {
	T      *testing.T
	Stdout *bytes.Buffer
	Cli    *mocks.MockCli
	ctx    context.Context
	logger *log.Logger
}

// RunMockCliIntegrationTest runs an ardi integration test with mock cli
func RunMockCliIntegrationTest(name string, t *testing.T, f func(env *MockCliIntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ctrl := gomock.NewController(st)
		defer cancel()
		defer CleanAll()

		cliInstance := mocks.NewMockCli(ctrl)
		cliInstance.EXPECT().InitSettings(gomock.Any()).AnyTimes()

		CleanAll()

		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.InfoLevel)

		env := MockCliIntegrationTestEnv{
			T:      st,
			Stdout: &b,
			Cli:    cliInstance,
			logger: logger,
			ctx:    ctx,
		}

		f(&env)
	})
}

// RunProjectInit initializes and ardi project directory in mock cli test
func (e *MockCliIntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project-init"}
	return e.Execute(projectInitArgs)
}

// ClearStdout clears integration test env stdout in mock cli test
func (e *MockCliIntegrationTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.logger.SetOutput(&b)
	e.Stdout = &b
}

// Execute executes the root command with given arguments for mock cli test
func (e *MockCliIntegrationTestEnv) Execute(args []string) error {
	cmdEnv := &commands.CommandEnv{Logger: e.logger, MockCli: e.Cli}
	rootCmd := commands.GetRootCmd(cmdEnv)
	rootCmd.SetOut(e.logger.Out)
	rootCmd.SetArgs(args)

	return rootCmd.ExecuteContext(e.ctx)
}
