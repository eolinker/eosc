package admin

import "github.com/eolinker/eosc"

func (a *Admin) ListEmployees(profession string) ([]interface{}, error) {
	panic("implement me")
}

func (a *Admin) ListEmployeeNames(profession string) ([]string, error) {
	panic("implement me")
}

func (a *Admin) Update(profession, name, driver string, data eosc.IData) (*eosc.WorkerInfo, error) {
	panic("implement me")
}

func (a *Admin) Delete(profession, name string) (*eosc.WorkerInfo, error) {
	panic("implement me")
}

func (a *Admin) GetEmployee(profession, name string) (interface{}, error) {
	panic("implement me")
}

func (a *Admin) SearchBySkill(profession string, skill []string) ([]eosc.WorkerInfo, error) {
	panic("implement me")
}

func (a *Admin) Render(profession, driver string) (*eosc.Render, error) {
	panic("implement me")
}

func (a *Admin) Renders(profession string) (map[string]*eosc.Render, error) {
	panic("implement me")
}

func (a *Admin) Drivers(profession string) ([]eosc.DriverInfo, error) {
	panic("implement me")
}

func (a *Admin) DriverInfo(profession, driver string) (eosc.DriverDetail, error) {
	panic("implement me")
}

func (a *Admin) DriversItem(profession string) ([]eosc.Item, error) {
	panic("implement me")
}

func (a *Admin) ListProfessions() []eosc.ProfessionInfo {
	panic("implement me")
}
