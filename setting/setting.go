package setting

import (
	"github.com/eolinker/eosc"
)

var (
	settings = newSettings()
)

func RegisterSetting(name string, driver eosc.ISetting) error {
	return nil
}

type ISettings interface {
	//RegisterSetting(name string, driver eosc.ISetting) error
	GetDriver(name string) (eosc.ISetting, bool)
}
type tSettings struct {
}

func (s *tSettings) GetDriver(name string) (eosc.ISetting, bool) {
	//TODO implement me
	panic("implement me")
}

func newSettings() *tSettings {
	return &tSettings{}
}

func GetSettings() ISettings {
	return settings
}
