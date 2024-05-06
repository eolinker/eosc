package api_apinto

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/utils"
	"github.com/eolinker/eosc/utils/hash"
)

func init() {
	Register(cmd.HSet, HashSet)
	Register(cmd.HGet, HashGet)
	Register(cmd.HDel, HashDel)
	Register(cmd.HDelAll, HashDelAll)
	Register(cmd.HGetAll, HashGetAll)
	Register(cmd.HExists, HashExists)
	Register(cmd.HMSet, HashMSet)
	Register(cmd.HRest, HashRest)
	Register(cmd.HKeys, HashKeys)
	Register(cmd.HMatch, HashMatch)
}

// HashRest 重置一个key, 与 HMSET的区别是会移除不在此请求中的field, 如果 HRESET key , 则等效于 HDELALL
// 格式为  HREST key [field value]...
func HashRest(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}

	key, err := msgs[1].String()
	if err != nil {
		return err
	}
	fvms := msgs[2:]
	kv := make(map[string]string)
	for len(fvms) >= 2 {
		fv := fvms[:2]
		fvms = fvms[2:]
		var f, v string
		err := fv.Scan(&f, &v)
		if err != nil {
			return err
		}
		kv[f] = v
	}
	if len(kv) == 0 {
		errDelete := session.Call(func(adminApi admin.AdminApiWrite) error {
			return adminApi.DeleteHashAll(context.Background(), key)
		})
		if errDelete != nil {
			return errDelete
		}
		session.Write(cmd.OK)
		return nil
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		return adminApi.SetHash(context.Background(), key, kv)
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil
}

// HashDelAll 删除指定的key
func HashDelAll(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	key, err := msgs[1].String()
	if err != nil {
		return err
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		return adminApi.DeleteHashAll(context.Background(), key)
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil

}

// HashMatch 用前缀匹配查找所有存在的hashmap,返回 key数组
// HMATCH keyPrefix
// 特例: HATCH * 会返回所有的 key
func HashMatch(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return nil
	}
	matchKey, err := msgs[1].String()
	if err != nil {
		return err
	}
	var keys []string
	_ = session.Call(func(adminApi admin.AdminApiWrite) error {
		keys = adminApi.ListHash(context.Background(), matchKey)
		return nil
	})
	session.WriteArray(utils.ArrayType(keys, func(t string) any { return t })...)
	return nil

}

// HashKeys 查找指定 hash的field列表
// HKEYS key
func HashKeys(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	key, err := msgs[1].String()
	if err != nil {
		return err
	}
	var hash hash.Hash[string, string]
	var has bool
	_ = session.Call(func(adminApi admin.AdminApiWrite) error {
		hash, has = adminApi.GetHashAll(context.Background(), key)
		return nil
	})
	if !has {
		session.Write(nil)
		return nil
	}
	keys := hash.Keys()

	session.WriteArray(utils.ArrayType(keys, func(t string) any {
		return t
	})...)
	return nil
}

// HashMSet 批量设置 key的 kv值
// HMSET key [field value]...
func HashMSet(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}

	key, err := msgs[1].String()
	if err != nil {
		return err
	}
	fvms := msgs[2:]
	kv := make(map[string]string)
	for len(fvms) >= 2 {
		fv := fvms[:2]
		fvms = fvms[2:]
		var f, v string
		err := fv.Scan(&f, &v)
		if err != nil {
			return err
		}
		kv[f] = v
	}
	if len(kv) == 0 {
		session.Write(cmd.OK)
		return nil
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		if len(kv) == 1 {
			for k, v := range kv {
				return adminApi.SetHashValue(context.Background(), key, k, v)
			}
		}
		ov, has := adminApi.GetHashAll(context.Background(), key)
		if !has {
			return adminApi.SetHash(context.Background(), key, kv)
		}
		for k, v := range kv {
			ov.Set(k, v)
		}
		return adminApi.SetHash(context.Background(), key, ov.Map())
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil

}

// HashExists 判断指定的 key 下的field 是否存在
// HEXISTS key field
func HashExists(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	if len(msgs) < 3 {
		return fmt.Errorf("field require")
	}
	var key, field string
	err = msgs[1:].Scan(&key, &field)
	if err != nil {
		return err
	}
	var exists bool
	_ = session.Call(func(adminApi admin.AdminApiWrite) error {
		_, exists = adminApi.GetHash(context.Background(), key, field)
		return nil
	})
	session.Write(exists)
	return nil
}

// HashGetAll 获取指定的hash的所有值
// HGETALL key
// 返回值为依次field,value的数组, 形如 field1,value1,field2,value2 ...
func HashGetAll(session ISession, message proto.IMessage) error {
	array, err := message.Array()
	if err != nil {
		return err
	}
	if len(array) < 2 {
		return fmt.Errorf("key require")
	}
	key, err := array[1].String()
	if err != nil {
		return err
	}
	var values hash.Hash[string, string]
	var has bool
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		values, has = adminApi.GetHashAll(context.Background(), key)
		return nil
	})
	if err != nil {
		return err
	}
	if has {
		m := values.Map()
		rs := make([]interface{}, 0, len(m)*2)
		for k, v := range m {
			rs = append(rs, k, v)
		}
		session.WriteArray(rs...)
	} else {
		session.WriteArray()
	}
	return nil

}

// 删除指定key下的多个域
func HashDel(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	if len(msgs) < 3 {
		return fmt.Errorf("field require")
	}
	key, err := msgs[1].String()
	if err != nil {
		return err
	}
	fieldsMsg := msgs[2:]
	fields := make([]string, 0, len(fieldsMsg))
	for _, m := range fieldsMsg {
		fv, err := m.String()
		if err != nil {
			return err
		}
		fields = append(fields, fv)
	}

	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		if len(fields) == 1 {
			return adminApi.DeleteHash(context.Background(), key, fields[0])
		}
		values, has := adminApi.GetHashAll(context.Background(), key)
		if !has {
			return nil
		}
		for _, f := range fields {
			values.Delete(f)
		}
		return adminApi.SetHash(context.Background(), key, values.Map())
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil
}

func HashGet(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}
	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	if len(msgs) < 3 {
		return fmt.Errorf("field require")
	}

	var key, field string
	err = msgs[1:].Scan(&key, &field)
	if err != nil {
		return err
	}
	var v string
	var has bool
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		v, has = adminApi.GetHash(context.Background(), key, field)
		return nil
	})
	if err != nil {
		return err
	}
	if has {
		session.Write(v)
	} else {
		session.Write(nil)
	}
	return nil
}

func HashSet(session ISession, message proto.IMessage) error {
	msgs, err := message.Array()
	if err != nil {
		return err
	}

	if len(msgs) < 2 {
		return fmt.Errorf("key require")
	}
	if len(msgs) < 3 {
		return fmt.Errorf("field require")
	}
	if len(msgs) < 4 {
		return fmt.Errorf("value require")
	}
	var key, field, value string
	err = msgs[1:].Scan(&key, &field, &value)
	if err != nil {
		return err
	}
	err = session.Call(func(adminApi admin.AdminApiWrite) error {
		return adminApi.SetHashValue(context.Background(), key, field, value)
	})
	if err != nil {
		return err
	}
	session.Write(cmd.OK)
	return nil
}
