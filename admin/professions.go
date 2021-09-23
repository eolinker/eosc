package admin

import "github.com/eolinker/eosc"

type IProfessions interface {
	GetProfession(name string) (IProfession, bool)
	List() []IProfession
	Infos() []*eosc.ProfessionInfo
}

type IProfession interface {
	Drivers() []*eosc.DriverInfo
	GetDriver(name string) (*eosc.DriverDetail, bool)
	HasDriver(name string) bool
	AppendAttr() []string
	DriversItem() []*eosc.Item
	Detail() *eosc.ProfessionDetail
}

type IProfessionEdit interface {
	SetDriver(name string, detail *eosc.DriverConfig) error
	DeleteDriver(name string) error
}

type IProfessionsEdit interface {
	Set(name string, profession *eosc.ProfessionConfig) error
	Delete(name string) error
}
