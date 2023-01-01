package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/arduino/arduino-cli/cli/output"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	log "github.com/sirupsen/logrus"
)

// Wrapper our wrapper around the arduino-cli interface
type Wrapper struct {
	ctx          context.Context
	cli          Cli
	inst         *rpc.Instance
	settingsPath string
	logger       *log.Logger
}

// BoardWithPort represents a single arduino Board with associated port
type BoardWithPort struct {
	FQBN string
	Name string
	Port string
}

// WrapperOption represents and option for the wrapper
type WrapperOption = func(w *Wrapper)

// NewCli return new arduino-cli wrapper
func NewCli(ctx context.Context, settingsPath string, logger *log.Logger, options ...WrapperOption) *Wrapper {
	w := &Wrapper{
		ctx:          ctx,
		logger:       logger,
		settingsPath: settingsPath,
	}

	for _, o := range options {
		o(w)
	}

	return w
}

// WithArduinoCli allows an injectable arduino cli interface
func WithArduinoCli(arduinoCli Cli) WrapperOption {
	return func(w *Wrapper) {
		w.cli = arduinoCli
		w.cli.InitSettings(w.settingsPath)
	}
}

// UpdateIndexFiles updates platform and library index files
func (w *Wrapper) UpdateIndexFiles() error {
	if err := w.UpdatePlatformIndex(); err != nil {
		return err
	}
	if err := w.UpdateLibraryIndex(); err != nil {
		return err
	}
	return nil
}

// UpdateLibraryIndex updates library index file
func (w *Wrapper) UpdateLibraryIndex() error {
	w.logger.Debug("Updating library index...")
	inst := w.getRPCInstance()

	return w.cli.UpdateLibrariesIndex(
		w.ctx,
		&rpc.UpdateLibrariesIndexRequest{
			Instance: inst,
		},
		w.getDownloadProgressFn(),
	)
}

// UpdatePlatformIndex updates platform index file
func (w *Wrapper) UpdatePlatformIndex() error {
	w.logger.Debug("Updating platform index...")
	inst := w.getRPCInstance()

	err := w.cli.UpdateIndex(
		w.ctx,
		&rpc.UpdateIndexRequest{
			Instance: inst,
		},
		w.getDownloadProgressFn(),
	)
	return err
}

// InstallPlatform installs a given platform
func (w *Wrapper) InstallPlatform(platform string) (string, string, error) {
	inst := w.getRPCInstance()

	pkg, arch, version := parsePlatform(platform)
	installedPlatform := fmt.Sprintf("%s:%s", pkg, arch)

	req := &rpc.PlatformInstallRequest{
		Instance:        inst,
		PlatformPackage: pkg,
		Architecture:    arch,
		Version:         version,
	}

	_, err := w.cli.PlatformInstall(
		w.ctx,
		req,
		w.getDownloadProgressFn(),
		w.getTaskProgressFn(),
	)
	if err != nil {
		return "", "", err
	}

	platforms, err := w.GetInstalledPlatforms()
	if err != nil {
		return "", "", err
	}

	foundVersion := version

	for _, plat := range platforms {
		if plat.GetId() == fmt.Sprintf("%s:%s", pkg, arch) {
			foundVersion = plat.GetInstalled()
		}
	}

	return installedPlatform, foundVersion, nil
}

// UninstallPlatform installs a given platform
func (w *Wrapper) UninstallPlatform(platform string) (string, error) {
	inst := w.getRPCInstance()

	pkg, arch, _ := parsePlatform(platform)

	removedPlatform := fmt.Sprintf("%s:%s", pkg, arch)

	req := &rpc.PlatformUninstallRequest{
		Instance:        inst,
		PlatformPackage: pkg,
		Architecture:    arch,
	}

	_, err := w.cli.PlatformUninstall(
		w.ctx,
		req,
		w.getTaskProgressFn(),
	)

	if err != nil {
		return "", err
	}

	return removedPlatform, nil
}

// GetInstalledPlatforms lists all installed platforms
func (w *Wrapper) GetInstalledPlatforms() ([]*rpc.Platform, error) {
	inst := w.getRPCInstance()

	req := &rpc.PlatformListRequest{
		Instance:      inst,
		UpdatableOnly: false,
		All:           false,
	}

	return w.cli.GetPlatforms(req)
}

// SearchPlatforms returns specified platform or all platforms if unspecified
func (w *Wrapper) SearchPlatforms() ([]*rpc.Platform, error) {
	if err := w.UpdatePlatformIndex(); err != nil {
		return nil, err
	}

	inst := w.getRPCInstance()

	req := &rpc.PlatformSearchRequest{
		Instance:    inst,
		AllVersions: false,
	}

	resp, err := w.cli.PlatformSearch(req)

	return resp.GetSearchOutput(), err
}

