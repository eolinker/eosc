package process_master

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/etcd"
	"net/url"
)

type EtcdSender struct {
	etcd.Etcd
}

func (e *EtcdSender) Send(event string, namespace string, key string, data []byte) error {
	dataKey := fmt.Sprintf("/%s/%s", namespace, url.PathEscape(key))
	switch event {
	case eosc.EventSet:
		return e.Etcd.Put(dataKey, data)
	case eosc.EventDel:
		return e.Etcd.Delete(dataKey)
	}
	return nil
}

func NewEtcdSender(etcd etcd.Etcd) *EtcdSender {
	return &EtcdSender{Etcd: etcd}
}
