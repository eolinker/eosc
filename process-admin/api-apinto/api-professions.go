package api_apinto

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/model"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/utils"
)

func init() {
	Register(cmd.PList, ListProfession)
	Register(cmd.PGet, GetProfession)
}

func ListProfession(session ISession, message proto.IMessage) error {
	var pl []*professions.Profession
	_ = session.Call(func(adminApi admin.AdminApiWrite) error {
		pl = adminApi.ListProfession(context.Background())
		return nil
	})

	pv := utils.ArrayType(utils.ArrayType(pl, model.TypeProfessionInfo), func(t *model.ProfessionInfo) any {
		return t
	})

	session.WriteArray(pv...)
	return nil
}

func GetProfession(session ISession, message proto.IMessage) error {
	ims, err := message.Array()
	if err != nil {
		return err
	}
	if len(ims) < 2 {
		return ErrorInvalidArg
	}
	professionName, err := ims[1].String()
	if err != nil {
		return fmt.Errorf("parse profession fail %v", err)
	}
	var p *professions.Profession
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		profession, ok := adminApi.GetProfession(context.Background(), professionName)
		if !ok {
			return proto.Nil
		}
		p = profession
		return nil
	})
	if err != nil {
		return err
	}
	session.Write(p.ProfessionConfig)
	return nil
}
