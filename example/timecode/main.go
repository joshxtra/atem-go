// Package main provides an example of requesting timecode from an ATEM switcher.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krombel/atem-go"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	client := atem.NewClient(log, "172.22.26.50")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err := client.Start(ctx)
	if err != nil {
		log.Error("Failed to start", "error", err)
		cancel()
		os.Exit(1) //nolint:gocritic // exitAfterDefer: acceptable in example code
	}
	log.Info("Connected")

	ticker := time.Tick(time.Second / 2)
	for {
		select {
		case <-ctx.Done():
			log.Error("exiting", "error", ctx.Err())
			os.Exit(1)
		case <-ticker:
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*2)
		tc, err := client.Timecode(ctx)
		if err != nil {
			log.Warn("failed to get time", "error", err)
		} else {
			fmt.Printf("Time: %#v\n", tc)
		}
		cancel()
		// stateSer, _ := json.MarshalIndent(client.State(), "", "  ")
		// fmt.Printf("State:\n%s\n", string(stateSer))
	}
}
