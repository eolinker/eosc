/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package admin

import (
	"context"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/professions"
)

type IAdminApi interface {
	GetEmployee(profession, name string) (*WorkerInfo, error)
	Export() map[string][]*WorkerInfo
	ListEmployees(profession string) ([]interface{}, error)
	GetProfession(profession string) (*professions.Profession, bool)
	Delete(id string) (*WorkerInfo, error)
	Update(profession, name, driver, version, desc string, data IData) (*WorkerInfo, error)
	Rebuild(id string) error
	//AllWorker() *WorkerInfo
	CheckDelete(ids ...string) (requires []string)
	Get(id string) (eosc.IWorker, bool)
	GetInfo(id string) (*WorkerInfo, bool)
	lock()
	unLock()
}

type IAdmin interface {
	Begin(ctx context.Context) ITransactionCtx
	GetEmployee(profession, name string) (*WorkerInfo, error)
	Export() map[string][]*WorkerInfo
	ListEmployees(profession string) ([]interface{}, error)
	GetProfession(profession string) (*professions.Profession, bool)
	CheckDelete(ids ...string) (requires []string)
	Rebuild(id string) error
}
