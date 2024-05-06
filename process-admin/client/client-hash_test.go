package client

import (
	"context"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	ctx := context.Background()

	var client HClient = createClient(t)
	matchKeys, err := client.HMatch(ctx, "*")
	if err != nil {
		return
	}
	for _, key := range matchKeys {
		t.Log("clear:", key)
		err := client.HDelAll(ctx, key)
		if err != nil {
			t.Error(err)
			return
		}
	}
	err = client.HSet(ctx, "subscription:1", "field1", "ok")
	if err != nil {
		t.Errorf("HSet err:%v", err)
		return
	}
	value, err := client.HGet(ctx, "subscription:1", "field1")
	if err != nil {
		t.Errorf("HGet err:%v", err)
		return
	}
	if value != "ok" {
		t.Errorf("HGet subscription:1 field1 want: ok but got:%v", value)
		return
	}
	t.Log("HSet OK")
	mapValue := map[string]string{
		"field1": "1",
		"field2": "2",
		"field3": "3",
	}
	err = client.HMSet(ctx, "sub:2", mapValue)
	if err != nil {
		t.Errorf("HMSet err:%v", err)
		return
	}
	value, err = client.HGet(ctx, "sub:2", "field1")
	if err != nil {
		t.Errorf("HGet err:%v", err)
		return
	}
	if value != "1" {
		t.Errorf("HGet sub:2 field1 want:1 but got:%v", value)
		return
	}
	t.Log("HMSet AND GET OK")
	all, err := client.HGetAll(ctx, "sub:2")
	if err != nil {
		t.Errorf("HGetAll sub:2 err:%v", err)
		return
	}

	if !reflect.DeepEqual(all, mapValue) {
		t.Errorf("HGetAll sub:2  want:%v but got:%v", mapValue, all)
		return
	}
	t.Log("HGetAll OK")

	match, err := client.HMatch(ctx, "*")
	if err != nil {
		t.Errorf("HMatch err:%v", err)
		return
	}
	t.Log("match:", match)
}
