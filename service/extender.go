package service

type IExtenderSetting interface {
	Reset(data []*ExtendsInfo) error
	Save(extenderName string, data *ExtendsInfo) error
	Del(extenderName string) (*ExtendsInfo, error)
	List(sortItem string, desc bool) ([]*ExtendsInfo, error)
	IExtenderPluginSetting
}

type IExtenderPluginSetting interface {
	GetPlugin(pluginID string) (*Plugin, error)
	Plugins(extenderName string) ([]*Plugin, error)
}
