package process_helper

import (
	"testing"

	"github.com/eolinker/eosc/service"
)

func TestProcess(t *testing.T) {
	es := []*service.ExtendsBasicInfo{
		{Group: "incloud", Project: "goku", Version: "2.7.0"},
	}
	getExtenders(es)
}
