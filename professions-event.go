package eosc

import "fmt"

var _ IStoreEventHandler = (*Professions)(nil)

func (ps *Professions) OnDel(v StoreValue) error {
	if p, has := ps.get(v.Profession); has {
		return p.del(v.Name)
	}
	return fmt.Errorf("%s:%w", v.Profession, ErrorProfessionNotExist)
}

func (ps *Professions) OnInit(vs []StoreValue) error {

	for i := range vs {
		if p, has := ps.get(vs[i].Profession); has {
			p.setId(vs[i].Name, vs[i].Id)
		} else {
			return fmt.Errorf("%s:%w", vs[i].Profession, ErrorProfessionNotExist)
		}
	}
	return nil

}

func (ps *Professions) OnChange(v StoreValue) error {
	if p, has := ps.get(v.Profession); has {
		return p.del(v.Name)
	}
	return fmt.Errorf("%s:%w", v.Profession, ErrorProfessionNotExist)
}
