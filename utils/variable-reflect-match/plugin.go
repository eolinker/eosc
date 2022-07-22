package main

//PluginConfig 普通插件配置，在router、service、upstream的插件格式
type PluginConfig struct {
	Disable bool        `json:"disable"`
	Config  interface{} `json:"config"`
}

type PluginInterfaceConfig struct {
	buf []byte
	v   interface{}
}

func (p *PluginInterfaceConfig) UnmarshalJSON(bytes []byte) error {
	p.buf = bytes
	return nil
}

//
//func (p *PluginConfig) ReBuild(target reflect.Type) error {
//
//	if p.Config == nil {
//		p.Config = new(PluginInterfaceConfig)
//	}
//	v := reflect.New(target).Interface()
//	err := json.Unmarshal(p.Config.buf, &v)
//	if err != nil {
//		return err
//	}
//	p.Config.v = v
//	return nil
//}
