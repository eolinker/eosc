package setting

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
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

func (s *tSettings) Update(name string, variable eosc.IVariable) (err error) {
	driver, has := s.GetDriver(name)
	if !has {
		return nil
	}
	if driver.ReadOnly() {
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
	variable.SetVariablesById(fmt.Sprintf("%s@setting", name), useVariable)

	return nil

}

func (s *tSettings) CheckVariable(name string, variable eosc.IVariable) (err error) {
	driver, has := s.GetDriver(name)
	if !has {
		return nil
	}
	if driver.ReadOnly() {
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
		if driver.ReadOnly() {
			return driver.Get()
		}
		v, yes := s.configs[name]
		if yes {
			return v
		}
		return driver.Get()
	}
	return nil
}

func (s *tSettings) Set(name string, org []byte, variable eosc.IVariable) (format interface{}, err error) {

	driver, has := s.GetDriver(name)
	if !has {
		return nil, eosc.ErrorDriverNotExist
	}

	if driver.ReadOnly() {
		return nil, eosc.ErrorStoreReadOnly
	}
	cfg, vs, err := variable.Unmarshal(org, driver.ConfigType())
	if err != nil {
		return nil, err
	}
	err = driver.Set(cfg)
	if err != nil {
		return nil, err
	}
	variable.SetVariablesById(fmt.Sprintf("%s@setting", name), vs)

	orgConfig := make(map[string]interface{})
	err = json.Unmarshal(org, &format)
	if err != nil {
		return nil, err
	}

	config := formatConfig(orgConfig, driver.ConfigType())
	s.lock.Lock()
	s.configs[name] = config
	s.orgConfig[name] = org
	s.lock.Unlock()
	return config, nil

}
func formatConfig(config map[string]interface{}, tp reflect.Type) interface{} {
	switch tp.Kind() {
	case reflect.Interface, reflect.Ptr:
		return formatConfig(config, tp.Elem())
	case reflect.Struct:
		nc := make(map[string]interface{}, tp.NumField())
		for i := 0; i < tp.NumField(); i++ {
			name := tp.Field(i).Name
			if fv, has := config[name]; has {
				nc[name] = fv
			}
		}
		return nc
	}
	return config
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
