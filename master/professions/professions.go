package professions


type Professions struct {

}

func NewProfessions() *Professions {
	return &Professions{}
}

func (p *Professions) ResetHandler(data []byte) error {
	panic("implement me")
}

func (p *Professions) CommitHandler(data []byte) error {
	panic("implement me")
}

func (p *Professions) Snapshot() []byte {
	panic("implement me")
}

func (p *Professions) ProcessHandler(propose []byte) (string, []byte, error) {
	panic("implement me")
}

