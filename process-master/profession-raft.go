package process_master

import (
	"os"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
	"github.com/golang/protobuf/proto"
)

type ProfessionRaft struct {
	eosc.IProfessionsData
}

func NewProfessionRaft(IProfessionsData eosc.IProfessionsData) *ProfessionRaft {
	return &ProfessionRaft{IProfessionsData: IProfessionsData}
}

func (p *ProfessionRaft) Set(name string, profession *eosc.ProfessionConfig) error {
	// todo raft sender
	return p.IProfessionsData.Set(name, profession)
}

func (p *ProfessionRaft) Delete(name string) error {
	// todo raft sender

	return p.IProfessionsData.Delete(name)
}

func (p *ProfessionRaft) encode() ([]byte, error) {
	list := p.IProfessionsData.All()
	pcd := &eosc.ProfessionConfigData{
		Data: list,
	}
	data, err := proto.Marshal(pcd)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (p *ProfessionRaft) decode(data []byte) ([]*eosc.ProfessionConfig, error) {
	pcd := new(eosc.ProfessionConfigData)
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
	p.IProfessionsData.Reset(ps)
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
