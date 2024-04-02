package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc/open-api"
)

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
	oe.events = append(oe.events, &open_api.EventResponse{
		Event:     "set",
		Namespace: "variable",
		Key:       namespace,
		Data:      data,
	})
	return nil
}
