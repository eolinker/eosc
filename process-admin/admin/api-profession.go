package admin

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/open-api"
)

func (oe *imlAdminApi) SetProfession(name string, profession *eosc.ProfessionConfig) error {
	old, has := oe.professionData.Get(name)
	err := oe.imlAdminData.setProfession(name, profession)
	if err != nil {
		return err
	}
	if has {
		oe.actions = append(oe.actions, newRollbackForSetProfession(name, old.ProfessionConfig))
	} else {
		oe.actions = append(oe.actions, newRollbackForAddProfession(name))
	}
	data, _ := json.Marshal(profession)
	oe.events = append(oe.events, &open_api.EventResponse{
		Event:     eosc.EventSet,
		Namespace: eosc.NamespaceProfession,
		Key:       name,
		Data:      data,
	})
	return nil
}

func (oe *imlAdminApi) ResetProfession(configs []*eosc.ProfessionConfig) {
	//TODO implement me
	panic("implement me")
}
