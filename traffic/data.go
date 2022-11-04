package traffic

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	cmuxMatch "github.com/eolinker/eosc/traffic/cmux-match"
	"github.com/eolinker/eosc/traffic/mixl"
	"net"
)

type MatcherData struct {
	data eosc.Untyped[int, *Matcher]
}

func NewMatcherData(tfConf ...*PbTraffic) *MatcherData {

	m := &MatcherData{
		data: eosc.BuildUntyped[int, *Matcher](),
	}
	listeners, err := toListeners(tfConf)
	log.Debug("read listeners: ", len(listeners))
	if err != nil {
		log.Warn("read listeners:", err)
	}

	lm := make(map[int][]*net.TCPListener)
	for _, l := range listeners {
		p := readPort(l.Addr())
		lm[p] = append(lm[p], l)
	}
	for p, ls := range lm {
		mixlistener := mixl.NewMixListener(p, ls...)
		m.Set(p, &Matcher{
			mixListener: mixlistener,
			CMuxMatch:   cmuxMatch.NewMatch(mixlistener),
		})
	}
	return m
}

func (m *MatcherData) Set(port int, mux *Matcher) {
	m.data.Set(port, mux)
}
func (m *MatcherData) Get(port int) *Matcher {
	o, has := m.data.Get(port)
	if !has {
		return nil
	}
	return o
}
func (m *MatcherData) Del(port int) (*Matcher, bool) {
	return m.data.Del(port)
}
func (m *MatcherData) Clone() *MatcherData {
	return &MatcherData{
		data: m.data.Clone(),
	}
}
func (m *MatcherData) All() map[int]*Matcher {
	return m.data.All()
}

type Matcher struct {
	mixListener *mixl.MixListener
	cmuxMatch.CMuxMatch
}

func NewMatcher(port int, ls ...*net.TCPListener) *Matcher {
	mixListener := mixl.NewMixListener(port, ls...)
	return &Matcher{mixListener: mixListener, CMuxMatch: cmuxMatch.NewMatch(mixListener)}
}

func (m *Matcher) Listeners() []*net.TCPListener {
	return m.mixListener.Listeners()
}
func (m *Matcher) Close() error {
	m.mixListener.Close()
	m.CMuxMatch.Close()
	return nil
}
