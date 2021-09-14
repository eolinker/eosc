package admin

type DriverInfo struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Label      string            `json:"label"`
	Desc       string            `json:"desc"`
	Profession string            `json:"profession"`
	Params     map[string]string `json:"params"`
}

type ProfessionInfo struct {
	Name         string       `json:"name"`
	LocalName    string       `json:"local_name"`
	Desc         string       `json:"desc"`
	Dependencies []string     `json:"dependencies"`
	AppendLabels []string     `json:"labels"`
	Drivers      []DriverInfo `json:"drivers"`
}
type ProfessionItem struct {
}

type WorkerInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Create string `json:"create_time"`
	Update string `json:"update_time"`
}

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type Render struct {

}

type IAdmin interface {
	ListProfessions() []ProfessionInfo
	ListEmployees(profession string) ([]interface{}, error)
	ListEmployeeNames(profession string) ([]string, error)
	Update(profession, name, driver string, data IData) (*WorkerInfo, error)
	Delete(profession, name string) (*WorkerInfo, error)
	GetEmployee(profession, name string) (interface{}, error)
	Render(profession, driver string) (*Render, error)
	Renders(profession string) (map[string]*Render, error)
	Drivers(profession string) ([]DriverInfo, error)
	DriverInfo(profession, driver string) (DriverDetail, error)
	DriversItem(profession string) ([]Item, error)
	SearchBySkill(profession string, skill []string) ([]WorkerInfo, error)
	//ExportByProfession(profession string) ([]StoreValue, error)
}

