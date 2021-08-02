package admin_open_api

import (
	"fmt"
	"os"
	"strings"

	"github.com/eolinker/eosc/log"

	"github.com/ghodss/yaml"
)

func export(value map[string][]interface{}, dir string, id string) (string, error) {
	if dir == "" {
		dir = "export"
	}
	dir = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), id)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return "", err
		}
	}

	for k, v := range value {
		newValue := map[string][]interface{}{
			k: v,
		}
		d, err := yaml.Marshal(newValue)
		if err != nil {
			log.Errorf("marshal error	%s	%s", k, err.Error())
			continue
		}
		err = writeFile(fmt.Sprintf("%s/%s.yml", dir, k), d)
		if err != nil {
			log.Errorf("write file error	%s	", k, err.Error())
			continue
		}
	}
	return dir, nil
}

func writeFile(fileName string, data []byte) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
