package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

var _ json.Unmarshaler = (*LogFormat)(nil)

var ErrUnknownFormat = errors.New("unknown log format")

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type LogFormat string

func (f *LogFormat) UnmarshalJSON(raw []byte) error {
	switch string(raw) {
	case "text":
		*f = LogFormatText
	case "json":
		*f = LogFormatJSON
	default:
		return fmt.Errorf("%w: %s", ErrUnknownFormat, raw)
	}

	return nil
}

func (f LogFormat) String() string {
	return string(f)
}

type Logging struct {
	AddSource bool       `json:"addSource"`
	Level     slog.Level `json:"level"`
	Format    LogFormat  `json:"format"`
}

func (l Logging) Handler(out io.Writer) slog.Handler {
	switch l.Format {
	case LogFormatText:
		return slog.NewTextHandler(out, l.Options())
	case LogFormatJSON:
		fallthrough
	default:
		return slog.NewJSONHandler(out, l.Options())
	}
}

func (l Logging) Options() *slog.HandlerOptions {
	opts := &slog.HandlerOptions{
		AddSource: l.AddSource,
		Level:     l.Level,
	}

	return opts
}
