package version

import (
	"encoding/json"
	"fmt"
)

var (
	APP_NAME = "geoip"

	// Version release version of the provider
	version string

	// GitCommitSHA is the Git SHA of the latest tag/release
	commit string

	// GitBranch
	branch string

	buildTime string

	// DevVersion string for the development version
	DevVersion = "dev"
)

type VersionInfo struct {
	AppName   string
	Version   string
	Branch    string
	Commit    string
	BuildTime string
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{
		AppName:   APP_NAME,
		Version:   BuildVersion(),
		Branch:    branch,
		Commit:    commit,
		BuildTime: buildTime,
	}
}

// BuildVersion returns current version of the provider
func BuildVersion() string {
	if len(version) == 0 {
		return DevVersion
	}
	return version
}

func Print() {
	b, _ := json.MarshalIndent(NewVersionInfo(), "", "\t")
	fmt.Println(string(b))
}
