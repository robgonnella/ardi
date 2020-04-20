package rpc

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	arduino "github.com/arduino/arduino-cli/cli"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var cli = arduino.ArduinoCli

// RPC represents our arduino-cli rpc wrapper
type RPC struct {
	// Connection grpc client connection availabel to defer close
	Connection *grpc.ClientConn
	client     rpc.ArduinoCoreClient
	instance   *rpc.Instance
	logger     *log.Logger
}

// Board represents a single arduino Board
type Board struct {
	FQBN string
	Name string
	Port string
}

// New return new RPC controller
func New(dataConfigPath string, logger *log.Logger) (*RPC, error) {
	logger.Info("Starting daemon")
	go startDaemon(dataConfigPath)

	logger.Info("Connecting to server")
	conn, err := getServerConnection()
	if err != nil {
		logger.WithError(err).Error("error connecting to arduino-cli rpc server")
		return nil, err
	}

	client := rpc.NewArduinoCoreClient(conn)
	instance, err := getRPCInstance(client, logger)
	if err != nil {
		return nil, err
	}

	rpc := &RPC{
		Connection: conn,
		logger:     logger,
		client:     client,
		instance:   instance,
	}

	return rpc, nil
}

// UpdateIndexFiles updates platform and library index files
func (r *RPC) UpdateIndexFiles() error {
	if err := r.UpdatePlatformIndex(); err != nil {
		return err
	}
	if err := r.UpdateLibraryIndex(); err != nil {
		return err
	}
	return nil
}

