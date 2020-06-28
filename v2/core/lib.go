package core

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	"github.com/robgonnella/ardi/v2/rpc"
)

// LibCore core module for lib commands
type LibCore struct {
	logger *log.Logger
	client rpc.Client
}

// NewLibCore Lib instance
func NewLibCore(client rpc.Client, logger *log.Logger) *LibCore {
	return &LibCore{
		logger: logger,
		client: client,
	}
}

// Search all available libraries with optional search filter
func (l *LibCore) Search(searchArg string) error {
	libraries, err := l.client.SearchLibraries(searchArg)
	if err != nil {
		return err
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].GetName() < libraries[j].GetName()
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Library\tLatest\tOther Releases")
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
		fmt.Fprintf(w, "%s\t%s\t%s\n", lib.GetName(), lib.GetLatest().GetVersion(), strings.Join(releases, ", "))
	}
	return nil
}

// Add library for project
func (l *LibCore) Add(lib string) (string, string, error) {
	libParts := strings.Split(lib, "@")
	library := libParts[0]
	version := ""
	if len(libParts) > 1 {
		version = libParts[1]
	}

	l.logger.Infof("Installing library: %s %s", library, version)

	installedVersion, err := l.client.InstallLibrary(library, version)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to install %s", library)
		return "", "", err
	}

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
