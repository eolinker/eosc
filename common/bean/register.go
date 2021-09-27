package bean

import (
	"fmt"
	"sync"
)

var (
	locker     = sync.RWMutex{}
	beanOfName = make(map[string]interface{})
)

func RegisterByName(namespace string, m interface{}) {
	locker.Lock()
	defer locker.Unlock()

	_, has := beanOfName[namespace]
	if has {
		panic(fmt.Sprintf("register name duplication:%s", namespace))
	}
	beanOfName[namespace] = m
}

func GetListByName(namespaces ...string) []interface{} {
	rs := make([]interface{}, 0, len(namespaces))
	locker.RLock()
	defer locker.RUnlock()

	for _, namespace := range namespaces {
		v, has := beanOfName[namespace]
		if has {
			rs = append(rs, v)
		}
	}

	return rs
}
func GetByName(namespace string) (interface{}, bool) {
	locker.RLock()
	defer locker.RUnlock()
	v, has := beanOfName[namespace]
	return v, has
}
