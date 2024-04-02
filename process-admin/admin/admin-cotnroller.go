package admin

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"github.com/eolinker/eosc/utils/config"
	"github.com/eolinker/eosc/variable"
)

type iAdminOperator interface {
	setWorker(id, profession, name, driverName, version, desc string, body []byte, updateAt, createAt string) (*WorkerInfo, error)
	delWorker(id string) (*WorkerInfo, error)
	setSetting(name string, data []byte) error
	setVariable(namespace string, values map[string]string) error
	setProfession(name string, profession *eosc.ProfessionConfig) error
	delProfession(name string) error
}

func (d *imlAdminData) setWorker(id, profession, name, driverName, version, desc string, body []byte, updateAt, createAt string) (*WorkerInfo, error) {

	log.Debug("set:", id, ",", profession, ",", name, ",", driverName)
	p, has := d.professionData.Get(profession)
	if !has {
		return nil, fmt.Errorf("%s:%w", profession, eosc.ErrorProfessionNotExist)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		driverName = name
	}
	driver, has := p.GetDriver(driverName)
	if !has {
		return nil, fmt.Errorf("%s,%w", driverName, eosc.ErrorDriverNotExist)
	}

	conf, usedVariables, err := d.variable.Unmarshal(body, driver.ConfigType())
	if err != nil {
		return nil, err
	}

	requires, err := config.CheckConfig(conf, d)
	if err != nil {
		return nil, err
	}
	if dc, ok := driver.(eosc.IExtenderConfigChecker); ok {
		if err = dc.Check(conf, requires); err != nil {
			return nil, err
		}
	}
	wInfo, hasInfo := d.workers.Get(id)
	if hasInfo && wInfo.worker != nil {

		if wInfo.config.Profession != profession {
			return nil, fmt.Errorf("%s:%w", version, eosc.ErrorNotAllowCreateForSingleton)
		}
		e := wInfo.worker.Reset(conf, requires)
		if e != nil {
			return nil, e
		}
		d.requireManager.Set(id, utils.ArrayType(utils.MapKey(requires), func(t config.RequireId) string { return string(t) }))
		wInfo.reset(driverName, version, desc, body, wInfo.worker, driver.ConfigType(), updateAt, createAt)

		d.variable.SetRequire(id, usedVariables)
		return wInfo, nil
	}
	// create
	worker, err := driver.Create(id, name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return nil, err
	}

	if !hasInfo {
		if updateAt == "" {
			updateAt = eosc.Now()
		}
		if createAt == "" {
			createAt = eosc.Now()
		}
		wInfo = NewWorkerInfo(worker, id, profession, name, driverName, version, desc, createAt, updateAt, body, driver.ConfigType())
	} else {
		if updateAt == "" {
			updateAt = eosc.Now()
		}

		createAt = wInfo.config.Create

		wInfo.reset(driverName, version, desc, body, worker, driver.ConfigType(), updateAt, createAt)
	}

	// store
	d.workers.Set(id, wInfo)
	d.requireManager.Set(id, utils.ArrayType(utils.MapKey(requires), func(t config.RequireId) string { return string(t) }))

	d.variable.SetRequire(id, usedVariables)
	log.Debug("worker-data set worker done:", id)

	return wInfo, nil
}

func (d *imlAdminData) delWorker(id string) (*WorkerInfo, error) {
	worker, has := d.workers.Get(id)
	if !has {
		return nil, eosc.ErrorWorkerNotExits
	}

	if d.requireManager.RequireByCount(id) > 0 {
		return nil, eosc.ErrorRequire
	}

	if destroy, ok := worker.worker.(eosc.IWorkerDestroy); ok {
		err := destroy.Destroy()
		if err != nil {
			return nil, err
		}
	}
	d.workers.Del(id)
	d.requireManager.Del(id)
	d.variable.RemoveRequire(id)
	return worker, nil
}

func (d *imlAdminData) setSetting(name string, data []byte) error {
	return d.settings.SettingWorker(name, data, d.variable)
}

func (d *imlAdminData) setVariable(namespace string, values map[string]string) error {
	log.Debug("check variable...")
	affectIds, clone, errCheck := d.variable.Check(namespace, values)
	if errCheck != nil {
		return errCheck
	}

	log.Debug("parse variable...")
	parse := variable.NewParse(clone)

	for _, id := range affectIds {
		profession, name, success := eosc.SplitWorkerId(id)
		if !success {
			continue
		}
		if profession != Setting {
			info, has := d.workers.Get(id)
			if !has {
				return fmt.Errorf("worker(%s) not found, error is %s", id, ErrorNotExist)
			}
			_, _, err := parse.Unmarshal(info.Body(), info.ConfigType())
			if err != nil {
				return fmt.Errorf("unmarshal error:%s,body is '%s'", err, string(info.Body()))
			}
			_, err = d.setWorker(id, info.config.Profession, info.config.Name, info.config.Driver, info.config.Version, info.config.Description, info.Body(), info.config.Create, info.config.Update)
			if err != nil {
				return err
			}
		} else {
			err := d.settings.CheckVariable(name, clone)
			if err != nil {
				return fmt.Errorf("setting %s unmarshal error:%s", name, err)
			}
			err = d.settings.Update(name, d.variable)
			if err != nil {
				return err
			}
		}

	}

	return d.variable.SetByNamespace(namespace, values)

}
