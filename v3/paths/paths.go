package paths

import (
	"path"
	"path/filepath"
)

// arduinoCliDataDir data directory name
const arduinoCliDataDir = ".ardi"

// arduinoCliDataConfig data directory config name
const arduinoCliDataConfig = "arduino-cli.yaml"

// ardi config name
const ardiConfig = "ardi.json"

// ArdiProjectConfig per-project ardi config
var ArdiProjectConfig, _ = filepath.Abs(path.Join(".", ardiConfig))

// ArduinoCliProjectDataDir per-project data config directory for cores, libraries etc
var ArduinoCliProjectDataDir, _ = filepath.Abs(path.Join(".", arduinoCliDataDir))

// ArduinoCliProjectConfig per-project arduino-cli config
var ArduinoCliProjectConfig, _ = filepath.Abs(path.Join(ArduinoCliProjectDataDir, arduinoCliDataConfig))
