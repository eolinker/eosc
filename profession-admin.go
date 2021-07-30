package eosc

import (
	"fmt"
)

var (
	_ IAdmin = (*Professions)(nil)
)

//ListProfessions
func (ps *Professions) ListProfessions() []ProfessionInfo {
	return ps.infos
}

//ListEmployees
func (ps *Professions) ListEmployees(profession string) ([]interface{}, error) {

	p, has :=ps.data.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	ids := p.ids()
	res := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		value, has := ps.store.Get(id)
		if has {
			p, has :=ps.data.get(value.Profession)
			if has {
				res = append(res, p.genInfo(&value))
			}
		}
	}
	return res, nil
}

func (ps *Professions) Delete(profession, name string)(*WorkerInfo, error ){
	if ps.store.ReadOnly() {
		return nil, ErrorStoreReadOnly
	}
	p, has :=ps.data.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	id, hasId := p.getId(name)
	if !hasId {
		id = name
	}
	v, has := ps.store.Get(id)
	if !has {
		return nil,ErrorWorkerNotExits
	}
	if v.Profession != profession {
		return nil,ErrorWorkerNotExits
	}
	if hasId && v.Name != name {
		return nil,ErrorWorkerNotExits
	}

	err:=ps.store.Del(id)
	if err!= nil{
		return nil,err
	}
	return &WorkerInfo{
		Id:     v.Id,
		Name:   v.Name,
		Driver: v.Driver,
		Create: v.CreateTime,
		Update: v.UpdateTime,
	},nil

}

func (ps *Professions) GetEmployee(profession, name string) (interface{}, error) {
	p, has :=ps.data.get(profession)
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
	p, has :=ps.data.get(profession)
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

	p, has :=ps.data.get(profession)
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

	p, has :=ps.data.get(profession)
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

	p, has :=ps.data.get(profession)
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
	p,has:=ps.data.get(profession)
	if !has{
		return nil,fmt.Errorf("%s:%w",profession,ErrorProfessionNotExist)
	}

	ids:=p.ids()
	res:=make([]WorkerInfo,0,len(ids))
	for _,id:=range ids{

		w,has:=ps.workers.Get(id)
		if has{
			for _,s:=range skill{
				if w.CheckSkill(s){
					v,has:=ps.store.Get(id)
					if has{
						res = append(res, WorkerInfo{
							Id:     id,
							Name:   v.Name,
							Driver: v.Driver,
							Create: v.CreateTime,
							Update: v.UpdateTime,
						})
					}
				}
			}
		}
	}
	return res,nil
}

func (ps *Professions) Update(profession, name, driver string, data IData) (*WorkerInfo, error) {
	if ps.store.ReadOnly() {
		return nil, ErrorStoreReadOnly
	}
	p, has := ps.data.get(profession)
	if !has {
		return nil, ErrorProfessionNotExist
	}
	id, hasId := p.getId(name)
	if !hasId {
		id = name
	}
	v, has := ps.store.Get(id)
	if !has {
		if driver == ""{
			return nil, fmt.Errorf("driver:%w",ErrorRequire)
		}

		if _,dhas:=p.getDriver(driver);!dhas{
			return nil, fmt.Errorf("%s:%w",driver,ErrorDriverNotExist)
		}

		id =  fmt.Sprintf("%s@%s", name, profession)
		v = StoreValue{
			Id:        id,
			Profession: profession,
			Name:       name,
			Driver:     driver,
			CreateTime: Now(),
			UpdateTime: Now(),
			IData:      data,
			Sing:       "",
		}
	}else{
		if  driver == ""{
			driver = v.Driver
		}else {
			v.Driver = driver
		}
	}
	v.IData = data
	v.UpdateTime = Now()
	err:=p.CheckerConfig(driver,v.IData,ps.workers)
	if err!= nil{
		return nil,err
	}

	if err:=ps.store.Set(v);err!= nil{
		return nil,err
	}
	p.setId(name,id)

	return &WorkerInfo{
		Id:     v.Id,
		Name:   v.Name,
		Driver: v.Profession,
		Create: v.CreateTime,
		Update: v.UpdateTime,
	}, nil
}

func (ps *Professions) Renders(profession string) (map[string]*Render, error) {
	p, has :=ps.data.get(profession)
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
