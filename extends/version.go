package extends

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/eolinker/eosc"
)

type VersionInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	IsLatest    bool   `json:"is_latest"`
}

func GetAvailableVersions(group, project string) ([]*VersionInfo, error) {
	plugins, err := ExtendersRequest(group, project)
	if err != nil {
		return nil, err
	}
	//var latest *VersionInfo
	versions := make([]*VersionInfo, 0, len(plugins))
	arch := Arch()
	for _, p := range plugins {
		for _, a := range p.Arches {
			if a == arch {
				versions = append(versions, p.VersionInfo)
				break
			}
		}
	}
	return versions, nil
}

//Arch 当前架构环境，[{go版本}-{eosc版本}-{架构}]
func Arch() string {
	return fmt.Sprintf("%s-%s-%s-%s", strings.TrimPrefix(runtime.Version(), "go"), eosc.Version(), runtime.GOOS, runtime.GOARCH)
}
