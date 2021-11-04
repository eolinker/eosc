package extends

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/env"
)

func LocalExtenderPath(group, project, version string) string {
	fileName := FormatFileName(group, project, version)

	return filepath.Join(env.ExtendersDir(), "repository", eosc.Version(), runtime.Version(), group, project, version, fileName)
}
func DecodeExtenderId(id string) (group, project, version string, err error) {
	vs := strings.Split(id, ":")
	if len(vs) == 2 {
		return vs[0], vs[1], "", nil
	}
	if len(vs) == 3 {
		return vs[0], vs[1], vs[2], nil
	}
	return "", "", "", fmt.Errorf("%w:%s", ErrorInvalidExtenderId, id)
}

func FormatDriverId(group, project, fac string) string {
	return fmt.Sprint(group, ":", project, ":", fac)
}

func FormatFileName(group, project, version string) string {
	return fmt.Sprint(group, "-", project, "-", version, "-", runtime.Version(), "-", eosc.Version(), "-", runtime.GOOS, "-", runtime.GOARCH)
}

func LocalExtendTarPath(group, project, version string) string {
	return filepath.Join(env.ExtendersDir(), "repository", eosc.Version(), runtime.Version(), group, project, version, FormatFileName(group, project, version), tarSuffix)
}
