package lib

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	"github.com/robgonnella/ardi/v2/rpc"
)

// Lib core module for lib commands
type Lib struct {
	ardiJSON *ardijson.ArdiJSON
	logger   *log.Logger
	client   *rpc.Client
}

// New Lib instance
func New(client *rpc.Client, logger *log.Logger) (*Lib, error) {
	if err := client.UpdateIndexFiles(); err != nil {
		logger.WithError(err).Error("Failed to update index files")
		return nil, err
	}

	ardiJSON, err := ardijson.New(logger)
	if err != nil {
		return nil, err
	}

	return &Lib{
		ardiJSON: ardiJSON,
		logger:   logger,
		client:   client,
	}, nil
}

// Search all available libraries with optional search filter
func (l *Lib) Search(searchArg string) error {
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
func (l *Lib) Add(libraries []string) error {
	for _, lib := range libraries {
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
			return err
		}
		if err := l.ardiJSON.AddLibrary(library, installedVersion); err != nil {
			l.logger.WithError(err).Error("Failed to update ardi.json")
			return err
		}
	}
	return nil
}

// Remove library either globally or for project
func (l *Lib) Remove(libraries []string) error {
	for _, lib := range libraries {
		l.logger.Infof("Removing library: %s", lib)
		if err := l.client.UninstallLibrary(lib); err != nil {
			return err
		}
		if err := l.ardiJSON.RemoveLibrary(lib); err != nil {
			return err
		}
	}
	return nil
}

// Install all libraries specified in ardi.json
func (l *Lib) Install() error {
	for name, version := range l.ardiJSON.Config.Libraries {
		lib := []string{fmt.Sprintf("%s@%s", name, version)}
		if err := l.Add(lib); err != nil {
			return err
		}
	}

	return nil
}
