package cli

import (
	"context"
	"io"

	"github.com/arduino/arduino-cli/cli/globals"
	"github.com/arduino/arduino-cli/cli/instance"
	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/commands/board"
	"github.com/arduino/arduino-cli/commands/compile"
	"github.com/arduino/arduino-cli/commands/core"
	"github.com/arduino/arduino-cli/commands/lib"
	"github.com/arduino/arduino-cli/commands/upload"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
)

// Cli reprents our wrapper around arduino-cli
//go:generate mockgen -destination=../mocks/mock_cli.go -package=mocks github.com/robgonnella/ardi/v2/cli-wrapper Cli
type Cli interface {
	CreateInstance() (*rpc.Instance, error)
	CreateInstanceIgnorePlatformIndexErrors() *rpc.Instance
	UpdateIndex(context.Context, *rpc.UpdateIndexReq, commands.DownloadProgressCB) (*rpc.UpdateIndexResp, error)
	UpdateLibrariesIndex(context.Context, *rpc.UpdateLibrariesIndexReq, commands.DownloadProgressCB) error
	PlatformUpgrade(context.Context, *rpc.PlatformUpgradeReq, commands.DownloadProgressCB, commands.TaskProgressCB) (*rpc.PlatformUpgradeResp, error)
	PlatformInstall(context.Context, *rpc.PlatformInstallReq, commands.DownloadProgressCB, commands.TaskProgressCB) (*rpc.PlatformInstallResp, error)
	PlatformUninstall(context.Context, *rpc.PlatformUninstallReq, func(curr *rpc.TaskProgress)) (*rpc.PlatformUninstallResp, error)
	GetPlatforms(*rpc.PlatformListReq) ([]*rpc.Platform, error)
	PlatformSearch(*rpc.PlatformSearchReq) (*rpc.PlatformSearchResp, error)
	ConnectedBoards(instanceID int32) (r []*rpc.DetectedPort, e error)
	Upload(context.Context, *rpc.UploadReq, io.Writer, io.Writer) (*rpc.UploadResp, error)
	Compile(context.Context, *rpc.CompileReq, io.Writer, io.Writer, bool) (*rpc.CompileResp, error)
	LibrarySearch(context.Context, *rpc.LibrarySearchReq) (*rpc.LibrarySearchResp, error)
	LibraryInstall(context.Context, *rpc.LibraryInstallReq, commands.DownloadProgressCB, commands.TaskProgressCB) error
	LibraryUninstall(context.Context, *rpc.LibraryUninstallReq, commands.TaskProgressCB) error
	LibraryList(context.Context, *rpc.LibraryListReq) (*rpc.LibraryListResp, error)
	Version() string
}

// ArduinoCli represents our wrapper around arduino-cli
type ArduinoCli struct{}

func newArduinoCli() *ArduinoCli {
	return &ArduinoCli{}
}

// CreateInstance wrapper around arduino-cli CreateInstance
func (c *ArduinoCli) CreateInstance() (*rpc.Instance, error) {
	return instance.CreateInstance()
}

// CreateInstanceIgnorePlatformIndexErrors wrapper around arduino-cli CreateInstanceIgnorePlatformIndexErrors
func (c *ArduinoCli) CreateInstanceIgnorePlatformIndexErrors() *rpc.Instance {
	return instance.CreateInstanceIgnorePlatformIndexErrors()
}

// UpdateIndex wrapper around arduino-cli UpdateIndex
func (c *ArduinoCli) UpdateIndex(ctx context.Context, req *rpc.UpdateIndexReq, fn commands.DownloadProgressCB) (*rpc.UpdateIndexResp, error) {
	return commands.UpdateIndex(ctx, req, fn)
}

// UpdateLibrariesIndex wrapper around arduino-cli UpdateLibrariesIndex
func (c *ArduinoCli) UpdateLibrariesIndex(ctx context.Context, req *rpc.UpdateLibrariesIndexReq, fn commands.DownloadProgressCB) error {
	return commands.UpdateLibrariesIndex(ctx, req, fn)
}

