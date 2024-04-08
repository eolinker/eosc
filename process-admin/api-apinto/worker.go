package api_apinto

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/marshal"
)

func init() {
	Register("get", GetWorker)
	Register("list", ListWorker)
	Register("delete", DeleteWorker)
	Register("set", SetWorker)
}

type workerBase struct {
	Driver      string `json:"driver"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func SetWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 3 {
		return ErrorInvalidArg
	}
	id, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse id fail %v", err)
	}
	body, err := ims[2].String()
	if err != nil {
		return fmt.Errorf("parse body fail %v", err)
	}
	profession, name, success := eosc.SplitWorkerId(id)
	if !success {
		return fmt.Errorf("invalid id %s", id)
	}
	data := marshal.JsonData(body)
	base := new(workerBase)
	err = data.UnMarshal(base)
	if err != nil {
		return err
	}

	err = session.Call(func(adminApi admin.AdminApiWrite) error {

		_, e := adminApi.SetWorker(context.Background(), profession, name, base.Driver, base.Version, base.Description, marshal.JsonData(body))
		return e
	})
	if err != nil {
		return err
	}
	session.Write("OK")
	return nil
}

func DeleteWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 2 {
		return ErrorInvalidArg
	}
	id, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse id fail %v", err)
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		_, e := adminApi.DeleteWorker(context.Background(), id)
		return e
	})
	if err != nil {
		return err
	}
	session.Write("OK")
	return nil
}

func ListWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 2 {
		return ErrorInvalidArg
	}
	profession, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse profession fail %v", err)
	}
	var list []any

	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		listWorker, errCall := adminApi.ListWorker(context.Background(), profession)
		if errCall != nil {
			return errCall
		}
		list = make([]any, 0, len(listWorker))
		for _, worker := range listWorker {
			list = append(list, worker.Detail())
		}
		return nil
	})
	if err != nil {
		return err
	}
	session.WriteArray(list...)
	return nil
}

func GetWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 2 {
		return ErrorInvalidArg
	}
	id, err := ims[1].String()
	if err != nil {

		errArg := fmt.Errorf("parse id fail %v", err)
		return errArg
	}
	var workInfo *admin.WorkerInfo
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		w, errGet := adminApi.GetWorker(context.Background(), id)
		if errGet != nil {
			return errGet
		}
		workInfo = w
		return nil
	})
	if err != nil {
		return err
	}
	session.Write(workInfo.Detail())
	return nil
}
