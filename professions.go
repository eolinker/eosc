package eosc

type IProfessions interface {
	Set(name string, profession *ProfessionConfig) error
	Delete(name string) error
	Reset([]*ProfessionConfig)
	Names() []string
	GetProfession(name string) (IProfession, bool)
	All() []*ProfessionConfig
}

type IProfession interface {
	Drivers() []*DriverConfig
	GetDriver(name string) (*DriverConfig, bool)
	HasDriver(name string) bool
	AppendAttr() []string
	Mod() ProfessionConfig_ProfessionMod
}
