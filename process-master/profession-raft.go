package process_master

import (
	"os"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"
)

type ProfessionRaft struct {
	eosc.IProfessions
}

func (p *ProfessionRaft) Append(cmd string, data []byte) error {
	return nil
}

func (p *ProfessionRaft) Complete() error {
	return nil
}

func NewProfessionRaft(IProfessionsData eosc.IProfessions) *ProfessionRaft {
	return &ProfessionRaft{IProfessions: IProfessionsData}
}

func (p *ProfessionRaft) Set(name string, profession *eosc.ProfessionConfig) error {
	// todo raft sender
	return p.IProfessions.Set(name, profession)
}

func (p *ProfessionRaft) Delete(name string) error {
	// todo raft sender

	return p.IProfessions.Delete(name)
}

func (p *ProfessionRaft) encode() ([]byte, error) {
	list := p.IProfessions.All()
	pcd := &eosc.ProfessionConfigs{
		Data: list,
	}
	data, err := proto.Marshal(pcd)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (p *ProfessionRaft) decode(data []byte) ([]*eosc.ProfessionConfig, error) {
	pcd := new(eosc.ProfessionConfigs)
	err := proto.Unmarshal(data, pcd)
	if err != nil {
		return nil, err
	}
	return pcd.Data, nil

}
func (p *ProfessionRaft) Encode(startIndex int) ([]byte, []*os.File, error) {

	data, err := p.encode()
	if err != nil {
		return nil, nil, err
	}
	return utils.EncodeFrame(data), nil, nil
}
func (p *ProfessionRaft) ResetHandler(data []byte) error {

	ps, err := p.decode(data)
	if err != nil {
		return err
	}
	p.IProfessions.Reset(ps)
	return nil
}

func (p *ProfessionRaft) CommitHandler(cmd string, data []byte) error {
	return nil
}

func (p *ProfessionRaft) Snapshot() []byte {

	data, err := p.encode()
	if err != nil {
		return nil
	}

	return data
}

func (p *ProfessionRaft) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {
	// todo
	return nil, nil, nil
}
