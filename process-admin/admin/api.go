package admin

import (
	open_api "github.com/eolinker/eosc/open-api"
)

var (
	_ AdminApiWrite    = (*imlAdminApi)(nil)
	_ AdminTransaction = (*imlAdminApi)(nil)
)

type imlAdminApi struct {
	*imlAdminData
	actions RollbackHandlerList
	events  []*open_api.EventResponse
}

func newImlAdminApi(data *imlAdminData) *imlAdminApi {
	return &imlAdminApi{imlAdminData: data}
}

func (oe *imlAdminApi) Commit() ([]*open_api.EventResponse, error) {

	events := make([]*open_api.EventResponse, len(oe.events))
	copy(events, oe.events)
	oe.events = oe.events[:0]
	oe.actions = oe.actions[:0]
	oe.unLock()
	return events, nil
}

func (oe *imlAdminApi) Rollback() error {
	defer func() {
		oe.actions = oe.actions[:0]
		oe.events = oe.events[:0]
		oe.unLock()
	}()
	return oe.actions.RollBack(oe.imlAdminData)

}
