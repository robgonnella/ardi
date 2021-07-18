package commands_test

import (
	"fmt"
	"testing"

	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	"github.com/robgonnella/ardi/v2/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAttachCommand(t *testing.T) {
	board := testutil.GenerateRPCBoard("Arduino Mega", "arduino:avr:mega")

	instance := &rpc.Instance{Id: int32(1)}

	platformReq := &rpc.PlatformListRequest{
		Instance:      instance,
		UpdatableOnly: false,
		All:           true,
	}

	boardItem := &rpc.BoardListItem{
		Name: board.Name,
		Fqbn: board.FQBN,
	}

	port := &rpc.DetectedPort{
		Address: board.Port,
		Boards:  []*rpc.BoardListItem{boardItem},
	}

	detectedPorts := []*rpc.DetectedPort{port}

	testutil.RunMockIntegrationTest("attaches to detected board", t, func(env *testutil.MockIntegrationTestEnv) {
		env.SerialPort.EXPECT().SetTargets(board.Port, 9600).MaxTimes(1)
		env.SerialPort.EXPECT().SetTargets("", 0).MaxTimes(1)
		env.SerialPort.EXPECT().Watch().MaxTimes(1)
		env.SerialPort.EXPECT().Close().MaxTimes(1)

		env.ArduinoCli.EXPECT().CreateInstance().Return(instance).MaxTimes(1)
		env.ArduinoCli.EXPECT().ConnectedBoards(instance.GetId()).Return(detectedPorts, nil).MaxTimes(1)
		env.ArduinoCli.EXPECT().GetPlatforms(platformReq).MaxTimes(1)

		args := []string{"attach"}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})

	testutil.RunMockIntegrationTest("attaches to provided port", t, func(env *testutil.MockIntegrationTestEnv) {
		port := "/some/port"
		baud := 14400

		env.SerialPort.EXPECT().SetTargets(port, baud).MaxTimes(1)
		env.SerialPort.EXPECT().SetTargets("", 0).MaxTimes(1)
		env.SerialPort.EXPECT().Watch().MaxTimes(1)
		env.SerialPort.EXPECT().Close().MaxTimes(1)

		args := []string{"attach", "--port", port, "--baud", fmt.Sprint(baud)}
		err := env.Execute(args)
		assert.NoError(env.T, err)
	})
}
