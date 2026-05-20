package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rickcern44/gopints/internal/api"
	internudp "github.com/rickcern44/gopints/internal/udp"
	"github.com/rickcern44/gopints/pkg/config"
	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
)

var version = "dev"

func main() {
	cfgPath := flag.String("config", "", "path to JSON config file (default: env vars then built-in defaults)")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	slog.Info("starting kegerator server", "version", version, "simulate", cfg.Server.Simulate)

	store, err := tap.NewSQLiteStore(cfg.Server.DBPath)
	if err != nil {
		slog.Error("failed to open database", "path", cfg.Server.DBPath, "err", err)
		os.Exit(1)
	}
	defer store.Close()

	// Seed configured taps so they exist in the DB from the start.
	seedTaps(store, cfg.Agent.Taps)

	events := make(chan flow.PourEvent, 64)
	meters := make(map[uint8]*flow.Meter, len(cfg.Agent.Taps))
	for _, t := range cfg.Agent.Taps {
		meters[t.ID] = flow.NewMeter(t.ID, cfg.Server.PulsesPerLiter, 0, events)
	}

	receiver, err := internudp.NewReceiver(cfg.Server.UDPAddr)
	if err != nil {
		slog.Error("failed to bind UDP", "addr", cfg.Server.UDPAddr, "err", err)
		os.Exit(1)
	}
	defer receiver.Close()

	slog.Info("UDP receiver listening", "addr", cfg.Server.UDPAddr)
	go receiveLoop(receiver, meters)

	srv := api.NewServer(cfg.Server.HTTPAddr, store, meters, cfg.Server.Simulate, version)

	// Single goroutine fans out events: broadcast to WS clients and persist PourEnded.
	go eventLoop(srv, store, events)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "err", err)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
}

// receiveLoop reads UDP datagrams and routes each pulse to the correct FlowMeter.
func receiveLoop(r *internudp.Receiver, meters map[uint8]*flow.Meter) {
	for {
		msg, err := r.Receive()
		if err != nil {
			// Receiver.Close() causes a read error — treat it as shutdown.
			return
		}
		slog.Debug("pulse received", "tap_id", msg.TapID)
		if m, ok := meters[msg.TapID]; ok {
			m.HandlePulse()
		} else {
			slog.Warn("pulse for unknown tap", "tap_id", msg.TapID)
		}
	}
}

// eventLoop is the single consumer of the events channel. It broadcasts every
// event to WebSocket clients and persists PourEnded events to the store.
func eventLoop(srv *api.Server, store tap.Store, events <-chan flow.PourEvent) {
	for e := range events {
		srv.Broadcast(e)
		switch e.Type {
		case flow.PourStarted:
			slog.Info("pour started", "tap_id", e.TapID)
		case flow.PourEnded:
			slog.Info("pour ended", "tap_id", e.TapID, "volume_ml", fmt.Sprintf("%.1f", e.VolumeMl))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := store.RecordPour(ctx, tap.Pour{
				TapID:     e.TapID,
				VolumeMl:  e.VolumeMl,
				StartedAt: e.StartedAt,
				EndedAt:   e.EndedAt,
			})
			cancel()
			if err != nil {
				slog.Error("failed to record pour", "tap_id", e.TapID, "err", err)
			}
		}
	}
}

func seedTaps(store tap.Store, taps []config.TapConfig) {
	ctx := context.Background()
	for _, t := range taps {
		if err := store.EnsureTap(ctx, t.ID); err != nil {
			slog.Warn("failed to seed tap", "tap_id", t.ID, "err", err)
		}
	}
}

func loadConfig(path string) (*config.AppConfig, error) {
	if path != "" {
		return (&config.FileLoader{Path: path}).Load()
	}
	return (&config.EnvLoader{}).Load()
}
