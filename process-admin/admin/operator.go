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
	setWorker(c *eosc.WorkerConfig) (*WorkerInfo, error)
	delWorker(id string) (*WorkerInfo, error)
	setSetting(name string, data []byte) error
	setVariable(namespace string, values map[string]string) error
	setProfession(name string, profession *eosc.ProfessionConfig) error
	delProfession(name string) error
	setHashValue(key string, values stringHash)
	delHashValue(key string)
}

func (d *imlAdminData) setWorker(cf *eosc.WorkerConfig) (*WorkerInfo, error) {

	log.Debug("set:", cf.Id, ",", cf.Profession, ",", cf.Name, ",", cf.Driver)
	p, has := d.professionData.Get(cf.Profession)
	if !has {
		return nil, fmt.Errorf("%s:%w", cf.Profession, eosc.ErrorProfessionNotExist)
	}
	if p.Mod == eosc.ProfessionConfig_Singleton {
		cf.Driver = cf.Name
	}
	driver, has := p.GetDriver(cf.Driver)
	if !has {
		return nil, fmt.Errorf("%s,%w", cf.Driver, eosc.ErrorDriverNotExist)
	}
	conf, usedVariables, err := d.variable.Unmarshal(cf.Body, driver.ConfigType())
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
	wInfo, hasInfo := d.workers.Get(cf.Id)
	if hasInfo && wInfo.worker != nil {

		if wInfo.config.Profession != cf.Profession {
			return nil, fmt.Errorf("%s:%w", cf.Profession, eosc.ErrorNotAllowCreateForSingleton)
		}
		e := wInfo.worker.Reset(conf, requires)
		if e != nil {
			return nil, e
		}
		d.requireManager.Set(cf.Id, utils.ArrayType(utils.MapKey(requires), func(t config.RequireId) string { return string(t) }))
		wInfo.reset(cf, wInfo.worker, driver.ConfigType())

		d.variable.SetRequire(cf.Id, usedVariables)
		return wInfo, nil
	}
	// create
	worker, err := driver.Create(cf.Id, cf.Name, conf, requires)
	if err != nil {
		log.Warn("worker-data set worker create:", err)
		return nil, err
	}

	if !hasInfo {
		if cf.Update == "" {
			cf.Update = eosc.Now()
		}
		if cf.Create == "" {
			cf.Create = eosc.Now()
		}
		wInfo = NewWorkerInfo(worker, cf, driver.ConfigType())
	} else {

		wInfo.reset(cf, worker, driver.ConfigType())
	}

	// store
	d.workers.Set(cf.Id, wInfo)
	d.requireManager.Set(cf.Id, utils.ArrayType(utils.MapKey(requires), func(t config.RequireId) string { return string(t) }))

	d.variable.SetRequire(cf.Id, usedVariables)
	log.Debug("worker-data set worker done:", cf.Id)

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
	return d.settings.SettingWorker(name, data)
}

func (d *imlAdminData) setVariable(namespace string, values map[string]string) (resultErr error) {
	log.Debug("check variable...")

	affectIds, oldValue, err := d.variable.SetByNamespace(namespace, values)
	if err != nil {
		return err
	}
	defer func() {
		if resultErr != nil {
			_, _, _ = d.variable.SetByNamespace(namespace, oldValue)
		}
	}()
	log.Debug("parse variable...")
	parse := variable.NewParse(d.variable)

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
			_, err = d.setWorker(info.config)
			if err != nil {
				return err
			}
		} else {
			err := d.settings.CheckVariable(name)
			if err != nil {
				return fmt.Errorf("setting %s unmarshal error:%s", name, err)
			}
			err = d.settings.Update(name)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *imlAdminData) delHashValue(key string) {
	d.customerHash.Del(key)
}
func (d *imlAdminData) setHashValue(key string, values stringHash) {

	d.customerHash.Set(key, values)
}

func (d *imlAdminData) setProfession(name string, profession *eosc.ProfessionConfig) error {
	return d.professionData.Set(name, profession)
}

func (d *imlAdminData) delProfession(name string) error {
	return d.professionData.Delete(name)
}
