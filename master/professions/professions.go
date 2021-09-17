package professions

import (
	"encoding/json"

	"github.com/eolinker/eosc"
)

const (
	SpaceProfession = "profession"
)

type Professions struct {
	fileName        string
	professionSlice []eosc.ProfessionConfig
	professions     eosc.IUntyped
	drivers         eosc.IUntyped
}

func (p *Professions) ResetHandler(data []byte) error {
	professions, err := readProfessionConfig(p.fileName)
	if err != nil {
		return err
	}
	p.professionSlice = professions

	return nil
}

func (p *Professions) CommitHandler(cmd string, data []byte) error {
	return nil
}

func (p *Professions) Snapshot() []byte {
	data, _ := json.Marshal(p.professionSlice)
	return data
}

func (p *Professions) ProcessHandler(cmd string, body []byte) ([]byte, error) {
	return nil, nil
}

func NewProfessions(fileName string) *Professions {
	return &Professions{
		fileName:        fileName,
		professionSlice: nil,
	}
}
