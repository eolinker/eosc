package traffic

import (
	"io"
	"net"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

type tTrafficData struct {
	data eosc.IUntyped
}

func (t *tTrafficData) remove(name string) {
	t.data.Del(name)
}
func (t *tTrafficData) Del(name string) (net.Listener, bool) {
	d, has := t.data.Del(name)
	if has {
		return d.(tListener).Listener, has
	}
	return nil, false
}
func newTTrafficData() *tTrafficData {
	return &tTrafficData{
		data: eosc.NewUntyped(),
	}
}

func (t *tTrafficData) Read(r io.Reader) {

	listeners, err := readListener(r)
	log.Debug("read listeners: ", len(listeners))
	if err != nil {
		log.Warn("read listeners:", err)
		return
	}
	for _, ln := range listeners {
		t.add(ln)
	}
}
func (t *tTrafficData) add(ln net.Listener) {
	name := toName(ln)
	log.Info("traffic add:", name)
	t.data.Set(name, newTTcpListener(ln, t))
}

func (t *tTrafficData) get(name string) (net.Listener, bool) {
	d, has := t.data.Get(name)
	if has {
		return d.(*tListener), has
	}
	return nil, false

}
func (t *tTrafficData) All() map[string]net.Listener {
	all := t.data.All()
	res := make(map[string]net.Listener)
	for n, v := range all {
		res[n] = v.(*tListener).Listener
	}
	return res
}
func (t *tTrafficData) list() []interface{} {
	return t.data.List()
}
func (t *tTrafficData) clone() *tTrafficData {
	return &tTrafficData{
		data: t.data.Clone(),
	}
}
