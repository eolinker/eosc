package admin

import (
	"fmt"

	"github.com/eolinker/eosc"
)

func (a *Admin) ListEmployees(profession string) ([]interface{}, error) {
	list, err := a.workers.GetList(profession)
	if err != nil {
		return nil, err
	}
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs, nil
}

func (a *Admin) Update(profession, name, driver string, data []byte) error {
	return a.workers.Set(profession, name, driver, data)
}

func (a *Admin) Delete(profession, name string) (*eosc.WorkerInfo, error) {
	id := eosc.ToWorkerId(name, profession)
	return a.workers.Delete(id)
}

func (a *Admin) GetEmployee(profession, name string) (interface{}, error) {
	id := eosc.ToWorkerId(name, profession)
	return a.workers.GetWork(id)
}

func (a *Admin) Drivers(profession string) ([]*eosc.DriverInfo, error) {
	ip, has := a.professions.GetProfession(profession)
	if !has {
		return nil, fmt.Errorf("%s %w", profession, ErrorNotExist)
	}
	return ip.Drivers(), nil
}

func (a *Admin) DriverInfo(profession, driver string) (*eosc.DriverDetail, error) {
	ip, has := a.professions.GetProfession(profession)
	if !has {
		return nil, fmt.Errorf("profession %s:%w", profession, ErrorNotExist)
	}
	d, b := ip.GetDriver(driver)
	if !b {
		return nil, fmt.Errorf("driver %s of %s:%w", driver, profession, ErrorNotExist)
	}
	return d, nil
}

func (a *Admin) DriversItem(profession string) ([]*eosc.Item, error) {
	ip, has := a.professions.GetProfession(profession)
	if !has {
		return nil, fmt.Errorf("profession %s:%w", profession, ErrorNotExist)
	}
	return ip.DriversItem(), nil
}

func (a *Admin) ListProfessions() []*eosc.ProfessionInfo {
	return a.professions.Infos()
}
