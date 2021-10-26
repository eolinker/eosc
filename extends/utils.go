package extends

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/env"
)

func LocalExtenderPath(group, name string) string {

	return filepath.Join(env.ExtendersDir(), "local", eosc.Version(), runtime.Version(), fmt.Sprintf("%s-%s", group, name))
}
func DecodeExtenderId(id string) (group, name string, err error) {
	vs := strings.Split(id, ":")
	if len(vs) == 2 {
		return vs[0], vs[1], nil
	}
	return "", "", fmt.Errorf("%w:%s", ErrorInvalidExtenderId, id)
}

func FormatDriverId(group, name, fac string) string {
	return fmt.Sprint(group, ":", name, ":", fac)
}
