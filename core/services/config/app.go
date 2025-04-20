package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/invopop/yaml"
)

var ErrUnsupportedConfigFormat = errors.New("unsupported config format")

type Telemetry struct {
	Logging         Logging       `json:"logging"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`
}

type ServerOptions struct {
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout"`
	ShutdownTimeout   time.Duration `json:"shutdownTimeout"`
}

type RequestOptions struct {
	MaxBodySize DataSize `json:"maxBodySize"`
}

type Server struct {
	Host           string         `json:"host"`
	Port           uint16         `json:"port"`
	ServerOptions  ServerOptions  `json:"serverOptions"`
	RequestOptions RequestOptions `json:"requestOptions"`
}

func LoadFromPath(path string) (App, error) {
	var unmarshaler func(data []byte, target any) error

	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		unmarshaler = func(data []byte, target any) error {
			return yaml.Unmarshal(data, target)
		}
	case ".json":
		unmarshaler = json.Unmarshal
	default:
		return App{}, fmt.Errorf("%w: %s", ErrUnsupportedConfigFormat, filepath.Ext(path))
	}

	configFileBytes, err := os.ReadFile(path)
	if err != nil {
		return App{}, fmt.Errorf("%w: %s", err, path)
	}

	cfg := Default()

	if err := unmarshaler(configFileBytes, &cfg); err != nil {
		return App{}, fmt.Errorf("unmarshaling app config: %w", err)
	}

	return cfg, nil
}

func Default() App {
	return App{
		Domains: make(DomainMapping),
		Server: Server{
			Host:           "0.0.0.0",
			Port:           3498,
			ServerOptions:  ServerOptions{ReadHeaderTimeout: 100 * time.Millisecond, ShutdownTimeout: 10 * time.Second},
			RequestOptions: RequestOptions{MaxBodySize: 10 * MegaByte},
		},
		Telemetry: Telemetry{
			Logging:         Logging{Level: slog.LevelInfo},
			ShutdownTimeout: 10 * time.Second,
		},
	}
}

type App struct {
	Domains   DomainMapping `json:"domains"`
	Server    Server        `json:"server"`
	Telemetry Telemetry     `json:"telemetry"`
}
