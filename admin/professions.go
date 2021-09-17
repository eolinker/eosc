package admin

type IProfessions interface {
	GetProfession(name string) (IProfession, bool)
}

type IProfession interface {
	HasDriver(name string) bool
	AppendAttr() []string
}
