package eosc_args

func GetDefaultArg(cfg *Config, name string, value string) string {
	vl, has := GetEnv(name)
	if has {
		return vl
	}
	return cfg.GetDefault(name, value)
}

func GetArg(cfg *Config, name string) (string, bool) {
	vl, has := GetEnv(name)
	if has {
		return vl, true
	}
	return cfg.Get(name)
}
