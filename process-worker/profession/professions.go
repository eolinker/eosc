package profession

import "github.com/eolinker/eosc"

type IProfessions interface {
	Get(name string) (*Profession, bool)
}
type Professions struct {
	configs []*eosc.ProfessionConfig
}

func (p *Professions) Get(name string) (*Profession, bool) {

}
