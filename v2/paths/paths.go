package paths

import (
	"os"
	"path"
)

var homeDir, _ = os.UserHomeDir()

// ArdiDataDir per-project data directory for storeage of cores, libraries etc.
var ArdiDataDir = ".ardi"

// ArdiDataConfig used to tell arduino-cli where to find libraries
var ArdiDataConfig = path.Join(ArdiDataDir, "arduino-cli.yaml")

// ArdiBuildConfig stores library and build details for a specific project
var ArdiBuildConfig = "ardi.json"
