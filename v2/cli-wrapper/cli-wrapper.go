package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/arduino/arduino-cli/cli/output"
	"github.com/arduino/arduino-cli/commands"
	"github.com/arduino/arduino-cli/configuration"
	rpc "github.com/arduino/arduino-cli/rpc/commands"
	"github.com/robgonnella/ardi/v2/types"
	log "github.com/sirupsen/logrus"
)

// Wrapper our wrapper around the arduino-cli interface
type Wrapper struct {
	ctx          context.Context
	settingsPath string
	cli          Cli
	logger       *log.Logger
}

// Board represents a single arduino Board
type Board struct {
	FQBN string
	Name string
	Port string
}

// NewCli return new arduino-cli wrapper
func NewCli(ctx context.Context, settingsPath string, svrSettings *types.ArduinoCliSettings, logger *log.Logger, cli Cli) *Wrapper {
	configuration.Settings = configuration.Init(settingsPath)
	if cli == nil {
		cli = newArduinoCli()
	}
	return &Wrapper{
		ctx:          ctx,
		settingsPath: settingsPath,
		logger:       logger,
		cli:          cli,
	}
}

// UpdateIndexFiles updates platform and library index files
func (c *Wrapper) UpdateIndexFiles() error {
	if err := c.UpdatePlatformIndex(); err != nil {
		return err
	}
	if err := c.UpdateLibraryIndex(); err != nil {
		return err
	}
	return nil
}

// UpdateLibraryIndex updates library index file
func (c *Wrapper) UpdateLibraryIndex() error {
	c.logger.Debug("Updating library index...")
	inst := c.cli.CreateInstanceIgnorePlatformIndexErrors()

	return c.cli.UpdateLibrariesIndex(
		c.ctx,
		&rpc.UpdateLibrariesIndexReq{
			Instance: inst,
		},
		c.getDownloadProgressFn(),
	)
}

// UpdatePlatformIndex updates platform index file
func (c *Wrapper) UpdatePlatformIndex() error {
	c.logger.Debug("Updating platform index...")
	inst := c.cli.CreateInstanceIgnorePlatformIndexErrors()
	_, err := c.cli.UpdateIndex(
		c.ctx,
		&rpc.UpdateIndexReq{
			Instance: inst,
		},
		c.getDownloadProgressFn(),
	)
	return err
}

// UpgradePlatform upgrades a given platform
func (c *Wrapper) UpgradePlatform(platform string) error {
	inst, err := c.cli.CreateInstance()
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
	_, err = c.cli.PlatformUpgrade(
		c.ctx,
		req,
		c.getDownloadProgressFn(),
		c.getTaskProgressFn(),
	)
	return err
}

