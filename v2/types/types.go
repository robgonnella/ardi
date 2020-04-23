package types

// BoardManager board_manager config for arduino-cli grpc server
type BoardManager struct {
	AdditionalUrls []string `yaml:"additional_urls"`
}

// Directories paths where arduino-cli grpc server can find data
type Directories struct {
	Data      string `yaml:"data"`
	Downloads string `yaml:"downloads"`
	User      string `yaml:"user"`
}

// Telemetry enable/disable flag for arduino-cli grpc server
type Telemetry struct {
	Enabled bool `yaml:"enabled"`
}

// DataConfig represents yaml config for telling arduino-cli where to find libraries
type DataConfig struct {
	BoardManager BoardManager `yaml:"board_manager"`
	Directories  Directories  `yaml:"directories"`
	Telemetry    Telemetry    `yaml:"telemetry"`
}

// ArdiBuildJSON represents the build properties in ardi.json
type ArdiBuildJSON struct {
	Path  string            `json:"path"`
	FQBN  string            `json:"fqbn"`
	Props map[string]string `json:"props"`
}

// ArdiConfig represents the ardi.json file
type ArdiConfig struct {
	Libraries map[string]string        `json:"libraries"`
	Builds    map[string]ArdiBuildJSON `json:"builds"`
}
