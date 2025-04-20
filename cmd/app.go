package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"

	"github.com/prskr/go-dito/core/services/config"
	"github.com/prskr/go-dito/handlers/cli"
)

type App struct {
	Serve cli.ServeHandler `cmd:"" name:"serve" help:"Run mock server"`

	ConfigPath string `name:"config" short:"c" default:"config.yaml" env:"DITO_CONFIG_PATH" help:"Path to config file" type:"existingfile"`
}

func (app *App) Execute(ctx context.Context) error {
	cliCtx := kong.Parse(
		app,
		kong.Name("dito"),
		kong.Description("go-dito"),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)

	return cliCtx.Run()
}

func (app *App) AfterApply(kongCtx *kong.Context) error {
	appConfig, err := config.LoadFromPath(app.ConfigPath)
	if err != nil {
		return err
	}

	logger := slog.New(appConfig.Telemetry.Logging.Handler(os.Stderr))

	slog.SetDefault(logger)

	kongCtx.Bind(appConfig, logger)

	return nil
}