// InstallPlatform installs a given platform
func (c *Wrapper) InstallPlatform(platform string) (string, string, error) {
	inst, err := c.cli.CreateInstance()
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

	_, err = c.cli.PlatformInstall(
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
func (c *Wrapper) UninstallPlatform(platform string) (string, error) {
	inst, err := c.cli.CreateInstance()
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

	_, err = c.cli.PlatformUninstall(
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
func (c *Wrapper) GetInstalledPlatforms() ([]*rpc.Platform, error) {
	inst, err := c.cli.CreateInstance()
	if err != nil {
		return nil, err
	}

	req := &rpc.PlatformListReq{
		Instance:      inst,
		UpdatableOnly: false,
		All:           false,
	}

	return c.cli.GetPlatforms(req)
}

// SearchPlatforms returns specified platform or all platforms if unspecified
func (c *Wrapper) SearchPlatforms() ([]*rpc.Platform, error) {
	inst, err := c.cli.CreateInstance()
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

	resp, err := c.cli.PlatformSearch(req)

	return resp.GetSearchOutput(), err
}

// ConnectedBoards returns a list of connected arduino boards
func (c *Wrapper) ConnectedBoards() []*Board {
	inst, err := c.cli.CreateInstance()
	if err != nil {
		c.logger.WithError(err).Warn("failed to get list of connected boards")
		return nil
	}

	boardList := []*Board{}

	ports, err := c.cli.ConnectedBoards(inst.GetId())
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
func (c *Wrapper) AllBoards() []*Board {
	inst, err := c.cli.CreateInstance()
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

	platforms, err := c.cli.GetPlatforms(req)
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
func (c *Wrapper) Upload(fqbn, sketchDir, device string) error {
	inst, err := c.cli.CreateInstance()
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

	_, err = c.cli.Upload(
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
func (c *Wrapper) Compile(opts CompileOpts) error {
	inst, err := c.cli.CreateInstance()
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

	_, err = c.cli.Compile(
		c.ctx,
		req,
		os.Stdout,
		os.Stderr,
		c.isVerbose(),
	)

	return err
}

// SearchLibraries searches available libraries for download
func (c *Wrapper) SearchLibraries(query string) ([]*rpc.SearchedLibrary, error) {
	inst, err := c.cli.CreateInstance()
	if err != nil {
		return nil, err
	}

	req := &rpc.LibrarySearchReq{
		Instance: inst,
		Query:    query,
	}

	searchResp, err := c.cli.LibrarySearch(c.ctx, req)

	return searchResp.GetLibraries(), err
}

// InstallLibrary installs specified version of a library
func (c *Wrapper) InstallLibrary(name, version string) (string, error) {
	inst := c.cli.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryInstallReq{
		Instance: inst,
		Name:     name,
		Version:  version,
	}

	err := c.cli.LibraryInstall(
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
func (c *Wrapper) UninstallLibrary(name string) error {
	inst := c.cli.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryUninstallReq{
		Instance: inst,
		// Assume spaces in name were intended to be underscores. This indicates
		// a potential bug in the arduino-cli package manager as names
		// potentially do not have a one-to-one mapping with regards to install
		// and remove commands. It seems as though arduino should be forcing
		// devs to name their library according to the github url.
		// @todo there has to be a better way - find it!
		Name: strings.ReplaceAll(name, " ", "_"),
	}

	err := c.cli.LibraryUninstall(
		c.ctx,
		req,
		c.getTaskProgressFn())

	return err
}

// GetInstalledLibs returns a list of installed libraries
func (c *Wrapper) GetInstalledLibs() ([]*rpc.InstalledLibrary, error) {
	inst := c.cli.CreateInstanceIgnorePlatformIndexErrors()

	req := &rpc.LibraryListReq{
		Instance: inst,
	}

	res, err := c.cli.LibraryList(c.ctx, req)
	return res.GetInstalledLibrary(), err
}

// GetTargetBoard returns target info for a connected & disconnected boards
func (c *Wrapper) GetTargetBoard(fqbn, port string, onlyConnected bool) (*Board, error) {
	if fqbn != "" && port != "" {
		return &Board{
			FQBN: fqbn,
			Port: port,
		}, nil
	}

	fqbnErr := errors.New("you must specify a board fqbn to compile - you can find a list of board fqbns for installed platforms above")
	connectedBoardsErr := errors.New("no connected boards detected")
	connectedBoards := c.ConnectedBoards()
	allBoards := c.AllBoards()

	if fqbn != "" {
		if onlyConnected {
			for _, b := range connectedBoards {
				if b.FQBN == fqbn {
					return b, nil
				}
			}
			return nil, connectedBoardsErr
		}
		return &Board{FQBN: fqbn}, nil
	}

	if len(connectedBoards) == 0 {
		if onlyConnected {
			return nil, connectedBoardsErr
		}
		c.printFQBNs(allBoards, c.logger)
		return nil, fqbnErr
	}

	if len(connectedBoards) == 1 {
		return connectedBoards[0], nil
	}

	if len(connectedBoards) > 1 {
		c.printFQBNs(connectedBoards, c.logger)
		return nil, fqbnErr
	}

	return nil, errors.New("error parsing target")
}

// ClientVersion returns version of arduino-cli
func (c *Wrapper) ClientVersion() string {
	return c.cli.Version()
}

// private methods
func (c *Wrapper) isVerbose() bool {
	return c.logger.GetLevel() == log.DebugLevel
}

func (c *Wrapper) getDownloadProgressFn() commands.DownloadProgressCB {
	if c.isVerbose() {
		return output.ProgressBar()
	}
	return noDownloadOutput
}

func (c *Wrapper) getTaskProgressFn() commands.TaskProgressCB {
	if c.isVerbose() {
		return output.TaskProgress()
	}
	return noTaskOutput
}

// private helpers
func (c *Wrapper) printFQBNs(boardList []*Board, logger *log.Logger) {
	sort.Slice(boardList, func(i, j int) bool {
		return boardList[i].Name < boardList[j].Name
	})

	c.printBoardsWithIndices(boardList, logger)
}

func (c *Wrapper) printBoardsWithIndices(boards []*Board, logger *log.Logger) {
	w := tabwriter.NewWriter(logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()
	w.Write([]byte("No.\tName\tFQBN\n"))
	for i, b := range boards {
		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\n", i, b.Name, b.FQBN)))
	}
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
