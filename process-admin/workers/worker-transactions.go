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
}

func (i *imlTransaction) CheckDelete(ids ...string) (requires []string) {
	//TODO implement me
	panic("implement me")
}

func (i *imlTransaction) Delete(id string) (*WorkerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (i *imlTransaction) Update(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (i *imlTransaction) Commit() {
	//TODO implement me
	panic("implement me")
}

func (i *imlTransaction) Rollback() {
	//TODO implement me
	panic("implement me")
}
