package extends

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc"
)

const (
	pluginInfoURI = "/plugins/info"
)

var (
	client       = &http.Client{}
	NotPluginErr = "the plugin %s is not found"
)

type PluginInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Group       string `json:"group"`
	Project     string `json:"project"`
	Version     string `json:"version"`
	Go          string `json:"go"`
	Arch        string `json:"arch"`
	Eosc        string `json:"eosc"`
	sha         string `json:"sha"`
	Status      int    `json:"status"`
	IsLatest    bool   `json:"is_latest"`
	Create      string `json:"create"`
	Update      string `json:"update"`
	URL         string `json:"url"`
}

func pluginInfoRequest(group, project, version string) (*PluginInfo, error) {
	uri := fmt.Sprintf("%s%s", env.ExtenderMarkAddr(), pluginInfoURI)
	query := url.Values{}
	query.Add("version", version)
	query.Add("group", group)
	query.Add("project", project)
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
		Code    string      `json:"code"`
		Data    *PluginInfo `json:"data"`
		Message string      `json:"message"`
	}
	var respResult result
	err = json.Unmarshal(body, &respResult)
	if err != nil {
		return nil, err
	}
	if respResult.Data == nil {
		if version == "" {
			version = "latest"
		}
		return nil, fmt.Errorf(NotPluginErr, fmt.Sprintf("%s:%s:%s", group, project, version))
	}
	return respResult.Data, err
}

//DownLoadToRepository 下载指定版本的插件项目，并解压到仓库
func DownLoadToRepository(group, project, version string) error {
	// todo 填充下载插件的代码
	// todo 插件市场地址为 ： env.ExtenderMarkAddr()
	// todo 保存目录为  filepath.Join(env.ExtendersDir(),eosc.Version(),runtime.Version(),group,project,version)
	info, err := pluginInfoRequest(group, project, version)
	if err != nil {
		return err
	}
	dest := LocalExtenderPath(group, project, version)
	tarPath := LocalExtendTarPath(group, project, version)
	resp, err := http.Get(info.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return eosc.Decompress(tarPath, dest)
}

//DownLoadToRepositoryById 下载插件， id格式为 {group}:{project}[:{version}]
func DownLoadToRepositoryById(id string) error {
	group, project, version, err := DecodeExtenderId(id)
	if err != nil {
		return err
	}
	//if version == "" {
	//	return DownLoadLatest(group, project)
	//}
	return DownLoadToRepository(group, project, version)
}

////DownLoadLatest 下载latest
//func DownLoadLatest(group, project string) error {
//	latest, err := FindLatest(group, project)
//	if err != nil {
//		return err
//	}
//	return DownLoadToRepository(group, project, latest)
//}
//
////FindLatest 查找目标项目的latest
//func FindLatest(group, project string) (string, error) {
//	// todo 填充获取插件的latest 版本的逻辑
//	// todo 插件市场地址为 ： env.ExtenderMarkAddr()
//	return "", fmt.Errorf("[%s:%s]:%w for %s-%s-%s-eosc%s", group, project, ErrorExtenderNotFindMark, runtime.GOOS, runtime.GOARCH, runtime.Version(), eosc.Version())
//}
