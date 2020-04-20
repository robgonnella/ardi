package types

// DataConfig represents yaml config for telling arduino-cli where to find libraries
type DataConfig struct {
	ProxyType      string                 `yaml:"proxy_type"`
	SketchbookPath string                 `yaml:"sketchbook_path"`
	ArduinoData    string                 `yaml:"arduino_data"`
	BoardManager   map[string]interface{} `yaml:"board_manager,flow"`
}
