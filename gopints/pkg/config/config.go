package config

import "time"

// Loader is the interface for loading application configuration.
type Loader interface {
	Load() (*AppConfig, error)
}

// AppConfig is the top-level configuration for both binaries.
type AppConfig struct {
	Agent  AgentConfig  `json:"agent"`
	Server ServerConfig `json:"server"`
}

// AgentConfig configures the agent binary.
type AgentConfig struct {
	Taps      []TapConfig `json:"taps"`
	ServerUDP string      `json:"server_udp"` // e.g. "localhost:9876"
	Simulate  bool        `json:"simulate"`
}

// TapConfig configures a single tap's GPIO pin.
type TapConfig struct {
	ID       uint8         `json:"id"`
	ChipName string        `json:"chip_name"`
	Pin      int           `json:"pin"`
	Edge     string        `json:"edge"`     // "rising" | "falling" | "both"
	Debounce time.Duration `json:"debounce"` // e.g. 5ms
}

// ServerConfig configures the server binary.
type ServerConfig struct {
	UDPAddr        string  `json:"udp_addr"`         // e.g. ":9876"
	HTTPAddr       string  `json:"http_addr"`        // e.g. ":8080"
	DBPath         string  `json:"db_path"`          // e.g. "kegerator.db"
	PulsesPerLiter float64 `json:"pulses_per_liter"` // flow sensor calibration
	Simulate       bool    `json:"simulate"`
}
