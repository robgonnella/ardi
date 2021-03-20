package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/arduino/arduino-cli/cli/globals"
	"github.com/arduino/arduino-cli/cli/instance"
	"github.com/arduino/arduino-cli/cli/output"
	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/commands/board"
	"github.com/arduino/arduino-cli/commands/compile"
	"github.com/arduino/arduino-cli/commands/core"
	"github.com/arduino/arduino-cli/commands/lib"
	"github.com/arduino/arduino-cli/commands/upload"
	"github.com/arduino/arduino-cli/configuration"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
)

// Client reprents our wrapper around arduino-cli
//go:generate mockgen -destination=../mocks/mock_cli.go -package=mocks github.com/robgonnella/ardi/v2/cli-wrapper Client
type Client interface {
	UpdateIndexFiles() error
	UpdateLibraryIndex() error
	UpdatePlatformIndex() error
	UpgradePlatform(platform string) error
	InstallPlatform(platform string) (string, string, error)
	UninstallPlatform(platform string) (string, error)
	GetInstalledPlatforms() ([]*rpc.Platform, error)
	GetPlatforms() ([]*rpc.Platform, error)
	ConnectedBoards() []*Board
	AllBoards() []*Board
	Upload(fqbn, sketchDir, device string) error
	Compile(o CompileOpts) error
	SearchLibraries(query string) ([]*rpc.SearchedLibrary, error)
	InstallLibrary(name, version string) (string, error)
	UninstallLibrary(name string) error
	GetInstalledLibs() ([]*rpc.InstalledLibrary, error)
	ClientVersion() string
}

// ArdiClient represents a client connection to arduino-cli grpc daemon
type ArdiClient struct {
	ctx          context.Context
	settingsPath string
	logger       *log.Logger
}

// Board represents a single arduino Board
type Board struct {
	FQBN string
	Name string
	Port string
}

