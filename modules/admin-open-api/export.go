package admin_open_api

import (
	"github.com/eolinker/eosc/log"
	"gopkg.in/yaml.v2"
)

func getExportData(value map[string][]interface{}) map[string][]byte {
	data := make(map[string][]byte)
	for k, v := range value {
		newValue := map[string][]interface{}{
			k: v,
		}
		d, err := yaml.Marshal(newValue)
		if err != nil {
			log.Errorf("marshal error	%s	%s", k, err.Error())
			continue
		}
		data[k] = d
	}
	return data
}

func export(value map[string][]interface{}, dir string, id string) (string, error) {
	//if dir == "" {
	//	dir = "export"
	//}
	//if id != "" {
	//	dir = fmt.Sprintf("%s/%s", strings.TrimSuffix(dir, "/"), id)
	//}
	//_, err := os.Stat(dir)
	//if err != nil {
	//	if !os.IsNotExist(err) {
	//		return "", err
	//	}
	//	err = os.MkdirAll(dir, 0777)
	//	if err != nil {
	//		return "", err
	//	}
	//}

	//for k, v := range value {
	//	newValue := map[string][]interface{}{
	//		k: v,
	//	}
	//	d, err := yaml.Marshal(newValue)
	//	if err != nil {
	//		log.Errorf("marshal error	%s	%s", k, err.Error())
	//		continue
	//	}
	//	_, err := writeFile(fmt.Sprintf("%s/%s.yml", dir, k), d)
	//	if err != nil {
	//		log.Errorf("write file error	%s	", k, err.Error())
	//		continue
	//	}
	//}
	return dir, nil
}

//func writeFile(fileName string, data []byte) ([]*os.ZipFile, error) {
//	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	_, err = file.Write(data)
//	if err != nil {
//		return err
//	}
//	return nil
//}
