package setting

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"reflect"
	"strings"
	"sync"
)

var (
	settings ISettings = newSettings()
)

func RegisterSetting(name string, driver eosc.ISetting) error {
	settings.registerSetting(name, driver)
	return nil
}

type ISettings interface {
	registerSetting(name string, driver eosc.ISetting) error
	eosc.ISettings
}
type tSettings struct {
	lock      sync.RWMutex
	data      map[string]eosc.ISetting
	configs   map[string]interface{}
	orgConfig map[string][]byte
}

//func (s *tSettings) Set(name string, body []byte, variable eosc.IVariable) (format interface{}, update []*eosc.WorkerConfig, delete []string, err error) {
//	driver, has := s.GetDriver(name)
//	if !has {
//		err = eosc.ErrorDriverNotExist
//		return
//	}
//	configType := driver.ConfigType()
//	switch driver.Mode() {
//	case eosc.SettingModeReadonly:
//		err = eosc.ErrorUnsupportedKind
//		return
//	case eosc.SettingModeSingleton:
//		{
//			orgConfig := make(map[string]interface{})
//			err = json.Unmarshal(body, &orgConfig)
//			if err != nil {
//				return
//			}
//
//			err = s.SettingWorker(name, body, variable)
//			if err != nil {
//				return
//			}
//
//			config := formatConfig(orgConfig, configType)
//			s.lock.Lock()
//			s.configs[name] = config
//			s.orgConfig[name] = body
//			s.lock.Unlock()
//			format = config
//			outputBody, _ := json.Marshal(config)
//			update = []*eosc.WorkerConfig{
//				{
//					Id:          fmt.Sprintf("%s@setting", name),
//					Profession:  "setting",
//					Name:        name,
//					Driver:      name,
//					Create:      eosc.Now(),
//					Update:      eosc.Now(),
//					Body:        outputBody,
//					Description: "",
//				},
//			}
//			delete = nil
//			return
//		}
//	case eosc.SettingModeBatch:
//		orgs := splitConfig(body)
//		cfgs := make([]interface{}, 0, len(orgs))
//		orgObjs := make([]map[string]interface{}, 0, len(orgs))
//		usagesAll := make([][]string, 0, len(orgs))
//		for _, org := range orgs {
//			cfg, usages, errI := variable.Unmarshal(org, configType)
//			if errI != nil {
//				err = errI
//				return
//			}
//			orgConfig := make(map[string]interface{})
//			err = json.Unmarshal(body, &orgConfig)
//			if err != nil {
//				return
//			}
//			orgObjs = append(orgObjs, orgConfig)
//			cfgs = append(cfgs, cfg)
//			usagesAll = append(usagesAll, usages)
//		}
//		var updateIds []string
//		updateIds, delete, err = driver.Set(cfgs...)
//		if err != nil {
//			return
//		}
//		for i := range cfgs {
//			variable.SetRequire(update[i].Id, usagesAll[i])
//		}
//		for _, id := range delete {
//			variable.RemoveRequire(id)
//		}
//
//		for i := range orgObjs {
//			orgObjs[i] = formatConfig(orgObjs[i], configType)
//		}
//		format = orgObjs
//	}
//	return
//}

func (s *tSettings) Update(name string, variable eosc.IVariable) error {
	log.Debug("setting update:", name)
	driver, has := s.GetDriver(name)
	if !has {
		return nil
	}
	if driver.Mode() != eosc.SettingModeSingleton {
		return nil
	}
	s.lock.RLock()
	org, has := s.orgConfig[name]
	s.lock.RUnlock()
	if !has {
		return nil
	}
	conf, useVariable, err := variable.Unmarshal(org, driver.ConfigType())
	if err != nil {
		return err
	}

	err = driver.Set(conf)
	if err != nil {
		return err
	}
	variable.SetRequire(fmt.Sprintf("%s@setting", name), useVariable)

	return nil

}

func (s *tSettings) CheckVariable(name string, variable eosc.IVariable) (err error) {
	driver, has := s.GetDriver(name)
	if !has {
		return nil
	}
	if driver.Mode() != eosc.SettingModeSingleton {
		return nil
	}
	s.lock.RLock()
	org, has := s.orgConfig[name]
	s.lock.RUnlock()
	if !has {
		return nil
	}
	_, _, err = variable.Unmarshal(org, driver.ConfigType())

	return err

}

func (s *tSettings) GetConfig(name string) interface{} {

	if driver, has := s.GetDriver(name); has {
		switch driver.Mode() {
		case eosc.SettingModeReadonly:
			{
				return driver.Get()
			}
		case eosc.SettingModeSingleton:
			v, yes := s.configs[name]
			if yes {
				return v
			}
			return driver.Get()
		case eosc.SettingModeBatch:
			return nil
		}
	}
	return nil
}

func (s *tSettings) SettingWorker(name string, org []byte, variable eosc.IVariable) (err error) {
	log.Debug("setting Set:", name, " org:", string(org))

	driver, has := s.GetDriver(name)
	if !has {
		return eosc.ErrorDriverNotExist
	}

	if driver.Mode() != eosc.SettingModeSingleton {
		return eosc.ErrorUnsupportedKind
	}
	configType := driver.ConfigType()
	cfg, vs, err := variable.Unmarshal(org, configType)
	if err != nil {
		return err
	}
	err = driver.Set(cfg)
	if err != nil {
		return err
	}
	variable.SetRequire(fmt.Sprintf("%s@setting", name), vs)
	config := FormatConfig(org, configType)
	s.lock.Lock()
	s.configs[name] = config
	s.orgConfig[name] = org
	s.lock.Unlock()
	return nil

}
func (s *tSettings) registerSetting(name string, driver eosc.ISetting) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	name = toDriverName(name)

	_, has := s.data[name]
	if has {
		return eosc.ErrorRegisterConflict
	}
	s.data[name] = driver
	return nil
}

func (s *tSettings) GetDriver(name string) (eosc.ISetting, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	name = toDriverName(name)
	st, has := s.data[name]
	return st, has
}

func newSettings() *tSettings {
	return &tSettings{
		data:      make(map[string]eosc.ISetting),
		configs:   make(map[string]interface{}),
		orgConfig: map[string][]byte{},
	}
}

func GetSettings() ISettings {
	return settings
}

func toDriverName(id string) string {
	if i := strings.Index(id, "@"); i > -1 {
		id = id[:i]
		if len(id) == 0 {
			id = id[i+1:]
		}
	}
	return id
}
func FormatConfig(input []byte, tp reflect.Type) interface{} {
	orgConfig := make(map[string]interface{})
	json.Unmarshal(input, &orgConfig)
	return formatConfig(orgConfig, tp)
}
func formatConfig(config map[string]interface{}, tp reflect.Type) map[string]interface{} {
	switch tp.Kind() {
	case reflect.Interface, reflect.Ptr:
		return formatConfig(config, tp.Elem())
	case reflect.Struct:
		nc := make(map[string]interface{}, tp.NumField())
		for i := 0; i < tp.NumField(); i++ {
			name, has := tp.Field(i).Tag.Lookup("json")
			if !has {
				name = tp.Field(i).Name
			}
			name = strings.Split(name, ",")[0]
			if fv, has := config[name]; has {
				nc[name] = fv
			}
		}
		return nc
	}
	return config
}
