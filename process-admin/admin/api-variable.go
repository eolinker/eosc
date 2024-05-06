package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
)

func (d *imlAdminData) AllVariables(ctx context.Context) map[string]map[string]string {
	return d.variable.All()
}

func (d *imlAdminData) GetVariables(ctx context.Context, namespace string) (map[string]string, bool) {
	values, has := d.variable.GetByNamespace(namespace)
	return values, has
}

func (d *imlAdminData) GetVariable(ctx context.Context, namespace, key string) (string, bool) {
	values, has := d.variable.GetByNamespace(namespace)
	if !has {
		return "", false
	}
	v, h := values[key]
	return v, h
}
func (oe *imlAdminApi) SetVariable(ctx context.Context, namespace string, values map[string]string) error {
	oldValue, has := oe.variable.GetByNamespace(namespace)

	err := oe.setVariable(namespace, values)
	if err != nil {
		return err
	}

	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	if !has {
		oe.actions = append(oe.actions, newRollbackForCreatVariable(namespace))
	} else {
		oe.actions = append(oe.actions, newRollbackForUpdateVariable(namespace, oldValue))
	}
	oe.events = append(oe.events, &service.Event{
		Command:   eosc.EventSet,
		Namespace: eosc.NamespaceVariable,
		Key:       namespace,
		Data:      data,
	})
	return nil
}
