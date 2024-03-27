package admin_o

import (
	"context"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils"
)

var (
	_ AdminApi         = (*imlAdminApi)(nil)
	_ AdminTransaction = (*imlAdminApi)(nil)
)

type imlAdminApi struct {
	data *imlAdminData
}

func newImlAdminApi(data *imlAdminData) *imlAdminApi {
	return &imlAdminApi{data}
}

func (a *imlAdminApi) ListWorker(ctx context.Context, profession string) ([]*WorkerInfo, error) {
	list := a.data.workers.List()
	return utils.ArrayFilter(list, func(i int, v *WorkerInfo) bool {
		return v.config.Profession == profession
	}), nil
}

func (a *imlAdminApi) GetWorker(ctx context.Context, id string) (*WorkerInfo, error) {
	workerInfo, has := a.data.workers.Get(id)
	if has {
		return workerInfo, nil
	}
	return nil, ErrorNotExist
}

func (a *imlAdminApi) DeleteWorker(ctx context.Context, id string) (*WorkerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) SetWorker(ctx context.Context, profession, name, driver, version, desc string, data IData) (*WorkerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) AllWorkers(ctx context.Context) map[string][]*WorkerInfo {
	return utils.GroupBy(a.data.workers.List(), func(v *WorkerInfo) string {
		return v.config.Profession
	})
}

func (a *imlAdminApi) GetProfession(ctx context.Context, profession string) (*professions.Profession, bool) {
	return a.data.professionData.Get(profession)
}

func (a *imlAdminApi) ListProfession(ctx context.Context, profession string) ([]*professions.Profession, error) {
	return a.data.professionData.List(), nil
}

func (a *imlAdminApi) GetSetting(ctx context.Context, name string) (any, bool) {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) SetSetting(ctx context.Context, name string, data IData) error {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) AllVariables(ctx context.Context) map[string]map[string]string {
	return a.data.variable.All()
}

func (a *imlAdminApi) GetVariables(ctx context.Context, namespace string) (map[string]string, bool) {
	values, has := a.data.variable.GetByNamespace(namespace)
	return values, has
}

func (a *imlAdminApi) GetVariable(ctx context.Context, namespace, key string) (string, bool) {
	values, has := a.data.variable.GetByNamespace(namespace)
	if !has {
		return "", false
	}
	v, h := values[key]
	return v, h
}

func (a *imlAdminApi) SetVariable(ctx context.Context, namespace string, values map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) Commit() ([]*open_api.EventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *imlAdminApi) Rollback() error {
	//TODO implement me
	panic("implement me")
}
