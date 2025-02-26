package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/handlers/cli"
	"github.com/prskr/go-dito/infrastructure/config"
)

type App struct {
	Serve cli.ServeHandler `cmd:"" name:"serve" help:"Run mock server"`

	ConfigPath string `name:"config" short:"c" default:"config.pkl" env:"DITO_CONFIG_PATH" help:"Path to config file" type:"existingfile"`
}

func (app *App) Execute(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cliCtx := kong.Parse(
		app,
		kong.Name("dito"),
		kong.Description("go-dito"),
		kong.BindTo(ports.CWD(os.DirFS(wd)), (*ports.CWD)(nil)),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)

	return cliCtx.Run()
}

func (app *App) AfterApply(ctx context.Context, kongCtx *kong.Context) error {
	appConfig, err := config.LoadFromPath(ctx, app.ConfigPath)
	if err != nil {
		return err
	}

	logger := slog.New(appConfig.Telemetry.Logging.Handler(os.Stderr))

	config.SetCurrent(appConfig)
	slog.SetDefault(logger)

	kongCtx.Bind(appConfig, logger)

	return nil
}
