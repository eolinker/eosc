package api_apinto

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/marshal"
)

func init() {
	Register(cmd.SGet, GetSetting)
	Register(cmd.SSet, SetSetting)
}

func GetSetting(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 2 {
		return ErrorInvalidArg
	}
	name, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse setting name fail %v", err)
	}
	var settingConfig any
	var has bool
	_ = session.Call(func(adminApi admin.AdminApiWrite) error {
		settingConfig, has = adminApi.GetSetting(context.Background(), name)
		return nil
	})
	if !has {
		return proto.Nil
	}
	session.Write(settingConfig)
	return nil

}

func SetSetting(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 3 {
		return ErrorInvalidArg
	}
	name, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse setting name fail %v", err)
	}
	body, err := ims[2].String()
	if err != nil {
		return fmt.Errorf("parse body fail %v", err)
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {

		return adminApi.SetSetting(context.Background(), name, marshal.JsonData(body))
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil
}
