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
	lock sync.RWMutex
	data map[string]eosc.ISetting
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

	return formatConfig(orgConfig, driver.ConfigType()), nil

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
	return &tSettings{}
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
