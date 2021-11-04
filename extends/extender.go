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
type RegisterFunc func(eosc.IExtenderDriverRegister)

var (
	ErrorInvalidExtenderId     = errors.New("invalid  extender id")
	ErrorExtenderNameDuplicate = errors.New("duplicate extender factory name")
	ErrorExtenderNotFindLocal  = errors.New("not find extender on local")
	ErrorExtenderNotFindMark   = errors.New("not find extender on market")
	ErrorExtenderNoLatest      = errors.New("no latest for extender")
)

func ReadExtenderProject(group, project, version string) (*ExtenderRegister, error) {
	funcs, err := look(group, project, version)
	if err != nil {
		return nil, err
	}

	register := NewExtenderRegister(group, project)
	for _, f := range funcs {
		f(register)
	}
	return register, nil
}

func look(group, project, version string) ([]RegisterFunc, error) {
	inners, has := lookInner(group, project)
	if has {
		return inners, nil
	}
	dir := LocalExtenderPath(group, project, version)
	_, err := os.Stat(dir)
	if err != nil {
		log.Error(err)
		log.Debug(dir)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s-%s:%w:%s", group, project, ErrorExtenderNotFindLocal, dir)
		}
		return nil, err
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/*.so", dir))
	if err != nil {
		log.Error(err)
		log.Debug(dir)
		return nil, fmt.Errorf("%s-%s:%w", group, project, ErrorExtenderNotFindLocal)
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
