package extends

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/env"
)

type PluginVersion struct {
	*VersionInfo
	Arches []string `json:"arches"`
}

func PluginsRequest(group, project string) ([]*PluginVersion, error) {
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
		Code    string           `json:"code"`
		Data    []*PluginVersion `json:"data"`
		Message string           `json:"message"`
	}
	var respResult result
	err = json.Unmarshal(body, &respResult)
	if err != nil {
		return nil, err
	}
	return respResult.Data, err
}

func PluginInfoRequest(group, project, version string) (*PluginInfo, error) {
	uri := fmt.Sprintf("%s/api/%s/%s/%s", env.ExtenderMarkAddr(), group, project, version)
	query := url.Values{}
	query.Add("go", runtime.Version())
	query.Add("arch", runtime.GOARCH)
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
		Code    string        `json:"code"`
		Data    []*PluginInfo `json:"data"`
		Message string        `json:"message"`
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
