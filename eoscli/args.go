package eoscli

import (
	eosc_args "github.com/eolinker/eosc/eosc-args"
)

func getDefaultArg(cfg *eosc_args.Config, name string, value string) string {
	vl, has := eosc_args.GetEnv(name)
	if has {
		return vl
	}
	return cfg.GetDefault(name, value)
}

func getArg(cfg *eosc_args.Config, name string) (string, bool) {
	vl, has := eosc_args.GetEnv(name)
	if has {
		return vl, true
	}
	return cfg.Get(name)
}
