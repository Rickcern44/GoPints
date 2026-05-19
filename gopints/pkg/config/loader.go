package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Default returns a config suitable for local development (simulate mode, one tap).
func Default() *AppConfig {
	return &AppConfig{
		Agent: AgentConfig{
			ServerUDP: "localhost:9876",
			Simulate:  true,
			Taps: []TapConfig{
				{
					ID:       1,
					ChipName: "gpiochip0",
					Pin:      17,
					Edge:     "rising",
					Debounce: 5 * time.Millisecond,
				},
			},
		},
		Server: ServerConfig{
			UDPAddr:        ":9876",
			HTTPAddr:       ":8080",
			DBPath:         "kegerator.db",
			PulsesPerLiter: 450,
			Simulate:       true,
		},
	}
}

// FileLoader loads configuration from a JSON file.
type FileLoader struct {
	Path string
}

func (f *FileLoader) Load() (*AppConfig, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", f.Path, err)
	}
	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", f.Path, err)
	}
	return cfg, nil
}

// EnvLoader loads configuration from KEGERATOR_* environment variables,
// falling back to Default() for any value not set.
type EnvLoader struct{}

func (e *EnvLoader) Load() (*AppConfig, error) {
	cfg := Default()

	if v := os.Getenv("KEGERATOR_SERVER_UDP"); v != "" {
		cfg.Agent.ServerUDP = v
	}
	if v := os.Getenv("KEGERATOR_SIMULATE"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("config: KEGERATOR_SIMULATE: %w", err)
		}
		cfg.Agent.Simulate = b
		cfg.Server.Simulate = b
	}
	if v := os.Getenv("KEGERATOR_UDP_ADDR"); v != "" {
		cfg.Server.UDPAddr = v
	}
	if v := os.Getenv("KEGERATOR_HTTP_ADDR"); v != "" {
		cfg.Server.HTTPAddr = v
	}
	if v := os.Getenv("KEGERATOR_DB_PATH"); v != "" {
		cfg.Server.DBPath = v
	}
	if v := os.Getenv("KEGERATOR_PULSES_PER_LITER"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("config: KEGERATOR_PULSES_PER_LITER: %w", err)
		}
		cfg.Server.PulsesPerLiter = f
	}

	return cfg, nil
}

// StaticLoader returns a fixed config. Useful for tests.
type StaticLoader struct {
	Config *AppConfig
}

func (s *StaticLoader) Load() (*AppConfig, error) {
	return s.Config, nil
}
