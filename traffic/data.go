package traffic

import (
	"github.com/eolinker/eosc/log"
	"net"
)

type TrafficData struct {
	data map[string]*net.TCPListener
	stop bool
}

func (t *TrafficData) IsStop() bool {
	return t.stop
}

func NewTrafficData(data map[string]*net.TCPListener) *TrafficData {
	if data == nil {
		data = map[string]*net.TCPListener{}
	}
	return &TrafficData{data: data}
}

func (t *TrafficData) clone() map[string]*net.TCPListener {
	ce := make(map[string]*net.TCPListener)
	for k, v := range t.data {
		ce[k] = v
	}
	return ce
}
func (t *TrafficData) All() map[string]*net.TCPListener {
	return t.data
}

func (t *TrafficData) replace(addrs []string) (*TrafficData, error) {

	old := t.clone()
	datas := make(map[string]*net.TCPListener)

	for _, ad := range addrs {
		log.Debug("check traffic:", ad)
		v, has := old[ad]
		if has {
			delete(old, ad)
		} else {
			log.Debug("create traffic:", ad)

			l, err := net.Listen("tcp", ad)
			if err != nil {
				log.Error("listen tcp:", err)
				return nil, err
			}
			v = l.(*net.TCPListener)
		}
		datas[ad] = v
	}
	for n, o := range old {
		log.Debug("close old :", n)
		o.Close()
		log.Debug("close old done:", n)
	}

	return NewTrafficData(datas), nil
}
func (t *TrafficData) Shutdown() {
	t.stop = true
	list := t.All()

	for _, it := range list {
		it.Close()
	}
}
func (t *TrafficData) Close() {

	for _, it := range t.data {
		it.Close()
	}
	t.data = map[string]*net.TCPListener{}

}

//
//type Matcher struct {
//	mixListener *mixl.MixListener
//	cmuxMatch.CMuxMatch
//}
//
//func NewMatcher(port int, ls ...*net.TCPListener) *Matcher {
//	mixListener := mixl.NewMixListener(port, ls...)
//	return &Matcher{mixListener: mixListener, CMuxMatch: cmuxMatch.NewMatch(mixListener)}
//}
//
////	func (m *Matcher) Listeners() []*net.TCPListener {
////		return m.mixListener.Listeners()
////	}
//func (m *Matcher) Close() error {
//	m.mixListener.Close()
//	m.CMuxMatch.Close()
//	return nil
//}