// AllBoards returns a list of all supported boards
func (w *Wrapper) AllBoards() []*BoardWithPort {
	inst := w.getRPCInstance()

	w.logger.Debug("Getting list of supported boards...")

	boardList := []*BoardWithPort{}

	req := &rpc.PlatformListRequest{
		Instance:      inst,
		UpdatableOnly: false,
		All:           true,
	}

	platforms, err := w.cli.GetPlatforms(req)
	if err != nil {
		w.logger.WithError(err).Warn("failed to get list of installed platforms")
		return boardList
	}

	for _, p := range platforms {
		for _, b := range p.GetBoards() {
			b := BoardWithPort{
				FQBN: b.GetFqbn(),
				Name: b.GetName(),
			}
			boardList = append(boardList, &b)
		}
	}
	return boardList
}

// SearchLibraries searches available libraries for download
func (w *Wrapper) SearchLibraries(query string) ([]*rpc.SearchedLibrary, error) {
	inst := w.getRPCInstance()

	req := &rpc.LibrarySearchRequest{
		Instance: inst,
		Query:    query,
	}

	searchResp, err := w.cli.LibrarySearch(w.ctx, req)

	return searchResp.GetLibraries(), err
}

// InstallLibrary installs specified version of a library
func (w *Wrapper) InstallLibrary(name, version string) (string, error) {
	inst := w.getRPCInstance()

	req := &rpc.LibraryInstallRequest{
		Instance: inst,
		Name:     name,
		Version:  version,
	}

	err := w.cli.LibraryInstall(
		w.ctx,
		req,
		w.getDownloadProgressFn(),
		w.getTaskProgressFn(),
	)

	if err != nil {
		return "", err
	}

	libs, err := w.GetInstalledLibs()
	if err != nil {
		return "", err
	}

	foundVersion := version

	for _, lib := range libs {
		if lib.GetLibrary().Name == name {
			foundVersion = lib.GetLibrary().Version
		}
	}

	return foundVersion, err
}

// UninstallLibrary removes specified library
func (w *Wrapper) UninstallLibrary(name string) error {
	inst := w.getRPCInstance()

	req := &rpc.LibraryUninstallRequest{
		Instance: inst,
		Name:     name,
	}

	err := w.cli.LibraryUninstall(
		w.ctx,
		req,
		w.getTaskProgressFn())

	return err
}

// GetInstalledLibs returns a list of installed libraries
func (w *Wrapper) GetInstalledLibs() ([]*rpc.InstalledLibrary, error) {
	inst := w.getRPCInstance()

	req := &rpc.LibraryListRequest{
		Instance: inst,
	}

	res, err := w.cli.LibraryList(w.ctx, req)
	return res.GetInstalledLibraries(), err
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
func (w *Wrapper) Compile(opts CompileOpts) error {
	inst := w.getRPCInstance()

	resolvedSketchPath, err := filepath.Abs(opts.SketchPath)
	if err != nil {
		return errors.New("could not resolve sketch path")
	}

	resolvedSketchDir, err := filepath.Abs(opts.SketchDir)
	if err != nil {
		return errors.New("could not resolve sketch directory")
	}

	exportDir := path.Join(resolvedSketchDir, "build")

	req := &rpc.CompileRequest{
		Instance:        inst,
		Fqbn:            opts.FQBN,
		SketchPath:      resolvedSketchPath,
		ExportDir:       exportDir,
		BuildProperties: opts.BuildProps,
		ShowProperties:  opts.ShowProps,
		Verbose:         w.isVerbose(),
	}

	_, err = w.cli.Compile(
		w.ctx,
		req,
		os.Stdout,
		os.Stderr,
		w.getTaskProgressFn(),
		w.isVerbose(),
	)

	return err
}

// ClientVersion returns version of arduino-cli
func (w *Wrapper) ClientVersion() string {
	return w.cli.Version()
}

// private methods
func (w *Wrapper) isVerbose() bool {
	return w.logger.GetLevel() == log.DebugLevel
}

func (w *Wrapper) getDownloadProgressFn() rpc.DownloadProgressCB {
	if w.isVerbose() {
		return output.ProgressBar()
	}
	return noDownloadOutput
}

func (w *Wrapper) getTaskProgressFn() rpc.TaskProgressCB {
	if w.isVerbose() {
		return output.TaskProgress()
	}
	return noTaskOutput
}

func (w *Wrapper) getRPCInstance() *rpc.Instance {
	if w.inst == nil {
		w.inst = w.cli.CreateInstance()
	}

	return w.inst
}

// private helpers
func parsePlatform(platform string) (string, string, string) {
	if platform == "" {
		return "", "", ""
	}

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
