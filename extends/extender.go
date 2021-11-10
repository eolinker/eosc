package extends

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/env"
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

type ExtenderVersion struct {
	*VersionInfo
	Arches []string `json:"arches"`
}

func ExtendersRequest(group, project string) ([]*ExtenderVersion, error) {
	uri := fmt.Sprintf("%s/api/%s/%s", env.ExtenderMarkAddr(), group, project)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type result struct {
		Code    string             `json:"code"`
		Data    []*ExtenderVersion `json:"data"`
		Message string             `json:"message"`
	}
	var respResult result
	err = json.Unmarshal(body, &respResult)
	if err != nil {
		return nil, err
	}
	return respResult.Data, err
}

func ExtenderInfoRequest(group, project, version string) (*ExtenderInfo, error) {
	uri := fmt.Sprintf("%s/api/%s/%s/%s", env.ExtenderMarkAddr(), group, project, version)
	query := url.Values{}
	query.Add("go", runtime.Version())
	query.Add("arch", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	query.Add("eosc", eosc.Version())
	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type result struct {
		Code    string          `json:"code"`
		Data    []*ExtenderInfo `json:"data"`
		Message string          `json:"message"`
	}
	var respResult result
	err = json.Unmarshal(body, &respResult)
	if err != nil {
		return nil, err
	}
	if respResult.Data == nil || len(respResult.Data) < 1 {
		if version == "latest" {
			return nil, ErrorExtenderNoLatest
		}
		return nil, ErrorExtenderNotFindMark
	}
	return respResult.Data[0], err
}