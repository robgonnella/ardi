package ardi

import (
	"context"
	"io"
	"path"
	"time"

	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"google.golang.org/grpc"
)

func updateIndexes(client rpc.ArduinoCoreClient, instance *rpc.Instance) {
	updatePlatformIndex(client, instance)
	updateLibraryIndex(client, instance)
}

func getRPCInstance(client rpc.ArduinoCoreClient, configPath string) *rpc.Instance {
	configDir := path.Dir(configPath)
	initRespStream, err := client.Init(context.Background(), &rpc.InitReq{
		Configuration: &rpc.Configuration{DataDir: DataDir, SketchbookDir: configDir},
	})
	if err != nil {
		logger.Fatalf("Error creating server instance: %s", err)
	}

	var instance *rpc.Instance
	// Loop and consume the server stream until all the setup procedures are done.
	for {
		initResp, err := initRespStream.Recv()
		// The server is done.
		if err == io.EOF {
			break
		}

		// There was an error.
		if err != nil {
			logger.Fatalf("Init error: %s", err)
		}

		// The server sent us a valid instance, let's print its ID.
		if initResp.GetInstance() != nil {
			instance = initResp.GetInstance()
			logger.Debugf("Got a new instance with ID: %v", instance.GetId())
		}

		// When a download is ongoing, log the progress
		if initResp.GetDownloadProgress() != nil {
			logger.Debugf("DOWNLOAD: %s", initResp.GetDownloadProgress())
		}

		// When an overall task is ongoing, log the progress
		if initResp.GetTaskProgress() != nil {
			logger.Debugf("TASK: %s", initResp.GetTaskProgress())
		}
	}

	return instance
}

func getServerConnection() *grpc.ClientConn {
	backgroundCtx := context.Background()
	ctx, _ := context.WithTimeout(backgroundCtx, 2*time.Second)
	// Establish a connection with the gRPC server, started with the command: arduino-cli daemon
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatal("error connecting to arduino-cli rpc server, you can start it by running `arduino-cli daemon`")
	}
	return conn
}

func startDaemon(pathToConfig string) {
	logger.Debug("Starting daemon")
	cli.SetArgs([]string{"daemon", "--config-file", pathToConfig})
	if err := cli.Execute(); err != nil {
		logger.WithError(err).Fatal("Failed to start rpc server")
	}
	logger.Debug("Daemon started")
}