// PlatformUpgrade wrapper around arduino-cli PlatformUpgrade
func (c *ArduinoCli) PlatformUpgrade(ctx context.Context, req *rpc.PlatformUpgradeReq, dlfn commands.DownloadProgressCB, tfn commands.TaskProgressCB) (*rpc.PlatformUpgradeResp, error) {
	return core.PlatformUpgrade(ctx, req, dlfn, tfn)
}

// PlatformInstall wrapper around arduino-cli PlatformInstall
func (c *ArduinoCli) PlatformInstall(ctx context.Context, req *rpc.PlatformInstallReq, dlfn commands.DownloadProgressCB, tfn commands.TaskProgressCB) (*rpc.PlatformInstallResp, error) {
	return core.PlatformInstall(ctx, req, dlfn, tfn)
}

// PlatformUninstall wrapper around arduino-cli PlatformUninstall
func (c *ArduinoCli) PlatformUninstall(ctx context.Context, req *rpc.PlatformUninstallReq, fn func(curr *rpc.TaskProgress)) (*rpc.PlatformUninstallResp, error) {
	return core.PlatformUninstall(ctx, req, fn)
}

// GetPlatforms wrapper around arduino-cli GetPlatforms
func (c *ArduinoCli) GetPlatforms(req *rpc.PlatformListReq) ([]*rpc.Platform, error) {
	return core.GetPlatforms(req)
}

// PlatformSearch wrapper around arduino-cli PlatformSearch
func (c *ArduinoCli) PlatformSearch(req *rpc.PlatformSearchReq) (*rpc.PlatformSearchResp, error) {
	return core.PlatformSearch(req)
}

// ConnectedBoards wrapper around arduino-cli board.List
func (c *ArduinoCli) ConnectedBoards(instanceID int32) (r []*rpc.DetectedPort, e error) {
	return board.List(instanceID)
}

// Upload wrapper around arduino-cli Upload
func (c *ArduinoCli) Upload(ctx context.Context, req *rpc.UploadReq, out io.Writer, err io.Writer) (*rpc.UploadResp, error) {
	return upload.Upload(ctx, req, out, err)
}

// Compile wrapper around arduino-cli Compile
func (c *ArduinoCli) Compile(ctx context.Context, req *rpc.CompileReq, out io.Writer, err io.Writer, verbose bool) (*rpc.CompileResp, error) {
	return compile.Compile(ctx, req, out, err, verbose)
}

// LibrarySearch wrapper around arduino-cli LibrarySearch
func (c *ArduinoCli) LibrarySearch(ctx context.Context, req *rpc.LibrarySearchReq) (*rpc.LibrarySearchResp, error) {
	return lib.LibrarySearch(ctx, req)
}

// LibraryInstall wrapper around arduino-cli LibraryInstall
func (c *ArduinoCli) LibraryInstall(ctx context.Context, req *rpc.LibraryInstallReq, dlfn commands.DownloadProgressCB, tfn commands.TaskProgressCB) error {
	return lib.LibraryInstall(ctx, req, dlfn, tfn)
}

// LibraryUninstall wrapper around arduino-cli LibraryUninstall
func (c *ArduinoCli) LibraryUninstall(ctx context.Context, req *rpc.LibraryUninstallReq, tfn commands.TaskProgressCB) error {
	return lib.LibraryUninstall(ctx, req, tfn)
}

// LibraryList wrapper around arduino-cli LibraryList
func (c *ArduinoCli) LibraryList(ctx context.Context, req *rpc.LibraryListReq) (*rpc.LibraryListResp, error) {
	return lib.LibraryList(ctx, req)
}

// Version wrapper around arduino-cli global version
func (c *ArduinoCli) Version() string {
	return globals.VersionInfo.String()
}
