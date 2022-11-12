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

func ExampleReadController2() {
	addrs := []string{"http://0.0.0.0:19001", "https://0.0.0.0:19001", "tcp://:9988"}
	traffic, err := ReadTraffic(nil, config.FormatListenUrl(addrs...)...)
	if err != nil {
		fmt.Println(err)
		return
	}
	export, files := Export(traffic, 3)
	fmt.Println("export:", len(files))
	data, err := json.Marshal(export)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	//output:
}
func ExampleReadController() {
	peer := config.UrlConfig{
		ListenUrl: config.ListenUrl{
			ListenUrls:    []string{"http://0.0.0.0", "https://0.0.0.0"},
			AdvertiseUrls: nil,
		},
		Certificate: nil,
	}
	clent := config.UrlConfig{
		ListenUrl: config.ListenUrl{
			ListenUrls:    []string{"http://0.0.0.0", "https://0.0.0.0"},
			AdvertiseUrls: nil,
		},
		Certificate: nil,
	}
	traffic, err := ReadTraffic(nil, config.GetListens(peer.ListenUrl, clent.ListenUrl)...)
	if err != nil {
		fmt.Println(err)
		return
	}
	export, files := Export(traffic, 3)
	fmt.Println("export:", len(files))
	data, err := json.Marshal(export)
	if err != nil {
		return
	}
	fmt.Println(string(data))
	//output:

}
