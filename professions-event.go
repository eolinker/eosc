package eosc

import "fmt"

var _ IStoreEventHandler = (*Professions)(nil)

func (ps *Professions) OnDel(v StoreValue) error {
	if p, has := ps.data.get(v.Profession); has {
		if id,y:= p.delId(v.Name);y{
			if w,ok:= ps.workers.Del(id);ok{
				return w.Stop()
			}
		}
	}
	return fmt.Errorf("%s:%w", v.Profession, ErrorProfessionNotExist)
}

func (ps *Professions) OnInit(vs []StoreValue) error {

	for i := range vs {
		if e:=ps.Save(vs[i]);e!=nil{
			return e
		}
	}
	return nil

}

func (ps *Professions) OnChange(v StoreValue) error {
	return ps.Save(v)
}

func (ps *Professions)Save(v StoreValue)error  {
	if p, has := ps.data.get(v.Profession); has {
		if oid,has:= p.getId(v.Name);has && oid != v.Id{
			if w,has:=ps.workers.Del(oid);has{
				w.Stop()
			}
		}

		if err:= p.ChangeWorker(v.Driver,v.Id,v.Name,v.IData,ps.workers);err !=nil{
			return err
		}

	}
	return fmt.Errorf("%s:%w", v.Profession, ErrorProfessionNotExist)
}