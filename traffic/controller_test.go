package traffic

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/config"
	"net"
	"testing"
	"time"
)

func Test_listener(t *testing.T) {
	tcp, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   nil,
		Port: 9988,
	})
	if err != nil {
		t.Error(err)
		return
	}
	tcp.SetDeadline(time.Now().Add(time.Second * 10))
	acceptTCP, err := tcp.AcceptTCP()
	if err != nil {
		t.Log(err)
		return
	}
	acceptTCP.Close()
}

func ExampleReadController() {
	peer := config.UrlConfig{
		ListenUrls:    []string{"http://0.0.0.0", "https://0.0.0.0"},
		Certificate:   nil,
		AdvertiseUrls: nil,
	}
	clent := config.UrlConfig{
		ListenUrls:    []string{"http://0.0.0.0", "https://0.0.0.0"},
		Certificate:   nil,
		AdvertiseUrls: nil,
	}
	controller, err := ReadTraffic(nil, config.GetListens(peer, clent)...)
	if err != nil {
		fmt.Println(err)
		return
	}
	export, files := controller.Export(3)
	fmt.Println("export:", len(files))
	data, err := json.Marshal(export)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	//output:

}