// NewClient return new RPC controller
func NewClient(ctx context.Context, settingsPath string, svrSettings *types.ArduinoCliSettings, logger *log.Logger) Client {
	configuration.Settings = configuration.Init(settingsPath)
	return &ArdiClient{
		ctx:          ctx,
		settingsPath: settingsPath,
		logger:       logger,
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
	inst := instance.CreateInstanceIgnorePlatformIndexErrors()

	return commands.UpdateLibrariesIndex(
		c.ctx,
		&rpc.UpdateLibrariesIndexReq{
			Instance: inst,
		},
		c.getDownloadProgressFn(),
	)
}

// UpdatePlatformIndex updates platform index file
func (c *ArdiClient) UpdatePlatformIndex() error {
	c.logger.Debug("Updating platform index...")
	inst := instance.CreateInstanceIgnorePlatformIndexErrors()
	_, err := commands.UpdateIndex(
		c.ctx,
		&rpc.UpdateIndexReq{
			Instance: inst,
		},
		c.getDownloadProgressFn(),
	)
	return err
}

// UpgradePlatform upgrades a given platform
func (c *ArdiClient) UpgradePlatform(platform string) error {
	inst, err := instance.CreateInstance()
	if err != nil {
		return err
	}

	pkg, arch, _ := parsePlatform(platform)
	c.logger.Debugf("Upgrading platform: %s:%s\n", pkg, arch)
	req := &rpc.PlatformUpgradeReq{
		Instance:        inst,
		PlatformPackage: pkg,
		Architecture:    arch,
	}
	_, err = core.PlatformUpgrade(
		c.ctx,
		req,
		c.getDownloadProgressFn(),
		c.getTaskProgressFn(),
	)
	return err
}

// InstallPlatform installs a given platform
func (c *ArdiClient) InstallPlatform(platform string) (string, string, error) {
	inst, err := instance.CreateInstance()
	if err != nil {
		return "", "", err
	}

	if platform == "" {
		err := errors.New("must specify a platform to install")
		c.logger.WithError(err).Error()
		return "", "", err
	}

	pkg, arch, version := parsePlatform(platform)
	installedPlatform := fmt.Sprintf("%s:%s", pkg, arch)

	req := &rpc.PlatformInstallReq{
		Instance:        inst,
		PlatformPackage: pkg,
		Architecture:    arch,
		Version:         version,
	}

	_, err = core.PlatformInstall(
		c.ctx,
		req,
		c.getDownloadProgressFn(),
		c.getTaskProgressFn(),
	)
	if err != nil {
		return "", "", err
	}

	platforms, err := c.GetInstalledPlatforms()
	if err != nil {
		return "", "", err
	}

	foundVersion := version

	for _, plat := range platforms {
		if plat.GetID() == fmt.Sprintf("%s:%s", pkg, arch) {
			foundVersion = plat.GetInstalled()
		}
	}

	return installedPlatform, foundVersion, nil
}

// UninstallPlatform installs a given platform
func (c *ArdiClient) UninstallPlatform(platform string) (string, error) {
	inst, err := instance.CreateInstance()
	if err != nil {
		return "", err
	}

	if platform == "" {
		err := errors.New("must specify a platform to install")
		c.logger.WithError(err).Error()
		return "", err
	}

	pkg, arch, _ := parsePlatform(platform)

	removedPlatform := fmt.Sprintf("%s:%s", pkg, arch)

	req := &rpc.PlatformUninstallReq{
		Instance:        inst,
		PlatformPackage: pkg,
		Architecture:    arch,
	}

	_, err = core.PlatformUninstall(
		c.ctx,
		req,
		output.NewTaskProgressCB(),
	)

	if err != nil {
		return "", err
	}

	return removedPlatform, nil
}

// GetInstalledPlatforms lists all installed platforms
func (c *ArdiClient) GetInstalledPlatforms() ([]*rpc.Platform, error) {
	inst, err := instance.CreateInstance()
	if err != nil {
		return nil, err
	}

	req := &rpc.PlatformListReq{
		Instance:      inst,
		UpdatableOnly: false,
		All:           false,
	}

	return core.GetPlatforms(req)
}

// GetPlatforms returns specified platform or all platforms if unspecified
func (c *ArdiClient) GetPlatforms() ([]*rpc.Platform, error) {
	inst, err := instance.CreateInstance()
	if err != nil {
		return nil, err
	}

	if err := c.UpdateIndexFiles(); err != nil {
		return nil, err
	}

	req := &rpc.PlatformSearchReq{
		Instance:    inst,
		AllVersions: true,
	}

	resp, err := core.PlatformSearch(req)

	return resp.GetSearchOutput(), err
}

// ConnectedBoards returns a list of connected arduino boards
func (c *ArdiClient) ConnectedBoards() []*Board {
	inst, err := instance.CreateInstance()
	if err != nil {
		c.logger.WithError(err).Warn("failed to get list of connected boards")
		return nil
	}

	boardList := []*Board{}

	ports, err := board.List(inst.GetId())
	if err != nil {
		return nil
	}

	for _, port := range ports {
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
	inst, err := instance.CreateInstance()
	if err != nil {
		return nil
	}

	c.logger.Debug("Getting list of supported boards...")

	boardList := []*Board{}

	req := &rpc.PlatformListReq{
		Instance:      inst,
		UpdatableOnly: false,
		All:           true,
	}

	platforms, err := core.GetPlatforms(req)
	if err != nil {
		c.logger.WithError(err).Warn("failed to get list of installed platforms")
		return nil
	}

	for _, p := range platforms {
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
	inst, err := instance.CreateInstance()
	if err != nil {
		return err
	}

	req := &rpc.UploadReq{
		Instance:   inst,
		Fqbn:       fqbn,
		SketchPath: sketchDir,
		Port:       device,
		Verbose:    c.isVerbose(),
	}

	_, err = upload.Upload(
		c.ctx,
		req,
		os.Stdout,
		os.Stderr,
	)

	return err
}

// CompileOpts represents the options passed to the compile command
type CompileOpts struct {
	FQBN       string
	SketchDir  string
	SketchPath string
	BuildProps []string
	ShowProps  bool
}

// Compile the specified sketch
func (c *ArdiClient) Compile(opts CompileOpts) error {
	inst, err := instance.CreateInstance()
	if err != nil {
		return err
	}

	exportDir := path.Join(opts.SketchDir, "build")

	req := &rpc.CompileReq{
		Instance:        inst,
		Fqbn:            opts.FQBN,
		SketchPath:      opts.SketchPath,
		ExportDir:       exportDir,
		BuildProperties: opts.BuildProps,
		ShowProperties:  opts.ShowProps,
		Verbose:         c.isVerbose(),
	}

	_, err = compile.Compile(
		c.ctx,
		req,
		os.Stdout,
		os.Stderr,
		c.isVerbose(),
	)

	return err
}

// SearchLibraries searches available libraries for download
func (c *ArdiClient) SearchLibraries(query string) ([]*rpc.SearchedLibrary, error) {
	inst, err := instance.CreateInstance()
	if err != nil {
		return nil, err
	}

	req := &rpc.LibrarySearchReq{
		Instance: inst,
		Query:    query,
	}

	searchResp, err := lib.LibrarySearch(c.ctx, req)

	return searchResp.GetLibraries(), err
}

// InstallLibrary installs specified version of a library
func (c *ArdiClient) InstallLibrary(name, version string) (string, error) {
	inst := instance.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryInstallReq{
		Instance: inst,
		Name:     name,
		Version:  version,
	}

	err := lib.LibraryInstall(
		c.ctx,
		req,
		c.getDownloadProgressFn(),
		c.getTaskProgressFn(),
	)

	if err != nil {
		return "", err
	}

	libs, err := c.GetInstalledLibs()
	if err != nil {
		return "", err
	}

	foundVersion := version

	for _, lib := range libs {
		if lib.GetLibrary().Name == strings.ReplaceAll(name, " ", "_") {
			foundVersion = lib.GetLibrary().Version
		}
	}

	return foundVersion, err
}

// UninstallLibrary removes specified library
func (c *ArdiClient) UninstallLibrary(name string) error {
	inst := instance.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryUninstallReq{
		Instance: inst,
		// Assume spaces in name were intended to be underscore. This indicates
		// a potential bug in the arduino-cli package manager as names
		// potentially do not have a one-to-one mapping with regards to install
		// and remove commands. It seems as though arduino should be forcing
		// devs to name their library according to the github url.
		// @todo there has to be a better way - find it!
		Name: strings.ReplaceAll(name, " ", "_"),
	}

	err := lib.LibraryUninstall(
		c.ctx,
		req,
		c.getTaskProgressFn())

	return err
}

// GetInstalledLibs returns a list of installed libraries
func (c *ArdiClient) GetInstalledLibs() ([]*rpc.InstalledLibrary, error) {
	inst := instance.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryListReq{
		Instance: inst,
	}

	res, err := lib.LibraryList(c.ctx, req)
	return res.GetInstalledLibrary(), err
}

// ClientVersion returns version of arduino-cli
func (c *ArdiClient) ClientVersion() string {
	return globals.VersionInfo.String()
}

// private methods
func (c *ArdiClient) isVerbose() bool {
	return c.logger.GetLevel() == log.DebugLevel
}

func (c *ArdiClient) getDownloadProgressFn() commands.DownloadProgressCB {
	if c.isVerbose() {
		return output.ProgressBar()
	}
	return noDownloadOutput
}

func (c *ArdiClient) getTaskProgressFn() commands.TaskProgressCB {
	if c.isVerbose() {
		return output.TaskProgress()
	}
	return noTaskOutput
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

func noDownloadOutput(msg *rpc.DownloadProgress) {
	// do nothing
}

func noTaskOutput(msg *rpc.TaskProgress) {
	// do nothing
}
