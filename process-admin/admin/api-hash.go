package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils"
	"github.com/eolinker/eosc/utils/hash"
	"strings"
)

func (d *imlAdminData) GetHash(ctx context.Context, key string, field string) (string, bool) {
	h, has := d.customerHash.Get(key)
	if has {
		return h.Get(field)
	}
	return "", false
}

func (d *imlAdminData) GetHashAll(ctx context.Context, key string) (stringHash, bool) {
	return d.customerHash.Get(key)
}

func (d *imlAdminData) ListHash(ctx context.Context, prefix string) []string {
	keys := d.customerHash.Keys()
	return utils.ArrayFilter(keys, func(i int, v string) bool {
		return strings.HasPrefix(v, prefix)
	})
}

func (oe *imlAdminApi) SetHash(ctx context.Context, key string, values map[string]string) error {
	if values == nil {
		return oe.DeleteHashAll(ctx, key)
	}
	ov, has := oe.customerHash.Get(key)
	oe.customerHash.Set(key, hash.NewHash(values))
	if !has {
		oe.actions = append(oe.actions, NewRollbackForAddHash(key))
	} else {
		oe.actions = append(oe.actions, NewRollBackForSetHash(key, ov))
	}
	data, _ := json.Marshal(values)
	oe.events = append(oe.events, &service.Event{
		Namespace: eosc.NamespaceCustomer,
		Command:   eosc.EventSet,
		Key:       key,
		Data:      data,
	})
	return nil
}

func (oe *imlAdminApi) SetHashValue(ctx context.Context, key string, field string, value string) error {
	ov, has := oe.customerHash.Get(key)
	var nv hash.Hash[string, string]
	var action RollbackHandler
	if !has {
		nv = hash.NewHash(make(map[string]string))
		action = NewRollbackForAddHash(key)
	} else {
		nv = ov.Clone()
		action = NewRollBackForSetHash(key, ov)
	}
	nv.Set(field, value)
	oe.customerHash.Set(key, nv)
	oe.actions = append(oe.actions, action)
	data, _ := json.Marshal(nv.Map())
	oe.events = append(oe.events, &service.Event{
		Namespace: eosc.NamespaceCustomer,
		Command:   eosc.EventSet,
		Key:       key,
		Data:      data,
	})

	return nil
}

func (oe *imlAdminApi) DeleteHash(ctx context.Context, key, field string) error {
	ovs, has := oe.customerHash.Get(key)
	if !has {
		return nil
	}

	_, has = ovs.Get(field)
	if !has {
		return nil
	}
	nvs := ovs.Clone()
	ovs.Delete(field)
	if nvs.Len() == 0 {
		// to delete
		oe.actions = append(oe.actions, NewRollBackForDeleteHash(key, ovs))
		oe.events = append(oe.events, &service.Event{
			Namespace: eosc.NamespaceCustomer,
			Command:   eosc.EventDel,
			Key:       key,
			Data:      nil,
		})
		return nil
	}
	data, _ := json.Marshal(nvs.Map())
	oe.actions = append(oe.actions, NewRollBackForSetHash(key, ovs))
	oe.events = append(oe.events, &service.Event{
		Namespace: eosc.NamespaceCustomer,
		Command:   eosc.EventSet,
		Key:       key,
		Data:      data,
	})
	return nil
}

func (oe *imlAdminApi) DeleteHashAll(ctx context.Context, key string) error {
	ov, has := oe.customerHash.Del(key)
	if !has {

		return nil
	}
	oe.actions = append(oe.actions, NewRollBackForDeleteHash(key, ov))
	oe.events = append(oe.events, &service.Event{
		Namespace: eosc.NamespaceCustomer,
		Command:   eosc.EventDel,
		Key:       key,
		Data:      nil,
	})

	return nil
}
