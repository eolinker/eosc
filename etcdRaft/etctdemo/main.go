package main

import (
	"etcd_embed_demo/etcdRaft"
	"fmt"
)

func main() {
	node, err := etcdRaft.NewEtcdNode("default", nil, nil, map[string][]string{"default": {"http://default"}})
	if err != nil {
		panic(err)

	}
	node.Put("/abc", "test")
	prefix, err := node.GetPrefix("")
	if err != nil {
		panic(err)
		return
	}
	for k, v := range prefix {
		fmt.Println(k, "=", string(v))
	}
	select {}
}
