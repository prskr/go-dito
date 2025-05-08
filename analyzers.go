//go:build analyzers

package main

import (
	_ "code.icb4dc0.de/prskr/bazel-golangci-lint-analyzers/contextcheck"
	_ "code.icb4dc0.de/prskr/bazel-golangci-lint-analyzers/copyloopvar"
	_ "code.icb4dc0.de/prskr/bazel-golangci-lint-analyzers/err113"
	_ "github.com/gordonklaus/ineffassign/pkg/ineffassign"
	_ "github.com/kisielk/errcheck/errcheck"
	_ "github.com/lasiar/canonicalheader"
	_ "github.com/sivchari/containedctx"
)
