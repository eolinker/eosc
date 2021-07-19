package eosc

import (
	"reflect"
	"sync"
)
type StoreEventDispatcher struct {
	locker sync.RWMutex
	handlers []IStoreEventHandler
}

func NewStoreDispatcher() *StoreEventDispatcher {
	return &StoreEventDispatcher{
		locker:sync.RWMutex{},
	}
}

func (s *StoreEventDispatcher) AddListen(h IStoreEventHandler)bool {

	v :=reflect.ValueOf(h)
	if v.IsNil(){
		return false
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	for _,i:=range  s.handlers{
		if  reflect.DeepEqual(i,h){
			return false
		}
	}
	s.handlers = append(s.handlers, h)
	return true
}

func (s *StoreEventDispatcher) DispatchDel(v StoreValue)error  {
	s.locker.RLock()
	defer s.locker.RUnlock()

	for _,i:=range  s.handlers{

		if err:=i.OnDel(v);err!=nil{
			return err
		}
	}
	return nil
}

func (s *StoreEventDispatcher)DispatchChange(v StoreValue) error{
	s.locker.RLock()
	defer s.locker.RUnlock()

	for _,i:=range  s.handlers{

		if err:=i.OnChange(v);err!=nil{
			return err
		}
	}
	return nil
}

func (s *StoreEventDispatcher) DispatchInit(vs []StoreValue) error {
	s.locker.RLock()
	defer s.locker.RUnlock()

	for _,i:=range  s.handlers{

		if err:=i.OnInit(vs);err!=nil{
			return err
		}
	}
	return nil
}