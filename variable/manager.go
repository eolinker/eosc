package variable

import (
	"github.com/eolinker/eosc/workers/require"
	"sync"
)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) (map[string]string, []string)
	GetByNamespace(namespace string) (map[string]string, bool)
	Get() map[string]string
}

type Manager struct {
	// variables 变量数据
	variables      map[string]string
	requireManager require.IRequires
	locker         sync.RWMutex
}
