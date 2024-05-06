package admin

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/service"
)

func (d *imlAdminData) GetProfession(ctx context.Context, profession string) (*professions.Profession, bool) {
	return d.professionData.Get(profession)
}

func (d *imlAdminData) ListProfession(ctx context.Context) []*professions.Profession {
	return d.professionData.List()
}

func (oe *imlAdminApi) SetProfession(ctx context.Context, name string, profession *eosc.ProfessionConfig) error {
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
	oe.events = append(oe.events, &service.Event{
		Command:   eosc.EventSet,
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
