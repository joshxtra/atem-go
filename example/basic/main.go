package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mraerino/atem-go"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	client := atem.NewClient(log, "172.22.26.50")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := client.Start(ctx)
	if err != nil {
		log.Error("Failed to start", "error", err)
		os.Exit(1)
	}
	log.Info("Connected")

	stateSer, _ := json.MarshalIndent(client.State(), "", "  ")
	fmt.Printf("State:\n%s\n", string(stateSer))

	// allow for some packets to come in
	<-time.After(time.Second * 30)
}
