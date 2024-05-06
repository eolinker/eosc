package admin

import (
	"github.com/eolinker/eosc/service"
)

var (
	_ AdminApiWrite    = (*imlAdminApi)(nil)
	_ AdminTransaction = (*imlAdminApi)(nil)
)

type imlAdminApi struct {
	*imlAdminData
	actions RollbackHandlerList
	events  []*service.Event
}

func newImlAdminApi(data *imlAdminData) *imlAdminApi {
	return &imlAdminApi{imlAdminData: data}
}

func (oe *imlAdminApi) Commit() error {
	if len(oe.events) == 0 {
		oe.events = oe.events[:0]
		oe.actions = oe.actions[:0]
		oe.unLock()
		return nil
	}
	events := make([]*service.Event, len(oe.events))
	copy(events, oe.events)
	oe.events = oe.events[:0]
	oe.actions = oe.actions[:0]
	oe.unLock()

	return sendEvent(events)
}

func (oe *imlAdminApi) Rollback() error {
	defer func() {
		oe.actions = oe.actions[:0]
		oe.events = oe.events[:0]
		oe.unLock()
	}()
	return oe.actions.RollBack(oe.imlAdminData)

}
