package paths

import (
	"os"
	"path"
	"path/filepath"
)

var homeDir, _ = os.UserHomeDir()

// arduinoCliDataDir data directory name
const ardiDataDir = ".ardi"

// arduinoCliDataConfig data directory config name
const arduinoCliDataConfig = "arduino-cli.yaml"

// ardi config name
const ardiConfig = "ardi.json"

// ArdiProjectConfig per-project ardi config
var ArdiProjectConfig, _ = filepath.Abs(path.Join(".", ardiConfig))

// ArdiProjectDataDir per-project data config directory for cores, libraries etc
var ArdiProjectDataDir, _ = filepath.Abs(path.Join(".", ardiDataDir))

// ArduinoCliProjectConfig per-project arduino-cli config
var ArduinoCliProjectConfig, _ = filepath.Abs(path.Join(ArdiProjectDataDir, arduinoCliDataConfig))

// ArdiGlobalDataDir global data directory for storage of cores, libaraires etc.
var ArdiGlobalDataDir, _ = filepath.Abs(path.Join(homeDir, ardiDataDir))

// ArduinoCliGlobalConfig used to configure the global data directory
var ArduinoCliGlobalConfig, _ = filepath.Abs(path.Join(ArdiGlobalDataDir, arduinoCliDataConfig))

// ArdiGlobalConfig global ardi config
var ArdiGlobalConfig, _ = filepath.Abs(path.Join(ArdiGlobalDataDir, ardiConfig))
