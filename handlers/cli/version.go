package cli

import (
	"fmt"

	"github.com/prskr/go-dito/core/ports"
)

const DefaultVersion = "dev"

var (
	Version = DefaultVersion
	Commit  = "unknown"
	Date    = ""
)

type VersionHandler struct {
	Short bool `help:"Print only the version"`
}

func (h VersionHandler) Run(stdout ports.STDOUT) error {
	if h.Short {
		_, _ = fmt.Fprintln(stdout, Version)
		return nil
	}

	if Version == DefaultVersion {
		_, _ = fmt.Fprintln(stdout, "Version is not set. This is a development build.")
		return nil
	}

	_, _ = fmt.Fprintf(stdout,
		`%s
Commit: %s
Built at %s
`, Version, Commit, Date)

	return nil
}
