package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rickcern44/gopints/internal/udp"
	"github.com/rickcern44/gopints/pkg/agent"
	"github.com/rickcern44/gopints/pkg/config"
	"github.com/rickcern44/gopints/pkg/protocol"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "simulate" {
		runSimulate(os.Args[2:])
		return
	}
	runAgent(os.Args[1:])
}

// runAgent is the normal operating mode: watch GPIO pins and send UDP pulses.
func runAgent(args []string) {
	fs := flag.NewFlagSet("agent", flag.ExitOnError)
	cfgPath := fs.String("config", "", "path to JSON config file")
	autoPour := fs.Bool("auto-pour", false, "fire a test pour on startup (simulate mode only)")
	fs.Parse(args)

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	slog.Info("starting kegerator agent", "version", version, "simulate", cfg.Agent.Simulate, "taps", len(cfg.Agent.Taps))

	sender, err := udp.NewSender(cfg.Agent.ServerUDP)
	if err != nil {
		slog.Error("failed to create UDP sender", "addr", cfg.Agent.ServerUDP, "err", err)
		os.Exit(1)
	}
	defer sender.Close()

	var sim *agent.SimulatorRequester
	if cfg.Agent.Simulate {
		sim = agent.NewSimulatorRequester()
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	for _, tap := range cfg.Agent.Taps {
		tap := tap

		handler := &udp.Handler{TapID: tap.ID, Sender: sender}

		var requester agent.LineRequester
		if sim != nil {
			requester = sim
		} else {
			requester = agent.DefaultRequester()
		}

		obs := agent.NewObserver(
			tap.ChipName,
			tap.Pin,
			handler,
			agent.WithEdge(edgeFromString(tap.Edge)),
			agent.WithDebounce(tap.Debounce),
			agent.WithRequester(requester),
		)

		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("listening on tap", "tap", tap.ID, "chip", tap.ChipName, "pin", tap.Pin)
			if err := obs.Listen(ctx); err != nil && err != context.Canceled {
				slog.Error("observer error", "tap", tap.ID, "err", err)
			}
		}()

		if sim != nil && *autoPour {
			go func() {
				time.Sleep(2 * time.Second)
				slog.Info("auto-pour: triggering test pour", "tap", tap.ID)
				sim.Trigger(tap.ChipName, tap.Pin, 50, 200)
			}()
		}
	}

	<-ctx.Done()
	slog.Info("shutting down agent")
	wg.Wait()
}

// runSimulate sends UDP pulse datagrams directly to the server without any
// GPIO hardware or observer. It exercises the full agent→UDP→server path.
//
// Usage:
//
//	kegerator-agent simulate [flags]
//	  -tap      uint    tap ID to simulate (default 1)
//	  -pulses   int     number of pulses per pour (default 200, ~444ml at 450 p/L)
//	  -hz       float   pulse rate in Hz (default 50)
//	  -interval duration repeat pours at this interval; 0 = run once then exit
//	  -config   string  path to JSON config file
func runSimulate(args []string) {
	fs := flag.NewFlagSet("simulate", flag.ExitOnError)
	cfgPath := fs.String("config", "", "path to JSON config file")
	tapID := fs.Uint("tap", 1, "tap ID to simulate")
	pulses := fs.Int("pulses", 200, "pulses per pour (~444ml at 450 pulses/L)")
	hz := fs.Float64("hz", 50, "pulse rate in Hz")
	interval := fs.Duration("interval", 0, "repeat interval; omit or 0 to run once")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: kegerator-agent simulate [flags]\n\n")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	sender, err := udp.NewSender(cfg.Agent.ServerUDP)
	if err != nil {
		slog.Error("failed to create UDP sender", "addr", cfg.Agent.ServerUDP, "err", err)
		os.Exit(1)
	}
	defer sender.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pour := func() {
		slog.Info("simulate: pouring", "tap", *tapID, "pulses", *pulses, "hz", *hz)
		sleep := time.Duration(float64(time.Second) / *hz)
		for i := 0; i < *pulses; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			_ = sender.Send(protocol.PulseMessage{
				TapID:     uint8(*tapID),
				Timestamp: uint64(time.Now().UnixNano()),
			})
			time.Sleep(sleep)
		}
		slog.Info("simulate: pour complete", "tap", *tapID)
	}

	pour()

	if *interval > 0 {
		ticker := time.NewTicker(*interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pour()
			}
		}
	}
}

func loadConfig(path string) (*config.AppConfig, error) {
	if path != "" {
		return (&config.FileLoader{Path: path}).Load()
	}
	return (&config.EnvLoader{}).Load()
}

func edgeFromString(s string) agent.Edge {
	switch s {
	case "falling":
		return agent.EdgeFalling
	case "both":
		return agent.EdgeBoth
	default:
		return agent.EdgeRising
	}
}
