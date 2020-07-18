package core

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// LibCore core module for lib commands
type LibCore struct {
	logger      *log.Logger
	client      rpc.Client
	initialized bool
}

// NewLibCore Lib instance
func NewLibCore(client rpc.Client, logger *log.Logger) *LibCore {
	return &LibCore{
		logger:      logger,
		client:      client,
		initialized: false,
	}
}

// Search all available libraries with optional search filter
func (l *LibCore) Search(searchArg string) error {
	l.init()

	libraries, err := l.client.SearchLibraries(searchArg)
	if err != nil {
		return err
	}
	if len(libraries) == 0 {
		return fmt.Errorf("no libraries found for %s", searchArg)
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].GetName() < libraries[j].GetName()
	})

	w := tabwriter.NewWriter(l.logger.Out, 0, 0, 8, ' ', 0)
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
func (l *LibCore) Add(lib string) (string, string, error) {
	l.init()

	libParts := strings.Split(lib, "@")
	library := libParts[0]
	version := ""
	if len(libParts) > 1 {
		version = libParts[1]
	}

	installedVersion, err := l.client.InstallLibrary(library, version)
	if err != nil {
		return "", "", err
	}

	l.logger.Infof("Installed library: %s %s", library, installedVersion)
	return library, installedVersion, nil
}

// Remove library either globally or for project
func (l *LibCore) Remove(library string) error {
	l.logger.Infof("Removing library: %s", library)
	if err := l.client.UninstallLibrary(library); err != nil {
		return err
	}

	return nil
}

// ListInstalled lists all installed libraries
func (l *LibCore) ListInstalled() error {
	libs, err := l.client.GetInstalledLibs()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(l.logger.Out, 0, 0, 8, ' ', 0)
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
func (l *LibCore) init() error {
	if !l.initialized {
		if err := l.client.UpdateLibraryIndex(); err != nil {
			l.logger.WithError(err).Warn("Failed to update library index file")
			return err
		}
		l.initialized = true
	}
	return nil
}
