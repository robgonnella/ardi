package testutil

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v3/cli-wrapper"
	"github.com/robgonnella/ardi/v3/commands"
	"github.com/robgonnella/ardi/v3/core"
	"github.com/robgonnella/ardi/v3/mocks"
	"github.com/robgonnella/ardi/v3/util"
)

// UnitTestEnv represents our unit test environment
type UnitTestEnv struct {
	T            *testing.T
	Ctx          context.Context
	Logger       *log.Logger
	ArduinoCli   *mocks.MockCli
	SerialPort   *mocks.MockSerialPort
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
		defer cancel()
		defer CleanAll()

		cliCtrl := gomock.NewController(st)
		defer cliCtrl.Finish()
		cliInstance := mocks.NewMockCli(cliCtrl)

		portCtrl := gomock.NewController(st)
		defer portCtrl.Finish()
		portInatance := mocks.NewMockSerialPort(portCtrl)
		logger := log.New()

		CleanAll()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		ardiConfig, svrSettings := util.GetAllSettings()
		settingsPath := util.GetCliSettingsPath()

		cliInstance.EXPECT().InitSettings(settingsPath).AnyTimes()
		withArduinoCli := core.WithArduinoCli(cliInstance)

		coreOpts := core.NewArdiCoreOpts{
			Ctx:                ctx,
			Logger:             logger,
			CliSettingsPath:    settingsPath,
			ArdiConfig:         *ardiConfig,
			ArduinoCliSettings: *svrSettings,
		}
		ardiCore := core.NewArdiCore(coreOpts, withArduinoCli)

		env := UnitTestEnv{
			T:          st,
			Ctx:        ctx,
			Logger:     logger,
			ArduinoCli: cliInstance,
			SerialPort: portInatance,
			ArdiCore:   ardiCore,
			Stdout:     &b,
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
	ardiConfig, svrSettings := util.GetAllSettings()
	cliSettingsPath := util.GetCliSettingsPath()

	coreOpts := core.NewArdiCoreOpts{
		Ctx:                e.ctx,
		Logger:             e.logger,
		CliSettingsPath:    cliSettingsPath,
		ArdiConfig:         *ardiConfig,
		ArduinoCliSettings: *svrSettings,
	}

	arduinoCli := cli.NewArduinoCli()
	withArduinoCli := core.WithArduinoCli(arduinoCli)
	ardiCore := core.NewArdiCore(coreOpts, withArduinoCli)

	env := &commands.CommandEnv{
		ArdiCore: ardiCore,
		Logger:   e.logger,
	}

	rootCmd := commands.NewRootCmd(env)
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

// MockIntegrationTestEnv represents our integration test environment with a mocked arduino cli
type MockIntegrationTestEnv struct {
	T          *testing.T
	Stdout     *bytes.Buffer
	ArdiCore   *core.ArdiCore
	ArduinoCli *mocks.MockCli
	SerialPort *mocks.MockSerialPort
	ctx        context.Context
	logger     *log.Logger
}

// RunMockIntegrationTest runs an ardi integration test with mock cli
func RunMockIntegrationTest(name string, t *testing.T, f func(env *MockIntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer CleanAll()

		cliCtrl := gomock.NewController(st)
		defer cliCtrl.Finish()
		cliInstance := mocks.NewMockCli(cliCtrl)
		cliInstance.EXPECT().InitSettings(gomock.Any()).AnyTimes()

		portCtrl := gomock.NewController(st)
		defer portCtrl.Finish()
		portInstance := mocks.NewMockSerialPort(portCtrl)

		CleanAll()

		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.InfoLevel)

		env := MockIntegrationTestEnv{
			T:          st,
			Stdout:     &b,
			ArduinoCli: cliInstance,
			SerialPort: portInstance,
			logger:     logger,
			ctx:        ctx,
		}

		f(&env)
	})
}

// RunProjectInit initializes and ardi project directory in mock cli test
func (e *MockIntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project-init"}
	return e.Execute(projectInitArgs)
}

// ClearStdout clears integration test env stdout in mock cli test
func (e *MockIntegrationTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.logger.SetOutput(&b)
	e.Stdout = &b
}

// Execute executes the root command with given arguments for mock cli test
func (e *MockIntegrationTestEnv) Execute(args []string) error {
	ardiConfig, svrSettings := util.GetAllSettings()
	cliSettingsPath := util.GetCliSettingsPath()

	coreOpts := core.NewArdiCoreOpts{
		Ctx:                e.ctx,
		Logger:             e.logger,
		CliSettingsPath:    cliSettingsPath,
		ArdiConfig:         *ardiConfig,
		ArduinoCliSettings: *svrSettings,
	}

	withArduinoCli := core.WithArduinoCli(e.ArduinoCli)
	ardiCore := core.NewArdiCore(coreOpts, withArduinoCli)

	e.ArdiCore = ardiCore

	env := &commands.CommandEnv{
		ArdiCore: ardiCore,
		Logger:   e.logger,
	}

	rootCmd := commands.NewRootCmd(env)
	rootCmd.SetOut(e.logger.Out)
	rootCmd.SetArgs(args)

	return rootCmd.ExecuteContext(e.ctx)
}
