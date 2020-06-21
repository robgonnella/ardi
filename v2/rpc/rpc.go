package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	arduino "github.com/arduino/arduino-cli/cli"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var cli = arduino.ArduinoCli
var port = "50051"

// Client reprents our wrapper around the arduino-cli rpc client
//go:generate mockgen -source=rpc.go -destination=mock_rpc/mock_rpc.go
type Client interface {
	UpdateIndexFiles() error
	UpdateLibraryIndex() error
	UpdatePlatformIndex() error
	UpgradePlatform(platform string) error
	InstallPlatform(platform string) error
	UninstallPlatform(platform string) error
	InstallAllPlatforms() error
	GetInstalledPlatforms() ([]*rpc.Platform, error)
	GetPlatforms() ([]*rpc.Platform, error)
	ConnectedBoards() []*Board
	AllBoards() []*Board
	Upload(fqbn, sketchDir, device string) error
	Compile(o CompileOpts) error
	SearchLibraries(query string) ([]*rpc.SearchedLibrary, error)
	InstallLibrary(name, version string) (string, error)
	UninstallLibrary(name string) error
}

// ArdiClient represents a client connection to arduino-cli grpc daemon
type ArdiClient struct {
	connection *grpc.ClientConn
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

// NewClient return new RPC controller
func NewClient(logger *log.Logger) (Client, error) {
	logger.Debug("Connecting to server")
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

	return &ArdiClient{
		connection: conn,
		logger:     logger,
		client:     client,
		instance:   instance,
	}, nil
}

//StartDaemon starts the arduino-cli grpc server locally
func StartDaemon(dataConfigPath string, verbose bool) {
	args := []string{
		"daemon",
		"--port",
		port,
		"--config-file",
		dataConfigPath,
	}
	if verbose {
		args = append(args, "--verbose")
	}
	cli.SetArgs(args)
	if err := cli.Execute(); err != nil {
		fmt.Printf("Error starting daemon: %s", err.Error())
	}
}

// UpdateIndexFiles updates platform and library index files
func (c *ArdiClient) UpdateIndexFiles() error {
	if err := c.UpdatePlatformIndex(); err != nil {
		return err
	}
	if err := c.UpdateLibraryIndex(); err != nil {
		return err
	}
	return nil
}

// UpdateLibraryIndex updates library index file
func (c *ArdiClient) UpdateLibraryIndex() error {
	c.logger.Debug("Updating library index...")

	libIdxUpdateStream, err := c.client.UpdateLibrariesIndex(
		context.Background(),
		&rpc.UpdateLibrariesIndexReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Error updating libraries index")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		resp, err := libIdxUpdateStream.Recv()
		if err == io.EOF {
			c.logger.Debug("Library index update done")
			return nil
		}

		if err != nil {
			c.logger.WithError(err).Error("Error updating libraries index")
			return err
		}

		if resp.GetDownloadProgress() != nil {
			c.logger.Debugf("DOWNLOAD: %s", resp.GetDownloadProgress())
		}
	}
}

// UpdatePlatformIndex updates platform index file
func (c *ArdiClient) UpdatePlatformIndex() error {
	c.logger.Debug("Updating platform index...")

	uiRespStream, err := c.client.UpdateIndex(
		context.Background(),
		&rpc.UpdateIndexReq{
			Instance: c.instance,
		},
	)
	if err != nil {
		c.logger.WithError(err).Error("Error updating platform index")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		uiResp, err := uiRespStream.Recv()

		// the server is done
		if err == io.EOF {
			c.logger.Debug("Platform index updated")
			return nil
		}

		// there was an error
		if err != nil {
			c.logger.WithError(err).Error("Error updating platform index")
			return err
		}

		// operations in progress
		if uiResp.GetDownloadProgress() != nil {
			c.logger.Debugf("DOWNLOAD: %s", uiResp.GetDownloadProgress())
		}
	}
}

// UpgradePlatform upgrades a given platform
func (c *ArdiClient) UpgradePlatform(platform string) error {
	pkg, arch, _ := parsePlatform(platform)
	c.logger.Debugf("Upgrading platform: %s:%s\n", pkg, arch)

	upgradeRespStream, err := c.client.PlatformUpgrade(
		context.Background(),
		&rpc.PlatformUpgradeReq{
			Instance:        c.instance,
			PlatformPackage: pkg,
			Architecture:    arch,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Error upgrading platform")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		upgradeResp, err := upgradeRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			c.logger.Debug("Upgrade done")
			return nil
		}

		// There was an error.
		if err != nil {
			if !strings.Contains(err.Error(), "platform already at latest version") {
				c.logger.WithError(err).Error("Cannot upgrade platform")
			}
			return err
		}

		// When a download is ongoing, log the progress
		if upgradeResp.GetProgress() != nil {
			c.logger.Debugf("DOWNLOAD: %s", upgradeResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if upgradeResp.GetTaskProgress() != nil {
			c.logger.Debugf("TASK: %s", upgradeResp.GetTaskProgress())
		}
	}
}

// InstallPlatform installs a given platform
func (c *ArdiClient) InstallPlatform(platform string) error {
	if platform == "" {
		err := errors.New("Must specify a platform to install")
		c.logger.WithError(err).Error()
		return err
	}

	pkg, arch, version := parsePlatform(platform)

	c.logger.Infof("Installing platform: %s:%s\n", pkg, arch)

	installRespStream, err := c.client.PlatformInstall(
		context.Background(),
		&rpc.PlatformInstallReq{
			Instance:        c.instance,
			PlatformPackage: pkg,
			Architecture:    arch,
			Version:         version,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Failed to install platform")
		return err
	}

	installedVersion := ""
	// Loop and consume the server stream until all the operations are done.
	for {
		installResp, err := installRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			c.logger.Infof("Installed: %s:%s@%s", pkg, arch, installedVersion)
			return nil
		}

		// There was an error.
		if err != nil {
			c.logger.WithError(err).Error("Failed to install platform")
			return err
		}

		// When a download is ongoing, log the progress
		if installResp.GetProgress() != nil {
			c.logger.Debugf("DOWNLOAD: %s", installResp.GetProgress())
		}

		// When an overall task is ongoing, log the progress
		if progress := installResp.GetTaskProgress(); progress != nil {
			c.logger.Debugf("TASK: %s", progress)
			name := progress.GetName()
			msg := progress.GetMessage()
			if strings.Contains(msg, "installed") {
				if parts := strings.Split(msg, "@"); len(parts) > 1 {
					installedVersion = strings.Replace(parts[1], " installed", "", 1)
				}
			}
			if strings.Contains(name, "already installed") {
				if parts := strings.Split(name, "@"); len(parts) > 1 {
					installedVersion = strings.Replace(parts[1], " already installed", "", 1)
				}
			}
		}
	}
}

// UninstallPlatform installs a given platform
func (c *ArdiClient) UninstallPlatform(platform string) error {
	if platform == "" {
		err := errors.New("Must specify a platform to install")
		c.logger.WithError(err).Error()
		return err
	}

	pkg, arch, _ := parsePlatform(platform)

	c.logger.Infof("Uninstalling platform: %s:%s\n", pkg, arch)

	uninstallRespStream, err := c.client.PlatformUninstall(
		context.Background(),
		&rpc.PlatformUninstallReq{
			Instance:        c.instance,
			PlatformPackage: pkg,
			Architecture:    arch,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Failed to uninstall platform")
		return err
	}

	// Loop and consume the server stream until all the operations are done.
	for {
		uninstallResp, err := uninstallRespStream.Recv()

		// The server is done.
		if err == io.EOF {
			c.logger.Infof("Uninstalled: %s:%s", pkg, arch)
			return nil
		}

		// There was an error.
		if err != nil {
			c.logger.WithError(err).Error("Failed to uninstall platform")
			return err
		}

		// When an overall task is ongoing, log the progress
		if progress := uninstallResp.GetTaskProgress(); progress != nil {
			c.logger.Debugf("TASK: %s", progress)
		}
	}
}

// InstallAllPlatforms installs and upgrades all platforms
func (c *ArdiClient) InstallAllPlatforms() error {

	searchResp, err := c.client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Search error")
		return err
	}

	platforms := searchResp.GetSearchOutput()

	for _, plat := range platforms {
		id := plat.GetID()
		latest := plat.GetLatest()
		platform := fmt.Sprintf("%s@%s", id, latest)
		c.logger.Debugf("Search result: %s", platform)
		// Ignore individual errors when installing and upgrading all platforms
		c.InstallPlatform(platform)
		c.UpgradePlatform(platform)
	}
	return nil
}

// GetInstalledPlatforms lists all installed platforms
func (c *ArdiClient) GetInstalledPlatforms() ([]*rpc.Platform, error) {
	listResp, err := c.client.PlatformList(
		context.Background(),
		&rpc.PlatformListReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("List error")
		return nil, err
	}

	return listResp.GetInstalledPlatform(), nil
}

// GetPlatforms returns specified platform or all platforms if unspecified
func (c *ArdiClient) GetPlatforms() ([]*rpc.Platform, error) {
	if err := c.UpdateIndexFiles(); err != nil {
		return nil, err
	}

	searchResp, err := c.client.PlatformSearch(
		context.Background(),
		&rpc.PlatformSearchReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Platform search error")
		return nil, err
	}

	return searchResp.GetSearchOutput(), nil
}

// ConnectedBoards returns a list of connected arduino boards
func (c *ArdiClient) ConnectedBoards() []*Board {
	boardList := []*Board{}

	boardListResp, err := c.client.BoardList(
		context.Background(),
		&rpc.BoardListReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Board list error")
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
func (c *ArdiClient) AllBoards() []*Board {
	c.logger.Debug("Getting list of supported boards...")

	boardList := []*Board{}

	listResp, err := c.client.PlatformList(
		context.Background(),
		&rpc.PlatformListReq{
			Instance: c.instance,
		},
	)

	if err != nil {
		c.logger.WithError(err).Error("Failed to get board list")
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
func (c *ArdiClient) Upload(fqbn, sketchDir, device string) error {
	uplRespStream, err := c.client.Upload(
		context.Background(),
		&rpc.UploadReq{
			Instance:   c.instance,
			Fqbn:       fqbn,
			SketchPath: sketchDir,
			Port:       device,
			Verbose:    c.isVerbose(),
		})

	if err != nil {
		c.logger.WithError(err).Error("Failed to upload")
		return err
	}

	for {
		uplResp, err := uplRespStream.Recv()
		if err == io.EOF {
			// target.Uploading = false
			c.logger.Info("Upload complete")
			return nil
		}

		if err != nil {
			c.logger.WithError(err).Error("Failed to upload")
			// target.Uploading = false
			return err
		}

		// When an operation is ongoing you can get its output
		if resp := uplResp.GetOutStream(); resp != nil {
			c.logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := uplResp.GetErrStream(); resperr != nil {
			c.logger.Debugf("STDERR: %s", resperr)
		}
	}
}

// CompileOpts represents the options passed to the compile command
type CompileOpts struct {
	FQBN       string
	SketchDir  string
	SketchPath string
	ExportName string
	BuildProps []string
	ShowProps  bool
}

// Compile the specified sketch
func (c *ArdiClient) Compile(opts CompileOpts) error {
	exportFile := ""
	if opts.ExportName != "" {
		exportFile = path.Join(opts.SketchDir, opts.ExportName)
	}

	compRespStream, err := c.client.Compile(
		context.Background(),
		&rpc.CompileReq{
			Instance:        c.instance,
			Fqbn:            opts.FQBN,
			SketchPath:      opts.SketchPath,
			ExportFile:      exportFile,
			BuildProperties: opts.BuildProps,
			ShowProperties:  opts.ShowProps,
			Verbose:         true,
		})

	if err != nil {
		c.logger.WithError(err).Error("Failed to compile")
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
			c.logger.WithError(err).Error("Failed to compile")
			return err
		}

		// When an operation is ongoing you can get its output
		if resp := compResp.GetOutStream(); resp != nil {
			c.logger.Debugf("STDOUT: %s", resp)
		}
		if resperr := compResp.GetErrStream(); resperr != nil {
			c.logger.Errorf("STDERR: %s", resperr)
		}
	}
}

// SearchLibraries searches available libraries for download
func (c *ArdiClient) SearchLibraries(query string) ([]*rpc.SearchedLibrary, error) {
	searchResp, err := c.client.LibrarySearch(
		context.Background(),
		&rpc.LibrarySearchReq{
			Instance: c.instance,
			Query:    query,
		},
	)
	if err != nil {
		c.logger.WithError(err).Error("Error searching libraries")
		return nil, err
	}

	return searchResp.GetLibraries(), nil
}

// InstallLibrary installs specified version of a library
func (c *ArdiClient) InstallLibrary(name, version string) (string, error) {
	installRespStream, err := c.client.LibraryInstall(
		context.Background(),
		&rpc.LibraryInstallReq{
			Instance: c.instance,
			Name:     name,
			Version:  version,
		})

	if err != nil {
		c.logger.WithError(err).Error("Error installing library")
		return "", err
	}

	foundVersion := ""

	for {
		installResp, err := installRespStream.Recv()
		if err == io.EOF {
			c.logger.Info("Lib install done")
			return foundVersion, nil
		}

		if err != nil {
			c.logger.WithError(err).Error("Library install error")
			return "", err
		}

		if installResp.GetProgress() != nil {
			c.logger.Infof("DOWNLOAD: %s\n", installResp.GetProgress())
		}
		if installResp.GetTaskProgress() != nil {
			msg := installResp.GetTaskProgress()
			lib := msg.GetName()
			c.logger.Infof("TASK: %s\n", msg)
			if foundVersion == "" {
				foundVersion = strings.Split(lib, "@")[1]
			}
		}
	}
}

// UninstallLibrary removes specified library
func (c *ArdiClient) UninstallLibrary(name string) error {
	uninstallRespStream, err := c.client.LibraryUninstall(
		context.Background(),
		&rpc.LibraryUninstallReq{
			Instance: c.instance,
			// Assume spaces in name were intended to be underscore. This indicates
			// a potential bug in the arduino-cli package manager as names
			// potentially do not have a one-to-one mapping with regards to install
			// and remove commands. It seems as though arduino should be forcing
			// devs to name their library according to the github url.
			// @todo there has to be a better way - find it!
			Name: strings.ReplaceAll(name, " ", "_"),
		})

	if err != nil {
		c.logger.WithError(err).Error("Error uninstalling library")
		return err
	}

	for {
		uninstallRespStream, err := uninstallRespStream.Recv()
		if err == io.EOF {
			c.logger.Info("Lib uninstall done")
			return nil
		}

		if err != nil {
			c.logger.WithError(err).Error("Library install error")
			return err
		}

		if uninstallRespStream.GetTaskProgress() != nil {
			c.logger.Infof("TASK: %s\n", uninstallRespStream.GetTaskProgress())
		}
	}
}

// private methods
func (c *ArdiClient) isVerbose() bool {
	return c.logger.Level == log.DebugLevel
}

// private helpers
func parsePlatform(platform string) (string, string, string) {
	version := ""
	arch := ""
	parts := strings.Split(platform, "@")

	platform = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}

	platParts := strings.Split(platform, ":")
	platform = platParts[0]

	if len(platParts) > 1 {
		arch = platParts[1]
	}

	return platform, arch, version
}

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
	addr := fmt.Sprintf("localhost:%s", port)
	// Establish a connection with the gRPC server, started with the command: arduino-cli daemon
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
