package extends

import (
	"fmt"
	"runtime"
	"testing"
)

func TestReadExtenderProject(t *testing.T) {
	fmt.Println(runtime.Version())
	register, err := ReadExtenderProject("incloud", "goku", "2.2.0")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(register.All())
}
