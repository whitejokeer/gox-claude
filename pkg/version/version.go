package version

import (
	"fmt"
	"runtime"
)

// Version will be set by the linker during build
var Version = "dev"

// Info contains version information
type Info struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	Os        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns version as string
func String() string {
	info := Get()
	return fmt.Sprintf("gox version %s %s/%s %s", 
		info.Version, info.Os, info.Arch, info.GoVersion)
}