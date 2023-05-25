package extends

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/fileLocker"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/service"
	"google.golang.org/protobuf/proto"
	"os"
)

func DownloadCheck(group string, project string, version string) error {
	path := LocalExtenderPath(group, project, version)
	err := os.MkdirAll(path, 0666)
	if err != nil {
		return errors.New("create extender path " + path + " error: " + err.Error())
	}
	locker := fileLocker.NewLocker(LocalExtenderPath(group, project, version), 30, fileLocker.CliLocker)
	err = locker.TryLock()
	if err != nil {
		return errors.New("locker error: " + err.Error())
	}
	defer locker.Unlock()
	err = DownLoadToRepositoryById(FormatDriverId(group, project, version))
	if err != nil {
		return errors.New("download extender to local error: " + err.Error())
	}
	return nil
}
func CheckExtends(exts ...string) ([]*service.ExtendsInfo, []*service.ExtendsBasicInfo, error) {
	return checkExtends(exts)
}
func CheckExtender(group, project, version string) (*service.ExtendsInfo, error) {
	id := FormatFileName(group, project, version)
	info, faile, err := CheckExtends(id)

	if err != nil {
		return nil, err
	}
	if len(faile) > 0 {
		return nil, errors.New(faile[0].Msg)
	}
	if len(info) == 0 {
		return nil, errors.New("unknown error")
	}
	return info[0], nil
}
func checkExtends(exts []string) ([]*service.ExtendsInfo, []*service.ExtendsBasicInfo, error) {
	data, err := json.Marshal(exts)
	if err != nil {
		return nil, nil, err
	}

	cmd, err := process.Cmd(eosc.ProcessHelper, nil)
	if err != nil {
		return nil, nil, err
	}
	cmd.Stdin = bytes.NewReader(data)
	buff := &bytes.Buffer{}
	cmd.Stdout = buff
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	cmd.Wait()

	response := new(service.ExtendsResponse)

	err = proto.Unmarshal(buff.Bytes(), response)
	if err != nil {
		return nil, nil, err
	}
	if response.Code != "000000" {
		return nil, nil, errors.New(response.Msg)
	}
	return response.Extends, response.FailExtends, nil
}
