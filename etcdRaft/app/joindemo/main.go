package main

import (
	"flag"
	"fmt"
	"github.com/eolinker/eosc/etcdRaft"
	"net/http"
	"time"
)

func main() {
	name := ""
	flag.StringVar(&name, "name", "node1", "")
	flag.Parse()
	var node *etcdRaft.EtcdServer
	var err error
	if name == "node1" {
		fmt.Println("node1")
		node, err = etcdRaft.NewEtcdNode("node1", []string{"http://localhost:8081"}, []string{"http://127.0.0.1:8081"}, map[string][]string{
			"node1": {"http://127.0.0.1:8081"},
		})
		if err != nil {
			panic(err)
		}
		node.Put("/abc", "test")
		node.Put("/abc1", "test1")
		node.Put("/abc2", "test2")
		http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
			node.Put(time.Now().String(), "123456")
		})
		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			prefix, err := node.GetPrefix("")
			if err != nil {
				panic(err)
			}
			for k, v := range prefix {
				w.Write([]byte(k+"="+string(v)))

			}
		})
		http.ListenAndServe(":9998", nil)
	} else {
		fmt.Println("node2")
		node, err = etcdRaft.NewEtcdNode("node2", []string{"http://localhost:8082"}, []string{"http://127.0.0.1:8082"}, map[string][]string{
			"node2": {"http://127.0.0.1:8082"},
		})
		if err != nil {
			panic(err)
		}
		node.Put("/node2", "origin")

		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			prefix, err := node.GetPrefix("")
			if err != nil {
				panic(err)
			}
			for k, v := range prefix {
				w.Write([]byte(k+"="+string(v)))

			}
		})
		http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
			err := node.Join("http://127.0.0.1:8081", []string{"http://127.0.0.1:8082"})
			if err != nil {
				w.Write([]byte(err.Error()))
			}
		})
		http.ListenAndServe(":9999", nil)
	}
	select {}
}
