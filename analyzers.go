//go:build analyzers

package main

import (
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/kisielk/errcheck"
	_ "github.com/lasiar/canonicalheader"
)
