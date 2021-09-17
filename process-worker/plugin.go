package process_worker

import (
	"fmt"
	"path/filepath"
	"plugin"

	eosc_args "github.com/eolinker/eosc/eosc-args"

	"github.com/eolinker/eosc/log"
)

//RegisterFunc 注册函数
type RegisterFunc func()

func LoadPlugins(dir string) {

	files, err := filepath.Glob(fmt.Sprintf("%s/*.so", dir))
	if err != nil {
		log.Error(err)
		return
	}

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
		r.(RegisterFunc)()
	}
}

func loadPluginEnv() {
	dir := eosc_args.GetDefault("plugin_dir", "plugins")
	LoadPlugins(dir)
}
