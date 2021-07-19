package eosc

import "github.com/eolinker/eosc/internal"

var(
	_storeDriverData = internal.NewUntyped()
)
func RegisterStoreDriver(name string, factory IStoreFactory)  {
	_storeDriverData.Set(name,factory)
}
func GetStoreDriver(name string)(IStoreFactory,bool)  {
	if o, has := _storeDriverData.Get(name);has{
		return o.(IStoreFactory),true
	}
	return nil,false
}