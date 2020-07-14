package types

// BoardManager board_manager config for arduino-cli grpc server
type BoardManager struct {
	AdditionalUrls []string `yaml:"additional_urls"`
}

// Daemon daemon configuration
type Daemon struct {
	Port string `yaml:"port"`
}

// Directories paths where arduino-cli grpc server can find data
type Directories struct {
	Data      string `yaml:"data"`
	Downloads string `yaml:"downloads"`
	User      string `yaml:"user"`
}

// Logging logging configuration
type Logging struct {
	File   string `yaml:"file"`
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

// Telemetry enable/disable flag for arduino-cli grpc server
type Telemetry struct {
	Addr    string `yaml:"addr"`
	Enabled bool   `yaml:"enabled"`
}

// DataConfig represents yaml config for telling arduino-cli where to find libraries
type DataConfig struct {
	BoardManager BoardManager `yaml:"board_manager"`
	Daemon       Daemon       `yaml:"daemon"`
	Directories  Directories  `yaml:"directories"`
	Logging      Logging      `yaml:"logging"`
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
	Platforms map[string]string        `json:"platforms"`
	BoardURLS []string                 `json:"board-urls"`
	Libraries map[string]string        `json:"libraries"`
	Builds    map[string]ArdiBuildJSON `json:"builds"`
}
