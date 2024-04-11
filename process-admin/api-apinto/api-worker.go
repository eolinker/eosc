package api_apinto

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/utils"
)

func init() {
	Register(cmd.WorkerGet, GetWorker)
	Register(cmd.WorkerList, ListWorker)
	Register(cmd.WorkerMatch, MatchWorker)
	Register(cmd.WorkerDel, DeleteWorker)
	Register(cmd.WorkerSet, SetWorker)
}

func SetWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 3 {
		return ErrorInvalidArg
	}
	var id string
	var body []byte
	err = ims[1:].Scan(&id, &body)
	if err != nil {
		return err
	}

	profession, name, success := eosc.SplitWorkerId(id)
	if !success {
		return fmt.Errorf("invalid id %s", id)
	}
	data := marshal.JsonData(body)
	cf := new(eosc.WorkerConfig)
	err = data.UnMarshal(cf)
	if err != nil {
		return err
	}
	cf.Profession = profession
	cf.Id = id
	cf.Name = name
	cf.Body = data
	err = session.Call(func(adminApi admin.AdminApiWrite) error {

		_, e := adminApi.SetWorker(context.Background(), cf)
		return e
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
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
	var id string
	err = ims[1:].Scan(&id)
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
	session.Write(cmd.OK)
	return nil
}

type MatchLabel struct {
	labels map[string]string
}

func (m *MatchLabel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m.labels)
}

func MatchWorker(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 3 {
		return ErrorInvalidArg
	}
	var profession string
	var matchLabels MatchLabel
	err = ims[1:].Scan(&profession, &matchLabels)
	if err != nil {
		return err
	}
	var list []any
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		listWorker, errCall := adminApi.ListWorker(context.Background(), profession)
		if errCall != nil {
			return errCall
		}

		utils.ArrayFilter(listWorker, func(i int, v *admin.WorkerInfo) bool {
			vl := v.MatchLabels()
			if vl != nil {
				return false
			}
			for label, value := range matchLabels.labels {
				if vl[label] != value {
					return false
				}
			}
			return true
		})

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