// UpdateLibraryIndex updates library index file
func (r *RPC) UpdateLibraryIndex() error {
	r.logger.Debug("Updating library index...")

	libIdxUpdateStream, err := r.client.UpdateLibrariesIndex(
		context.Background(),
		&rpc.UpdateLibrariesIndexReq{
			Instance: r.instance,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Error updating libraries index")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		resp, err := libIdxUpdateStream.Recv()
		if err == io.EOF {
			r.logger.Debug("Library index update done")
			return nil
		}

		if err != nil {
			r.logger.WithError(err).Error("Error updating libraries index")
			return err
		}

		if resp.GetDownloadProgress() != nil {
			r.logger.Debugf("DOWNLOAD: %s", resp.GetDownloadProgress())
		}
	}
}

// UpdatePlatformIndex updates platform index file
func (r *RPC) UpdatePlatformIndex() error {
	r.logger.Debug("Updating platform index...")

	uiRespStream, err := r.client.UpdateIndex(
		context.Background(),
		&rpc.UpdateIndexReq{
			Instance: r.instance,
		},
	)
	if err != nil {
		r.logger.WithError(err).Error("Error updating platform index")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		uiResp, err := uiRespStream.Recv()

		// the server is done
		if err == io.EOF {
			r.logger.Debug("Platform index updated")
			return nil
		}

		// there was an error
		if err != nil {
			r.logger.WithError(err).Error("Error updating platform index")
			return err
		}

		// operations in progress
		if uiResp.GetDownloadProgress() != nil {
			r.logger.Debugf("DOWNLOAD: %s", uiResp.GetDownloadProgress())
		}
	}
}

// UpgradePlatform upgrades a given platform
func (r *RPC) UpgradePlatform(platPackage, arch string) error {
	r.logger.Debugf("Upgrading platform: %s:%s\n", platPackage, arch)

	upgradeRespStream, err := r.client.PlatformUpgrade(
		context.Background(),
		&rpc.PlatformUpgradeReq{
			Instance:        r.instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Error upgrading platform")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			r.logger.Debug("Upgrade done")
			return nil
		}

		// There was an error.
		if err != nil {
			if !strings.Contains(err.Error(), "platform already at latest version") {
				r.logger.WithError(err).Error("Cannot upgrade platform")
			}
			return err
		}

		// When a download is ongoing, log the progress
		if upgradeResp.GetProgress() != nil {
			r.logger.Debugf("DOWNLOAD: %s", upgradeResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if upgradeResp.GetTaskProgress() != nil {
			r.logger.Debugf("TASK: %s", upgradeResp.GetTaskProgress())
		}
	}
}

// InstallPlatform installs a given platform
func (r *RPC) InstallPlatform(platPackage, arch, version string) error {
	if err := r.UpdateIndexFiles(); err != nil {
		r.logger.WithError(err).Error("Failed to update index files")
		return err
	}

	r.logger.Debugf("Installing platform: %s:%s\n", arch, version)

	installRespStream, err := r.client.PlatformInstall(
		context.Background(),
		&rpc.PlatformInstallReq{
			Instance:        r.instance,
			PlatformPackage: platPackage,
			Architecture:    arch,
			Version:         version,
		})

	if err != nil {
		r.logger.WithError(err).Warn("Failed to install platform")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		installResp, err := installRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			r.logger.Debug("Install done")
			return nil
		}

		// There was an error.
		if err != nil {
			r.logger.WithError(err).Error("Failed to install platform")
			return err
		}

		// When a download is ongoing, log the progress
		if installResp.GetProgress() != nil {
			r.logger.Debugf("DOWNLOAD: %s", installResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if installResp.GetTaskProgress() != nil {
			r.logger.Debugf("TASK: %s", installResp.GetTaskProgress())
		}
	}
}

// InstallAllPlatforms installs and upgrades all platforms
func (r *RPC) InstallAllPlatforms() error {
	if err := r.UpdateIndexFiles(); err != nil {
		r.logger.WithError(err).Error("Failed to update index files")
		return err
	}

	searchResp, err := r.client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance: r.instance,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Search error")
		return err
	}

	platforms := searchResp.GetSearchOutput()

	for _, plat := range platforms {
		id := plat.GetID()
		idParts := strings.Split(id, ":")
		platPackage := idParts[0]
		arch := idParts[len(idParts)-1]
		latest := plat.GetLatest()
		r.logger.Debugf("Search result: %s: %s - %s", platPackage, id, latest)
		// Ignore individual errors when installing and upgrading all platforms
		r.InstallPlatform(platPackage, arch, latest)
		r.UpgradePlatform(platPackage, arch)
	}
	return nil
}

// ListInstalledPlatforms lists all installed platforms
func (r *RPC) ListInstalledPlatforms() error {
	listResp, err := r.client.PlatformList(
		context.Background(),
		&rpc.PlatformListReq{
			Instance: r.instance,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("List error")
		return err
	}

	r.logger.Debug("------INSTALLED PLATFORMS------")
	for _, plat := range listResp.GetInstalledPlatform() {
		r.logger.Debugf("Installed platform: %s - %s", plat.GetID(), plat.GetInstalled())
	}
	r.logger.Debug("-------------------------------")
	return nil
}

// GetPlatforms returns specified platform or all platforms if unspecified
func (r *RPC) GetPlatforms(query string) ([]*rpc.Platform, error) {
	if err := r.UpdateIndexFiles(); err != nil {
		return nil, err
	}

	searchResp, err := r.client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance:   r.instance,
			SearchArgs: query,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Platform search error")
		return nil, err
	}

	return searchResp.GetSearchOutput(), nil
}

// ConnectedBoards returns a list of connected arduino boards
func (r *RPC) ConnectedBoards() []*Board {
	boardList := []*Board{}

	boardListResp, err := r.client.BoardList(
		context.Background(),
		&rpc.BoardListReq{
			Instance: r.instance,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Board list error")
		return boardList
	}

	for _, port := range boardListResp.GetPorts() {
		for _, board := range port.GetBoards() {
			boardWithPort := Board{
				FQBN: board.GetFQBN(),
				Name: board.GetName(),
				Port: port.GetAddress(),
			}
			boardList = append(boardList, &boardWithPort)
		}
	}

	return boardList
}

// AllBoards returns a list of all supported boards
func (r *RPC) AllBoards() []*Board {
	r.logger.Debug("Getting list of supported boards...")

	boardList := []*Board{}

	listResp, err := r.client.PlatformList(
		context.Background(),
		&rpc.PlatformListReq{
			Instance: r.instance,
		},
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to get board list")
		return boardList
	}

	for _, p := range listResp.GetInstalledPlatform() {
		for _, b := range p.GetBoards() {
			b := Board{
				FQBN: b.GetFqbn(),
				Name: b.GetName(),
			}
			boardList = append(boardList, &b)
		}
	}
	return boardList
}

// Upload a sketch to target board
func (r *RPC) Upload(fqbn, sketchDir, device string) error {
	uplRespStream, err := r.client.Upload(
		context.Background(),
		&rpc.UploadReq{
			Instance:   r.instance,
			Fqbn:       fqbn,
			SketchPath: sketchDir,
			Port:       device,
			Verbose:    r.isVerbose(),
		})

	if err != nil {
		r.logger.WithError(err).Error("Failed to upload")
		return err
	}

	for {
		uplResp, err := uplRespStream.Recv()
		if err == io.EOF {
			// target.Uploading = false
			r.logger.Info("Upload complete")
			return nil
		}

		if err != nil {
			r.logger.WithError(err).Error("Failed to upload")
			// target.Uploading = false
			return err
		}

		// When an operation is ongoing you can get its output
		if resp := uplResp.GetOutStream(); resp != nil {
			r.logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := uplResp.GetErrStream(); resperr != nil {
			r.logger.Debugf("STDERR: %s", resperr)
		}
	}
}

// Compile the specified sketch
func (r *RPC) Compile(fqbn, sketchDir string, buildProps []string, showProps bool) error {

	compRespStream, err := r.client.Compile(
		context.Background(),
		&rpc.CompileReq{
			Instance:        r.instance,
			Fqbn:            fqbn,
			SketchPath:      sketchDir,
			BuildProperties: buildProps,
			ShowProperties:  showProps,
			Verbose:         true,
		})

	if err != nil {
		r.logger.WithError(err).Error("Failed to compile")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		compResp, err := compRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			return nil
		}

		// There was an error.
		if err != nil {
			r.logger.WithError(err).Error("Failed to compile")
			return err
		}

		// When an operation is ongoing you can get its output
		if resp := compResp.GetOutStream(); resp != nil {
			r.logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := compResp.GetErrStream(); resperr != nil {
			r.logger.Errorf("STDERR: %s", resperr)
		}
	}
}

// SearchLibraries searches available libraries for download
func (r *RPC) SearchLibraries(query string) ([]*rpc.SearchedLibrary, error) {
	searchResp, err := r.client.LibrarySearch(
		context.Background(),
		&rpc.LibrarySearchReq{
			Instance: r.instance,
			Query:    query,
		},
	)
	if err != nil {
		r.logger.WithError(err).Error("Error searching libraries")
		return nil, err
	}

	return searchResp.GetLibraries(), nil
}

// InstallLibrary installs specified version of a library
func (r *RPC) InstallLibrary(name, version string) (string, error) {
	installRespStream, err := r.client.LibraryInstall(
		context.Background(),
		&rpc.LibraryInstallReq{
			Instance: r.instance,
			Name:     name,
			Version:  version,
		})

	if err != nil {
		r.logger.WithError(err).Error("Error installing library")
		return "", err
	}

	foundVersion := ""

	for {
		installResp, err := installRespStream.Recv()
		if err == io.EOF {
			r.logger.Info("Lib install done")
			return foundVersion, nil
		}

		if err != nil {
			r.logger.WithError(err).Error("Library install error")
			return "", err
		}

		if installResp.GetProgress() != nil {
			r.logger.Infof("DOWNLOAD: %s\n", installResp.GetProgress())
		}
		if installResp.GetTaskProgress() != nil {
			msg := installResp.GetTaskProgress()
			lib := msg.GetName()
			r.logger.Infof("TASK: %s\n", msg)
			if foundVersion == "" {
				foundVersion = strings.Split(lib, "@")[1]
			}
		}
	}
}

// UninstallLibrary removes specified library
func (r *RPC) UninstallLibrary(name string) error {
	uninstallRespStream, err := r.client.LibraryUninstall(
		context.Background(),
		&rpc.LibraryUninstallReq{
			Instance: r.instance,
			// Assume spaces in name were intended to be underscore. This indicates
			// a potential bug in the arduino-cli package manager as names
			// potentially do not have a one-to-one mapping with regards to install
			// and remove commands. It seems as though arduino should be forcing
			// devs to name their library according to the github url.
			// @todo there has to be a better way - find it!
			Name: strings.ReplaceAll(name, " ", "_"),
		})

	if err != nil {
		r.logger.WithError(err).Error("Error uninstalling library")
		return err
	}

	for {
		uninstallRespStream, err := uninstallRespStream.Recv()
		if err == io.EOF {
			r.logger.Info("Lib uninstall done")
			return nil
		}

		if err != nil {
			r.logger.WithError(err).Error("Library install error")
			return err
		}

		if uninstallRespStream.GetTaskProgress() != nil {
			r.logger.Infof("TASK: %s\n", uninstallRespStream.GetTaskProgress())
		}
	}
}

// private
func (r *RPC) isVerbose() bool {
	return r.logger.Level == log.DebugLevel
}

// helpers
func getRPCInstance(client rpc.ArduinoCoreClient, logger *log.Logger) (*rpc.Instance, error) {
	initRespStream, err := client.Init(
		context.Background(),
		&rpc.InitReq{},
	)
	if err != nil {
		logger.Error("Error creating server instance: %s")
		return nil, err
	}

	var instance *rpc.Instance

	// Loop and consume the server stream until all the setup procedures are done.
	for {
		initResp, err := initRespStream.Recv()
		// The server is done.
		if err == io.EOF {
			return instance, nil
		}

		// There was an error.
		if err != nil {
			logger.WithError(err).Error("Init error")
			return nil, err
		}

		// The server sent us a valid instance, let's print its ID.
		if instance = initResp.GetInstance(); instance != nil {
			logger.Debugf("Got a new instance with ID: %v", instance.GetId())
			return instance, nil
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
}

func getServerConnection() (*grpc.ClientConn, error) {
	backgroundCtx := context.Background()
	ctx, cancel := context.WithTimeout(backgroundCtx, 2*time.Second)
	defer cancel()
	// Establish a connection with the gRPC server, started with the command: arduino-cli daemon
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func startDaemon(dataConfigPath string) {
	cli.SetArgs([]string{"daemon", "--config-file", dataConfigPath})
	if err := cli.Execute(); err != nil {
		fmt.Printf("Error starting daemon: %s", err.Error())
	}
}
