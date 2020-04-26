package paths

import (
	"os"
	"path"
	"path/filepath"
)

var homeDir, _ = os.UserHomeDir()

// ardiDataDir data directory name
const ardiDataDir = ".ardi"

// ardiDataConfig data directory config name
const ardiDataConfig = "arduino-cli.yaml"

// build config name
const ardiBuildConfig = "ardi.json"

// ArdiProjectBuildConfig per-project build config
var ArdiProjectBuildConfig, _ = filepath.Abs(path.Join(".", ardiBuildConfig))

// ArdiProjectDataDir per-project data config directory for cores, libraries etc
var ArdiProjectDataDir, _ = filepath.Abs(path.Join(".", ardiDataDir))

// ArdiProjectDataConfig per-project arduino-cli config
var ArdiProjectDataConfig, _ = filepath.Abs(path.Join(ArdiProjectDataDir, ardiDataConfig))

// ArdiGlobalDataDir global data directory for storage of cores, libaraires etc.
var ArdiGlobalDataDir, _ = filepath.Abs(path.Join(homeDir, ardiDataDir))

// ArdiGlobalDataConfig used to configure the global data directory
var ArdiGlobalDataConfig, _ = filepath.Abs(path.Join(ArdiGlobalDataDir, ardiDataConfig))
