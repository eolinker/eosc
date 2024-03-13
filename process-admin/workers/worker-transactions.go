/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package workers

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
)

const (
	defaultBuffAction = 20
)

type ITransaction interface {
	Begin(ctx context.Context) ITransactionCtx
}
type TransactionEvent struct {
	ID     string
	Action string
	Data   []*WorkerInfo
}
type ITransactionCtx interface {
	Commit()
	Rollback()
	CheckDelete(ids ...string) (requires []string)
	Delete(id string) (*WorkerInfo, error)
	Update(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error)
	Rebuild(id string) error
	GetEmployee(profession, name string) (*WorkerInfo, error)
	Export() map[string][]*WorkerInfo
	ListEmployees(profession string) ([]interface{}, error)
}

var (
	_ ITransactionCtx = (*imlTransaction)(nil)
)

type imlTransaction struct {
	*imlWorkers
	ctx     context.Context
	actions []*actionContent
}

func newImlTransaction(ctx context.Context, imlWorkers *imlWorkers) *imlTransaction {
	return &imlTransaction{ctx: ctx, imlWorkers: imlWorkers, actions: make([]*actionContent, 0, defaultBuffAction)}
}

func (t *imlTransaction) Delete(id string) (*WorkerInfo, error) {
	info, err := t.imlWorkers.delete(id)
	if err != nil {
		return nil, err
	}
	t.actions = append(t.actions, newActionContent(actionDelete, id, info.config))
	return info, nil
}

func (t *imlTransaction) Update(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error) {
	body, err := data.Encode()
	if err != nil {
		return nil, err
	}
	id, ok := eosc.ToWorkerId(name, profession)
	if !ok {
		return nil, fmt.Errorf("%s@%s:invalid id", name, profession)
	}

	log.Debug("update:", id, " ", profession, ",", name, ",", driver, ",", body)
	old, exits := t.data.GetInfo(id)
	if exits {
		// update
		if driver == "" {
			driver = old.Driver()
		}

		info, err := t.imlWorkers.set(id, profession, name, driver, version, desc, body, eosc.Now(), old.config.Create)
		if err != nil {
			return nil, err
		}
		t.actions = append(t.actions, newActionContent(actionSet, id, old.config))
		return info, nil
	}
	// create
	if driver == "" {
		return nil, fmt.Errorf("require driver")
	}
	info, err := t.imlWorkers.set(id, profession, name, driver, version, desc, body, eosc.Now(), eosc.Now())
	if err != nil {
		return nil, err
	}
	t.actions = append(t.actions, newActionContent(actionCreate, id, nil))
	return info, nil

}

func (t *imlTransaction) Commit() {
	t.actions = nil
	t.imlWorkers.lockTransaction.Unlock()
}

func (t *imlTransaction) Rollback() {
	defer func() {
		t.imlWorkers.lockTransaction.Lock()
		t.actions = nil
	}()
	max := len(t.actions)
	for i := max - 1; i >= 0; i-- {
		a := t.actions[i]
		switch a.action {
		case actionCreate:
			_, err := t.imlWorkers.delete(a.id)
			if err != nil {
				log.Errorf("rollback create %s to delete error:%s", a.id, err)
				continue
			}
		case actionDelete:
			_, err := t.imlWorkers.set(a.id, a.config.Profession, a.config.Name, a.config.Driver, a.config.Version, a.config.Description, a.config.Body, a.config.Update, a.config.Create)
			if err != nil {
				log.Errorf("rollback delete %s to create error:%s", a.id, err)

				continue
			}
		case actionSet:

			_, err := t.imlWorkers.set(a.id, a.config.Profession, a.config.Name, a.config.Driver, a.config.Version, a.config.Description, a.config.Body, a.config.Update, a.config.Create)
			if err != nil {
				log.Errorf("rollback delete %s to create error:%s", a.id, err)

				continue
			}
		}
	}
}
