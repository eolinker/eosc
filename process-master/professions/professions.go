package professions

import (
	"os"

	"github.com/eolinker/eosc/utils"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	data untypeProfessionData
}

func (p *Professions) Encode(startIndex int) ([]byte, []*os.File, error) {

	data, err := p.encode()
	if err != nil {
		return nil, nil, err
	}
	return utils.EncodeFrame(data), nil, nil
}

func (p *Professions) encode() ([]byte, error) {
	list := p.data.Data()
	pcd := &eosc.ProfessionConfigData{
		Data: list,
	}
	data, err := proto.Marshal(pcd)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (p *Professions) decode(data []byte) ([]*eosc.ProfessionConfig, error) {
	pcd := new(eosc.ProfessionConfigData)
	err := proto.Unmarshal(data, pcd)
	if err != nil {
		return nil, err
	}
	return pcd.Data, nil

}

func (p *Professions) Set(name string, profession *eosc.ProfessionConfig) error {
	adminProfession := NewProfession(profession)

	p.data.Set(name, adminProfession)
	return nil
}

func (p *Professions) Delete(name string) error {
	p.data.Del(name)
	return nil
}

func (p *Professions) List() []eosc.IProfessionData {
	professions := p.data.List()
	ps := make([]eosc.IProfessionData, 0, len(professions))
	for _, pv := range professions {
		ps = append(ps, pv)
	}
	return ps
}

func (p *Professions) Infos() []*eosc.ProfessionInfo {
	professions := p.data.List()
	ps := make([]*eosc.ProfessionInfo, 0, len(professions))
	for _, pv := range professions {
		ps = append(ps, pv.info)
	}
	return ps
}

func (p *Professions) GetProfession(name string) (eosc.IProfessionData, bool) {
	vl, has := p.data.Get(name)

	return vl, has
}

func (p *Professions) Reset(professions []*eosc.ProfessionConfig) {
	pfs := NewProfessionData()
	for _, pf := range professions {
		adminProfession := NewProfession(pf)
		pfs.Set(pf.Name, adminProfession)
	}
	p.data = pfs
}

func (p *Professions) ResetHandler(data []byte) error {

	ps, err := p.decode(data)
	if err != nil {
		return err
	}
	p.Reset(ps)
	return nil
}

func (p *Professions) CommitHandler(cmd string, data []byte) error {
	return nil
}

func (p *Professions) Snapshot() []byte {

	data, err := p.encode()
	if err != nil {
		return nil
	}

	return data
}

func (p *Professions) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	return nil, nil
}

func NewProfessions() *Professions {
	return &Professions{
		data: NewProfessionData(),
	}
}
