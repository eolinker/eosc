package store

import (
	"fmt"
	"github.com/eolinker/eosc"
)

type Factory struct {

}

func (f *Factory) Create(params map[string]string) (eosc.IStore, error) {
	if params == nil{
		return nil,fmt.Errorf("create yaml store error:%w",eosc.ErrorParamsIsNil)
	}

	file,has:= params["file"]
	if !has{
		return nil,fmt.Errorf("create yaml store  file:%w ",eosc.ErrorParamNotExist)
	}

	return NewStore(file)

}
