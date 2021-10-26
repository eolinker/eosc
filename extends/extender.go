package extends

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"
)

//RegisterFunc 注册函数
type RegisterFunc func(eosc.IExtenderRegister)

var (
	ErrorInvalidExtenderId     = errors.New("invalid  extender id")
	ErrorExtenderNameDuplicate = errors.New("duplicate extender factory name")
	ErrorExtenderNotFind       = errors.New("not find extender on local")
)

func ReadExtenderProject(group, name string) (*ExtenderRegister, error) {
	funcs, err := look(group, name)
	if err != nil {
		return nil, err
	}

	register := NewExtenderRegister(group, name)
	for _, f := range funcs {
		f(register)
	}
	return register, nil
}
func look(group, name string) ([]RegisterFunc, error) {
	inners, has := innerExtender[fmt.Sprint(group, "-", name)]
	if has {
		return inners, nil
	}
	dir := LocalExtenderPath(group, name)
	_, err := os.Stat(dir)
	if err != nil {
		log.Error(err)
		log.Debug(dir)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s-%s:%w:%s", group, name, ErrorExtenderNotFind, dir)
		}
		return nil, err
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/*.so", dir))
	if err != nil {
		log.Error(err)
		log.Debug(dir)
		return nil, fmt.Errorf("%s-%s:%w", group, name, ErrorExtenderNotFind)
	}
	registerFuncList := make([]RegisterFunc, 0, len(files))
	for _, f := range files {

		p, err := plugin.Open(f)
		if err != nil {
			log.Errorf("error to open plugin %s:%s", f, err.Error())
			continue
		}

		r, err := p.Lookup("Register")
		if err != nil {
			log.Errorf("call register from  plugin : %s : %s", f, err.Error())
			continue
		}
		f, ok := r.(RegisterFunc)
		if ok {
			registerFuncList = append(registerFuncList, f)
		}
	}
	return registerFuncList, nil
}
