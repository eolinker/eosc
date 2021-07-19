package store

import "github.com/eolinker/eosc"

func Register()  {
	eosc.RegisterStoreDriver("yaml",new(Factory))
}
