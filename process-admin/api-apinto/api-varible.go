package api_apinto

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
)

func init() {
	Register(cmd.VGet, GetVariable)
	Register(cmd.VSet, SetVariable)
}

func SetVariable(session ISession, message proto.IMessage) error {
	array, err := message.Array()
	if err != nil {
		return err
	}
	if len(array) < 2 {
		return errors.New("invalid arg")
	}
	var namespace string
	var body []byte
	err = array.Scan(&namespace, &body)
	if err != nil {
		return err
	}
	variables := make(map[string]string)
	err = json.Unmarshal(body, &variables)
	if err != nil {
		return err
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		return adminApi.SetVariable(context.Background(), namespace, variables)
	})
	if err != nil {
		return nil
	}
	session.Write(cmd.OK)
	return nil
}

func GetVariable(session ISession, message proto.IMessage) error {
	array, err := message.Array()
	if err != nil {
		return err
	}
	switch len(array) {
	case 0:
		return proto.Nil
	case 1:
		var variables map[string]map[string]string
		_ = session.Call(func(adminApi admin.AdminApiWrite) error {
			variables = adminApi.AllVariables(context.Background())
			return nil
		})
		session.Write(variables)
	case 2:
		var cmd string
		var namespace string
		err := array.Scan(&cmd, &namespace)
		if err != nil {
			return err
		}
		var values map[string]string
		var has bool
		_ = session.Call(func(adminApi admin.AdminApiWrite) error {
			values, has = adminApi.GetVariables(context.Background(), namespace)
			return nil
		})
		if !has {
			return proto.Nil
		}
		session.Write(values)
	default:
		var cmd string
		var namespace string
		var key string
		err := array.Scan(&cmd, &namespace, &key)
		if err != nil {
			return err
		}
		var value string
		var has bool
		_ = session.Call(func(adminApi admin.AdminApiWrite) error {
			value, has = adminApi.GetVariable(context.Background(), namespace, key)
			return nil
		})
		if !has {
			return proto.Nil
		}
		session.Write(value)
	}
	return nil
}
