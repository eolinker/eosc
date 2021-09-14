package eosc_args

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/process"
)

var args = map[string]string{}

func init() {
	loadArgs()
}

func loadArgs() {
	// 参数配置文件格式：分行获取
	data, err := ioutil.ReadFile(fmt.Sprintf("%s.args", process.AppName()))
	if err != nil {
		log.Errorf("load args error: %s", err.Error())
		return
	}
	lineData := strings.Split(string(data), "\n")
	for _, d := range lineData {
		d = strings.TrimSpace(d)
		index := strings.Index(d, "=")
		if index == -1 {
			args[d] = ""
			continue
		}
		args[d[:index]] = d[index+1:]
	}
	return
}

func WriteArgsToFile() error {
	builder := strings.Builder{}
	for key, value := range args {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
		builder.WriteString("\n")
	}
	f, err := os.OpenFile(fmt.Sprintf("%s.args", process.AppName()), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// offset
	_, err = f.Write([]byte(builder.String()))
	if err != nil {
		return err
	}
	log.Info("write args file succeed!")

	return nil
}
