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
func DecodeExtenderId(id string) (group, name, version string, err error) {
	vs := strings.Split(id, ":")
	if len(vs) == 2 {
		return vs[0], vs[1], "", nil
	}
	if len(vs) == 3 {
		return vs[0], vs[1], vs[2], nil
	}
	return "", "", "", fmt.Errorf("%w:%s", ErrorInvalidExtenderId, id)
}

func FormatDriverId(group, name, fac string) string {
	return fmt.Sprint(group, ":", name, ":", fac)
}

func FormatFileName(group, project, version string) string {
	return fmt.Sprint(group, "-", project, "-", version, "-", runtime.Version(), "-", eosc.Version(), "-", runtime.GOOS, "-", runtime.GOARCH)
}
