package eosc

type IProfessionsData interface {
	Set(name string, profession *ProfessionConfig) error
	Delete(name string) error
	GetProfession(name string) (IProfessionData, bool)

	Infos() []*ProfessionInfo
	Reset([]*ProfessionConfig)
	All() []*ProfessionConfig
}

type IProfessionData interface {
	Drivers() []*DriverInfo
	GetDriver(name string) (*DriverDetail, bool)
	HasDriver(name string) bool
	AppendAttr() []string
	DriversItem() []*Item
	Detail() *ProfessionDetail
}

type IProfessionDataEdit interface {
	SetDriver(name string, detail *DriverConfig) error
	DeleteDriver(name string) error
}
