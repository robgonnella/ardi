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
		ctx := context.Background()

		ctrl := gomock.NewController(st)
		defer ctrl.Finish()

		cliInstance := mocks.NewMockCli(ctrl)
		logger := log.New()

		CleanAll()

		var b bytes.Buffer
		logger.SetOutput(&b)
		logger.SetLevel(log.DebugLevel)

		ardiConfig, svrSettings := util.GetAllSettings()
		settingsPath := util.GetCliSettingsPath()

		cliInstance.EXPECT().InitSettings(settingsPath)
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
		CleanAll()
	})
}

// IntegrationTestEnv represents our integration test environment
type IntegrationTestEnv struct {
	T      *testing.T
	Stdout *bytes.Buffer
	logger *log.Logger
}

// RunIntegrationTest runs an ardi integration test
func RunIntegrationTest(name string, t *testing.T, f func(env *IntegrationTestEnv)) {
	t.Run(name, func(st *testing.T) {
		CleanAll()

		var b bytes.Buffer
		logger := log.New()
		logger.Out = &b
		logger.SetLevel(log.InfoLevel)

		env := IntegrationTestEnv{
			T:      st,
			Stdout: &b,
			logger: logger,
		}

		f(&env)
		CleanAll()
	})
}

// RunProjectInit initializes and ardi project directory
func (e *IntegrationTestEnv) RunProjectInit() error {
	projectInitArgs := []string{"project-init"}
	return e.Execute(projectInitArgs)
}

// Execute executes the root command with given arguments
func (e *IntegrationTestEnv) Execute(args []string) error {
	rootCmd := commands.GetRootCmd(e.logger, nil)
	rootCmd.SetOut(e.logger.Out)
	rootCmd.SetArgs(args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

// ExecuteWithMockCli executes command with injected mock cli instance
func (e *IntegrationTestEnv) ExecuteWithMockCli(args []string, inst *mocks.MockCli) error {
	inst.EXPECT().InitSettings(gomock.Any())
	rootCmd := commands.GetRootCmd(e.logger, inst)
	rootCmd.SetOut(e.logger.Out)
	rootCmd.SetArgs(args)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return rootCmd.ExecuteContext(ctx)
}

// ClearStdout clears integration test env stdout
func (e *IntegrationTestEnv) ClearStdout() {
	var b bytes.Buffer
	e.logger.SetOutput(&b)
	e.Stdout = &b
}
