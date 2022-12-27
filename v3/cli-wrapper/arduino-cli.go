package cli

import (
	"context"

	"github.com/arduino/arduino-cli/cli/globals"
	"github.com/arduino/arduino-cli/cli/instance"
	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/commands/core"
	"github.com/arduino/arduino-cli/commands/lib"
	"github.com/arduino/arduino-cli/configuration"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
)

// Cli reprents our wrapper around arduino-cli
//go:generate mockgen -destination=../mocks/mock_cli.go -package=mocks github.com/robgonnella/ardi/v3/cli-wrapper Cli
type Cli interface {
	InitSettings(string)
	CreateInstance() *rpc.Instance
	UpdateIndex(context.Context, *rpc.UpdateIndexRequest, rpc.DownloadProgressCB) error
	UpdateLibrariesIndex(context.Context, *rpc.UpdateLibrariesIndexRequest, rpc.DownloadProgressCB) error
	PlatformInstall(context.Context, *rpc.PlatformInstallRequest, rpc.DownloadProgressCB, rpc.TaskProgressCB) (*rpc.PlatformInstallResponse, error)
	PlatformUninstall(context.Context, *rpc.PlatformUninstallRequest, func(curr *rpc.TaskProgress)) (*rpc.PlatformUninstallResponse, error)
	GetPlatforms(*rpc.PlatformListRequest) ([]*rpc.Platform, error)
	PlatformSearch(*rpc.PlatformSearchRequest) (*rpc.PlatformSearchResponse, error)
	LibrarySearch(context.Context, *rpc.LibrarySearchRequest) (*rpc.LibrarySearchResponse, error)
	LibraryInstall(context.Context, *rpc.LibraryInstallRequest, rpc.DownloadProgressCB, rpc.TaskProgressCB) error
	LibraryUninstall(context.Context, *rpc.LibraryUninstallRequest, rpc.TaskProgressCB) error
	LibraryList(context.Context, *rpc.LibraryListRequest) (*rpc.LibraryListResponse, error)
	Version() string
}

// ArduinoCli represents our wrapper around arduino-cli
type ArduinoCli struct{}

// NewArduinoCli returns a new instance of ArduinoCli
func NewArduinoCli() *ArduinoCli {
	return &ArduinoCli{}
}

// InitSettings initializes settings from the path to arduino-cli.yaml
func (c *ArduinoCli) InitSettings(settingsPath string) {
	configuration.Settings = configuration.Init(settingsPath)
}

// CreateInstance wrapper around arduino-cli CreateInstance
func (c *ArduinoCli) CreateInstance() *rpc.Instance {
	return instance.CreateAndInit()
}

// UpdateIndex wrapper around arduino-cli UpdateIndex
func (c *ArduinoCli) UpdateIndex(ctx context.Context, req *rpc.UpdateIndexRequest, fn rpc.DownloadProgressCB) error {
	return commands.UpdateIndex(ctx, req, fn)
}

// UpdateLibrariesIndex wrapper around arduino-cli UpdateLibrariesIndex
func (c *ArduinoCli) UpdateLibrariesIndex(ctx context.Context, req *rpc.UpdateLibrariesIndexRequest, fn rpc.DownloadProgressCB) error {
	return commands.UpdateLibrariesIndex(ctx, req, fn)
}

// PlatformInstall wrapper around arduino-cli PlatformInstall
func (c *ArduinoCli) PlatformInstall(ctx context.Context, req *rpc.PlatformInstallRequest, dlfn rpc.DownloadProgressCB, tfn rpc.TaskProgressCB) (*rpc.PlatformInstallResponse, error) {
	return core.PlatformInstall(ctx, req, dlfn, tfn)
}

// PlatformUninstall wrapper around arduino-cli PlatformUninstall
func (c *ArduinoCli) PlatformUninstall(ctx context.Context, req *rpc.PlatformUninstallRequest, fn func(curr *rpc.TaskProgress)) (*rpc.PlatformUninstallResponse, error) {
	return core.PlatformUninstall(ctx, req, fn)
}

// GetPlatforms wrapper around arduino-cli GetPlatforms
func (c *ArduinoCli) GetPlatforms(req *rpc.PlatformListRequest) ([]*rpc.Platform, error) {
	return core.GetPlatforms(req)
}

// PlatformSearch wrapper around arduino-cli PlatformSearch
func (c *ArduinoCli) PlatformSearch(req *rpc.PlatformSearchRequest) (*rpc.PlatformSearchResponse, error) {
	return core.PlatformSearch(req)
}

// LibrarySearch wrapper around arduino-cli LibrarySearch
func (c *ArduinoCli) LibrarySearch(ctx context.Context, req *rpc.LibrarySearchRequest) (*rpc.LibrarySearchResponse, error) {
	return lib.LibrarySearch(ctx, req)
}

// LibraryInstall wrapper around arduino-cli LibraryInstall
func (c *ArduinoCli) LibraryInstall(ctx context.Context, req *rpc.LibraryInstallRequest, dlfn rpc.DownloadProgressCB, tfn rpc.TaskProgressCB) error {
	return lib.LibraryInstall(ctx, req, dlfn, tfn)
}

// LibraryUninstall wrapper around arduino-cli LibraryUninstall
func (c *ArduinoCli) LibraryUninstall(ctx context.Context, req *rpc.LibraryUninstallRequest, tfn rpc.TaskProgressCB) error {
	return lib.LibraryUninstall(ctx, req, tfn)
}

// LibraryList wrapper around arduino-cli LibraryList
func (c *ArduinoCli) LibraryList(ctx context.Context, req *rpc.LibraryListRequest) (*rpc.LibraryListResponse, error) {
	return lib.LibraryList(ctx, req)
}

// Version wrapper around arduino-cli global version
func (c *ArduinoCli) Version() string {
	return globals.VersionInfo.String()
}
