package eosc

type IProfessionsData interface {
	GetProfession(name string) (IProfessionData, bool)
	List() []IProfessionData
	Infos() []*ProfessionInfo
	IDataMarshaler
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

type IProfessionsDataEdit interface {
	Set(name string, profession *ProfessionConfig) error
	Delete(name string) error
}
