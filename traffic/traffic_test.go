package traffic

import (
	"fmt"
	"github.com/eolinker/eosc/config"

	"testing"
)

func ExampleListen() {
	addrs := []string{"http://0.0.0.0:19001", "https://0.0.0.0:19001", "tcp://"}
	tfData, err := NewTrafficData(nil).replace(config.FormatListenUrl(addrs...))
	if err != nil {
		fmt.Println("replace:", err)
		return
	}
	tfData = NewTrafficData(tfData.data)
	tf := NewTraffic(tfData)

	tcps, ssls := tf.Listen(addrs...)
	fmt.Println("tcp")
	for _, l := range tcps {

		fmt.Println(l.Addr().String())
	}
	fmt.Println("ssl")
	for _, l := range ssls {

		fmt.Println(l.Addr().String())
	}
	tfData.Shutdown()
	//output:
}

func Test_readAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "ipv4",
			args: args{
				addr: "http://0.0.0.0:8099",
			},
			want:  "0.0.0.0:8099",
			want1: false,
		},
		{
			name: "ipv6",
			args: args{
				addr: "https://[::]:8099",
			},
			want:  "[::]:8099",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readAddr(tt.args.addr)
			if got != tt.want {
				t.Errorf("readAddr() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readAddr() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
