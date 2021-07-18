package core

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	cli "github.com/robgonnella/ardi/v2/cli-wrapper"
)

// LibCore core module for lib commands
type LibCore struct {
	logger      *log.Logger
	cli         *cli.Wrapper
	initialized bool
}

// LibCoreOption represents options for LibCore
type LibCoreOption = func(c *LibCore)

// NewLibCore Lib instance
func NewLibCore(logger *log.Logger, options ...LibCoreOption) *LibCore {
	c := &LibCore{
		logger:      logger,
		initialized: false,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithLibCliWrapper allows an injectable cli wrapper
func WithLibCliWrapper(wrapper *cli.Wrapper) LibCoreOption {
	return func(c *LibCore) {
		c.cli = wrapper
	}
}

// Search all available libraries with optional search filter
func (c *LibCore) Search(searchArg string) error {
	c.init()

	libraries, err := c.cli.SearchLibraries(searchArg)
	if err != nil {
		return err
	}
	if len(libraries) == 0 {
		return fmt.Errorf("no libraries found for %s", searchArg)
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].GetName() < libraries[j].GetName()
	})

	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()

	w.Write([]byte("Library\tLatest\tOther Releases\n"))
	for _, lib := range libraries {
		releases := []string{}
		for _, rel := range lib.GetReleases() {
			releases = append(releases, rel.GetVersion())
		}
		sort.Slice(releases, func(i, j int) bool {
			return releases[i] > releases[j]
		})
		if len(releases) > 1 {
			releases = releases[1:]
		} else {
			releases = []string{}
		}
		if len(releases) > 4 {
			releases = releases[:4]
			releases = append(releases, "...")
		}
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", lib.GetName(), lib.GetLatest().GetVersion(), strings.Join(releases, ", "))))
	}
	return nil
}

// Add library for project
func (c *LibCore) Add(lib string) (string, string, error) {
	c.init()

	libParts := strings.Split(lib, "@")
	library := libParts[0]
	version := ""
	if len(libParts) > 1 {
		version = libParts[1]
	}

	installedVersion, err := c.cli.InstallLibrary(library, version)
	if err != nil {
		return "", "", err
	}

	c.logger.Infof("Installed library: %s %s", library, installedVersion)
	return library, installedVersion, nil
}

// Remove library from project
func (c *LibCore) Remove(library string) error {
	c.logger.Infof("Removing library: %s", library)
	if err := c.cli.UninstallLibrary(library); err != nil {
		return err
	}

	return nil
}

// ListInstalled lists all installed libraries
func (c *LibCore) ListInstalled() error {
	libs, err := c.cli.GetInstalledLibs()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(c.logger.Out, 0, 0, 8, ' ', 0)
	defer w.Flush()

	w.Write([]byte("Library\tVersion\tDescription\n"))
	for _, l := range libs {
		library := l.GetLibrary()
		name := library.GetName()
		version := library.Version
		desc := library.GetSentence()
		fields := fmt.Sprintf("%s\t%s\t%s\n", name, version, desc)
		w.Write([]byte(fields))
	}

	return nil
}

// private
func (c *LibCore) init() error {
	if !c.initialized {
		if err := c.cli.UpdateLibraryIndex(); err != nil {
			c.logger.WithError(err).Warn("Failed to update library index file")
			return err
		}
		c.initialized = true
	}
	return nil
}
