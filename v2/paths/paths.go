package paths

import (
	"os"
	"path"
)

var homeDir, _ = os.UserHomeDir()

// ArdiGlobalDataDir global data directory, ~/.ardi, used to avoid polluting an
// existing arduino-cli installation. Stores cores, libraries and the likes.
var ArdiGlobalDataDir = path.Join(homeDir, ".ardi")

// ArdiDataConfig used to tell arduino-cli where to find libraries
var ArdiDataConfig = "ardi.yaml"

// ArdiBuildConfig stores library and build details for a specific project
var ArdiBuildConfig = "ardi.json"

// ArdiGlobalDataConfig returns path to global library directory config file
var ArdiGlobalDataConfig = path.Join(ArdiGlobalDataDir, ArdiDataConfig)
