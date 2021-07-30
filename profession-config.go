package eosc

import (
	"fmt"
)

type DriverConfig struct {
	ID    string
	Name  string
	Label string
	Desc  string

	Params map[string]string
}

func (dc *DriverConfig) toInfo(professionName string) DriverInfo {
	return DriverInfo{
		Id:         dc.ID,
		Name:       dc.Name,
		Label:      dc.Label,
		Desc:       dc.Desc,
		Profession: professionName,
		Params:     dc.Params,
	}
}
func toInfo(dcs []DriverConfig, professionName string) []DriverInfo {
	res := make([]DriverInfo, 0, len(dcs))
	for _, d := range dcs {
		res = append(res, d.toInfo(professionName))
	}
	return res
}

type ProfessionConfig struct {
	Name         string
	Label        string
	Desc         string
	Dependencies []string
	AppendLabel  []string
	Drivers      []DriverConfig
}

func (pc *ProfessionConfig) create(driverRegister IDriverRegister) (IProfession, error) {

	ds := make([]*_TProfessionDriver, 0, len(pc.Drivers))
	for _, dc := range pc.Drivers {
		if driverFactory, has := driverRegister.GetProfessionDriver(dc.ID); has {
			d, err := driverFactory.Create(pc.Name, dc.Name, dc.Label, dc.Desc, dc.Params)
			if err != nil {
				return nil, fmt.Errorf("create driver %s:%w", dc.ID, err)
			}
			ds = append(ds,
				newProfessionDriver(d, driverFactory.ExtendInfo(), DriverInfo{
					Name:       dc.Name,
					Label:      dc.Label,
					Desc:       dc.Desc,
					Profession: pc.Name,
					Params:     dc.Params,
				}))
		} else {
			return nil, fmt.Errorf("%s:%w", dc.ID, ErrorDriverNotExist)
		}
	}
	p := &Profession{
		name:         pc.Name,
		label:        pc.Label,
		desc:         pc.Desc,
		dependencies: pc.Dependencies,
		appendLabels: pc.AppendLabel,
		drivers:      NewDrivers(ds),
		data:         NewUntyped(),
	}
	return p, nil
}

type ProfessionConfigs []ProfessionConfig

func (pcs ProfessionConfigs) Gen(driverRegister IDriverRegister, store IStore) (*Professions, error) {
	infos := make([]ProfessionInfo, 0, len(pcs))

	for _, p := range pcs {

		infos = append(infos, ProfessionInfo{
			Name:         p.Name,
			LocalName:    p.Label,
			Desc:         p.Desc,
			Dependencies: p.Dependencies,
			AppendLabels: p.AppendLabel,
			Drivers:      toInfo(p.Drivers, p.Name),
		})
	}

	infos, err := checkProfessions(infos)
	if err != nil {
		return nil, err
	}
	ps := &Professions{
		infos:   infos,
		store:   store,
		data:    newTProfessionUntyped(),
		workers: NewWorkers(),
	}
	for _, p := range pcs {
		profession, err := p.create(driverRegister)
		if err != nil {
			return nil, err
		}
		ps.data.add(p.Name, profession)
	}

	store.GetListener().AddListen(ps)
	return ps, nil
}


func checkProfessions(infos []ProfessionInfo) ([]ProfessionInfo, error) {

	less := make([]ProfessionInfo, len(infos))
	copy(less, infos)
	plist := make([]ProfessionInfo, 0, len(less))
	exist := make(map[string]int)
	do := 1
	for do > 0 && len(less) > 0 {
		do = 0
		ls := less
		less = make([]ProfessionInfo, 0, len(ls))
	FIND:
		for _, v := range ls {

			for _, d := range v.Dependencies {
				if _, has := exist[d]; !has {
					less = append(less, v)
					continue FIND
				}
			}
			plist = append(plist, v)
			exist[v.Name] = 1
			do++

		}
	}
	if len(less) > 0 {
		return nil, ErrorProfessionDependencies
	}
	return plist, nil

}