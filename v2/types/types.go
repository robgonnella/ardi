package types

// Project represents and arduino project
type Project struct {
	Directory string
	Sketch    string
	Baud      int
}

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

// ArduinoCliSettings represents yaml config for arduino-cli daemon
type ArduinoCliSettings struct {
	BoardManager BoardManager `yaml:"board_manager"`
	Daemon       Daemon       `yaml:"daemon"`
	Directories  Directories  `yaml:"directories"`
	Logging      Logging      `yaml:"logging"`
	Telemetry    Telemetry    `yaml:"telemetry"`
}

// ArdiBuild represents the build properties in ardi.json
type ArdiBuild struct {
	Directory string            `json:"directory"`
	Sketch    string            `json:"sketch"`
	Baud      int               `json:"baud"`
	FQBN      string            `json:"fqbn"`
	Props     map[string]string `json:"props"`
}

// ArdiDaemonConfig represents daemon config in ardi.json
type ArdiDaemonConfig struct {
	Port     string `json:"port"`
	LogLevel string `json:"logLevel"`
}

// ArdiConfig represents the ardi.json file
type ArdiConfig struct {
	Daemon    ArdiDaemonConfig     `json:"daemon"`
	Platforms map[string]string    `json:"platforms"`
	BoardURLS []string             `json:"boardUrls"`
	Libraries map[string]string    `json:"libraries"`
	Builds    map[string]ArdiBuild `json:"builds"`
}
