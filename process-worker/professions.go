package process_worker

import (
	"io"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

var _ IProfessions = (*Professions)(nil)

type IProfessions interface {
	Get(name string) (*Profession, bool)
	Sort() []*Profession
	List() []*Profession
}

type Professions struct {
	data eosc.IUntyped
}

func (ps *Professions) List() []*Profession {
	list := ps.data.List()
	rs := make([]*Profession, len(list))
	for i, v := range list {
		rs[i] = v.(*Profession)
	}
	return rs
}

func (ps *Professions) Sort() []*Profession {
	list := ps.List()

	sl := make([]*Profession, 0, len(list))
	sm := make(map[string]int)
	index := 0

	for len(list) > 0 {
		sc := 0
		for i, v := range list {
			dependenciesHas := 0
			for _, dep := range v.Dependencies {
				if _, has := sm[dep]; !has {
					break
				}

				dependenciesHas++
			}
			if dependenciesHas == len(v.Dependencies) {
				sl = append(sl, v)
				sm[v.Name] = index
				index++
				sc++
				list[i] = nil
			}
		}
		if sc == 0 {
			// todo profession dependencies error
			break
		}
		tmp := list[:0]
		for _, v := range list {
			if v != nil {
				tmp = append(tmp, v)
			}
		}
		list = tmp
	}
	return sl
}

func NewProfessions() *Professions {

	return &Professions{}
}

func (ps *Professions) init(configs []*eosc.ProfessionConfig, extends eosc.IExtenderDrivers) {
	data := eosc.NewUntyped()
	for _, c := range configs {
		log.Debug("add profession config:", c)
		p := NewProfession(c, extends)
		data.Set(c.Name, p)
	}
	ps.data = data
}
func (ps *Professions) Get(name string) (*Profession, bool) {
	p, b := ps.data.Get(name)
	log.Debug("get profession:", name, ":", b, "->", p)
	if !b {
		log.Debug("professions data:", ps.data)
	}
	if b {

		return p.(*Profession), true
	}
	return nil, false
}
func ReadProfessionData(r io.Reader) ([]*eosc.ProfessionConfig, error) {
	frame, err := utils.ReadFrame(r)
	if err != nil {
		return nil, err
	}

	pd := new(eosc.ProfessionConfigData)
	if e := proto.Unmarshal(frame, pd); e != nil {
		return nil, e
	}
	return pd.Data, nil
}
