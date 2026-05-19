package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/rickcern44/gopints/pkg/config"
)

// --- Default ---

func TestDefault_KeyValues(t *testing.T) {
	cfg := config.Default()
	if cfg.Agent.ServerUDP != "localhost:9876" {
		t.Errorf("ServerUDP: want localhost:9876, got %q", cfg.Agent.ServerUDP)
	}
	if !cfg.Agent.Simulate {
		t.Error("default Simulate should be true")
	}
	if cfg.Server.PulsesPerLiter != 450 {
		t.Errorf("PulsesPerLiter: want 450, got %f", cfg.Server.PulsesPerLiter)
	}
	if cfg.Server.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr: want :8080, got %q", cfg.Server.HTTPAddr)
	}
	if len(cfg.Agent.Taps) == 0 {
		t.Error("default config should have at least one tap")
	}
}

// --- FileLoader ---

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestFileLoader_ValidJSON(t *testing.T) {
	path := writeTempConfig(t, map[string]any{
		"server": map[string]any{
			"http_addr":        ":9090",
			"pulses_per_liter": 500,
		},
	})

	cfg, err := (&config.FileLoader{Path: path}).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.HTTPAddr != ":9090" {
		t.Errorf("HTTPAddr: want :9090, got %q", cfg.Server.HTTPAddr)
	}
	if cfg.Server.PulsesPerLiter != 500 {
		t.Errorf("PulsesPerLiter: want 500, got %f", cfg.Server.PulsesPerLiter)
	}
	// Fields not in JSON should retain defaults.
	if cfg.Server.UDPAddr != ":9876" {
		t.Errorf("UDPAddr should default to :9876, got %q", cfg.Server.UDPAddr)
	}
}

func TestFileLoader_MissingFile(t *testing.T) {
	_, err := (&config.FileLoader{Path: "/nonexistent/path/config.json"}).Load()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFileLoader_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad*.json")
	f.WriteString("{not valid json}")
	f.Close()

	_, err := (&config.FileLoader{Path: f.Name()}).Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// --- EnvLoader ---

func TestEnvLoader_ServerUDP(t *testing.T) {
	t.Setenv("KEGERATOR_SERVER_UDP", "192.168.1.10:9999")
	cfg, err := (&config.EnvLoader{}).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Agent.ServerUDP != "192.168.1.10:9999" {
		t.Errorf("ServerUDP: want 192.168.1.10:9999, got %q", cfg.Agent.ServerUDP)
	}
}

func TestEnvLoader_Simulate(t *testing.T) {
	t.Setenv("KEGERATOR_SIMULATE", "false")
	cfg, err := (&config.EnvLoader{}).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Agent.Simulate {
		t.Error("Simulate should be false")
	}
	if cfg.Server.Simulate {
		t.Error("Server.Simulate should be false")
	}
}

func TestEnvLoader_PulsesPerLiter(t *testing.T) {
	t.Setenv("KEGERATOR_PULSES_PER_LITER", "600")
	cfg, err := (&config.EnvLoader{}).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.PulsesPerLiter != 600 {
		t.Errorf("PulsesPerLiter: want 600, got %f", cfg.Server.PulsesPerLiter)
	}
}

func TestEnvLoader_InvalidSimulate(t *testing.T) {
	t.Setenv("KEGERATOR_SIMULATE", "notabool")
	if _, err := (&config.EnvLoader{}).Load(); err == nil {
		t.Fatal("expected error for invalid bool, got nil")
	}
}

func TestEnvLoader_InvalidPulsesPerLiter(t *testing.T) {
	t.Setenv("KEGERATOR_PULSES_PER_LITER", "not-a-number")
	if _, err := (&config.EnvLoader{}).Load(); err == nil {
		t.Fatal("expected error for invalid float, got nil")
	}
}

// --- StaticLoader ---

func TestStaticLoader_Passthrough(t *testing.T) {
	in := &config.AppConfig{
		Agent: config.AgentConfig{
			ServerUDP: "10.0.0.1:1234",
			Simulate:  false,
			Taps: []config.TapConfig{
				{ID: 3, ChipName: "gpiochip1", Pin: 22, Edge: "falling", Debounce: 10 * time.Millisecond},
			},
		},
	}
	got, err := (&config.StaticLoader{Config: in}).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != in {
		t.Error("StaticLoader should return the exact config pointer provided")
	}
}
