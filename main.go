package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/prskr/go-dito/cmd"
)

var app cmd.App

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := app.Execute(ctx)
	stop()

	if err != nil {
		slog.Error("Failed to execute app", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
