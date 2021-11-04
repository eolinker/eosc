package eosc

type IProfessions interface {
	Set(name string, profession *ProfessionConfig) error
	Delete(name string) error
	GetProfession(name string) (IProfession, bool)
	Infos() []*ProfessionInfo
	Reset([]*ProfessionConfig)
	All() []*ProfessionConfig
}

type IProfession interface {
	Drivers() []*DriverInfo
	GetDriver(name string) (*DriverDetail, bool)
	HasDriver(name string) bool
	AppendAttr() []string
	DriversItem() []*Item
	Mod() ProfessionConfig_ProfessionMod
	Detail() *ProfessionDetail
}
