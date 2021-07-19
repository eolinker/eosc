package store_memory_yaml

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	yaml_store "github.com/eolinker/eosc/modules/store-yaml"
)

func Register()  {
	eosc.RegisterStoreDriver("memory-yaml",new(Factory))
}

type Factory struct {

}

func (f *Factory) Create(params map[string]string) (eosc.IStore, error) {
	if params != nil{
		file,has:= params["file"]
		if has && file != ""{
			yamlStore, err := yaml_store.NewStore(file)
			if err != nil{
				log.Warnf("crate memory store from yaml fail:%s",err.Error())
			}else{
				return NewStore(yamlStore)
			}

 		}

	}

	return NewStore(nil)
}
