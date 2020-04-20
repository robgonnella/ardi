package lib

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"

	ardijson "github.com/robgonnella/ardi/v2/core/ardi-json"
	"github.com/robgonnella/ardi/v2/core/rpc"
)

// Lib core module for lib commands
type Lib struct {
	ardiJSON *ardijson.ArdiJSON
	logger   *log.Logger
	RPC      *rpc.RPC
}

// New Lib instance
func New(dataConfig string, logger *log.Logger) (*Lib, error) {
	rpc, err := rpc.New(dataConfig, logger)
	if err != nil {
		return nil, err
	}

	ardiJSON, err := ardijson.New(logger)
	if err != nil {
		return nil, err
	}

	return &Lib{
		ardiJSON: ardiJSON,
		logger:   logger,
		RPC:      rpc,
	}, nil
}

// Search all available libraries with optional search filter
func (l *Lib) Search(searchArg string) error {
	libraries, err := l.RPC.SearchLibraries(searchArg)
	if err != nil {
		return err
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].GetName() < libraries[j].GetName()
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Library\tLatest\tReleases")
	for _, lib := range libraries {
		releases := []string{}
		for _, rel := range lib.GetReleases() {
			releases = append(releases, rel.GetVersion())
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", lib.GetName(), lib.GetLatest().GetVersion(), strings.Join(releases, ", "))
	}
	return nil
}

// Add library for project
func (l *Lib) Add(name, version string) error {
	l.logger.Infof("Installing library: %s %s", name, version)
	installedVersion, err := l.RPC.InstallLibrary(name, version)
	if err != nil {
		l.logger.WithError(err).Errorf("Failed to install %s", name)
		return err
	}
	return l.ardiJSON.AddLibrary(name, installedVersion)
}

// Remove library either globally or for project
func (l *Lib) Remove(name string) error {
	l.logger.Infof("Removing library: %s", name)
	if err := l.RPC.UninstallLibrary(name); err != nil {
		return err
	}
	return l.ardiJSON.RemoveLibrary(name)
}

// Install all libraries specified in ardi.json
func (l *Lib) Install() error {
	for name, version := range l.ardiJSON.Config.Libraries {
		if err := l.Add(name, version); err != nil {
			return err
		}
	}

	return nil
}
