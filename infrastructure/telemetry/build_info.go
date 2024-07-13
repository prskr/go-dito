package telemetry

import (
	"runtime/debug"
	"sync"
)

var (
	buildInfo    *debug.BuildInfo
	infoIsSet    bool
	readInfoOnce sync.Once
)

func initBuildInfo() {
	readInfoOnce.Do(func() {
		buildInfo, infoIsSet = debug.ReadBuildInfo()
	})
}
