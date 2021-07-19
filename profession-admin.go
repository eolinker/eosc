package eosc

var (
	_ IAdmin = (*Professions)(nil)
)

func (ps *Professions) ListProfessions() []ProfessionInfo {
	return ps.infos
}

func (ps *Professions) ListEmployees(profession string) ([]interface{}, error) {

	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	ids := p.ids()
	res := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		value, has := ps.store.Get(id)
		if has {
			p, has := ps.get(value.Profession)
			if has {
				res = append(res, p.genInfo(&value))
			}
		}
	}
	return res, nil
}

func (ps *Professions) Delete(profession, name string) error {
	if ps.store.ReadOnly() {
		return ErrorStoreReadOnly
	}
	p, has := ps.get(profession)
	if !has {
		return ErrorProfessionNotExist
	}
	id, hasId := p.getId(name)
	if !hasId {
		id = name
	}
	v, has := ps.store.Get(id)
	if !has {
		return ErrorWorkerNotExits
	}
	if v.Profession != profession {
		return ErrorWorkerNotExits
	}
	if hasId && v.Name != name {
		return ErrorWorkerNotExits
	}

	ps.store.Del(id)

	return nil
}

func (ps *Professions) Get(profession, name string) (interface{}, error) {
	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	id, has := p.getId(name)
	if !has {
		id = name
	}

	value, has := ps.store.Get(id)
	if !has {
		return nil, ErrorWorkerNotExits
	}
	item := make(map[string]interface{})
	e := value.IData.UnMarshal(&item)
	if e != nil {
		return nil, e
	}
	item["name"] = value.Name
	item["id"] = id
	item["driver"] = value.Driver
	item["create_time"] = value.CreateTime
	item["update_time"] = value.UpdateTime
	return item, nil
}

func (ps *Professions) Render(profession, driver string) (*Render, error) {
	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	d, has := p.getDriver(driver)
	if !has {
		return nil, ErrorDriverNotExist
	}

	return GenRender(d.ConfigType()), nil

}

func (ps *Professions) DriverInfo(profession, driver string) (DriverDetail, error) {

	p, has := ps.get(profession)
	if !has {
		return DriverDetail{}, ErrorProfessionNotExist
	}
	d, has := p.getDriver(driver)
	if !has {
		return DriverDetail{}, ErrorDriverNotExist
	}

	return DriverDetail{
		DriverInfo: d.DriverInfo(),
		Extends:    d.ExtendInfo(),
	}, nil
}
func (ps *Professions) DriversItem(profession string) ([]Item, error) {

	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}

	ds := p.getDrivers()
	list := make([]Item, 0, len(ds))
	for _, d := range ds {
		list = append(list, Item{
			Value: d.DriverInfo().Name,
			Label: d.DriverInfo().Label,
		})
	}
	return list, nil
}
func (ps *Professions) Drivers(profession string) ([]DriverInfo, error) {

	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}

	ds := p.getDrivers()
	list := make([]DriverInfo, 0, len(ds))
	for _, d := range ds {
		list = append(list, d.DriverInfo())
	}
	return list, nil
}

func (ps *Professions) SearchBySkill(profession string, skill []string) ([]WorkerInfo, error) {
	panic("implement me")
}

func (ps *Professions) Update(profession, name, driver string, data IData) (*WorkerInfo, error) {
	if ps.store.ReadOnly() {
		return nil, ErrorStoreReadOnly
	}
	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	id, hasId := p.getId(name)
	if !hasId {
		id = name
	}
	v, has := ps.store.Get(id)
	if !has {
		return nil, ErrorWorkerNotExits
	}

	return &WorkerInfo{
		Id:     v.Id,
		Name:   v.Name,
		Driver: v.Profession,
		Create: v.CreateTime,
		Update: v.UpdateTime,
	}, nil
}

func (ps *Professions) Renders(profession string) (map[string]*Render, error) {
	p, has := ps.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	ds := p.getDrivers()
	renders := make(map[string]*Render, len(ds))
	for _, d := range ds {
		renders[d.DriverInfo().Name] = GenRender(d.ConfigType())
	}
	return renders, nil
}
