package admin

import "github.com/eolinker/eosc"

type IProfessions interface {
	GetProfession(name string) (IProfession, bool)
	List() []IProfession
	Infos() []*ProfessionInfo
}

type IProfession interface {
	Drivers() []*eosc.DriverInfo
	GetDriver(name string) (*eosc.DriverDetail, bool)
	HasDriver(name string) bool
	AppendAttr() []string
	Render(driver string) (*Render, bool)
	Renders() map[string]*Render
	DriversItem() []Item
	Info() *ProfessionInfo
}

type IProfessionEdit interface {
	SetDriver(name string, detail *eosc.DriverDetail) error
	DeleteDriver(name string) error
}

type IProfessionsEdit interface {
	Set(name string, profession *ProfessionInfo) error
	Delete(name string) error
}
