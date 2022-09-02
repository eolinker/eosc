package traffic

import (
	"fmt"
	"net"
	"reflect"
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
func Test_rebuildAddr(t *testing.T) {
	type args struct {
		addrs []*net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want map[int][]net.IP
	}{
		{
			name: "nil",
			args: args{
				addrs: []*net.TCPAddr{{
					IP:   net.ParseIP(""),
					Port: 80,
				}},
			},
			want: map[int][]net.IP{80: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rebuildAddr(tt.args.addrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rebuildAddr() = %v, want %v", got, tt.want)
			} else {
				for p, ips := range got {
					fmt.Printf("port:%d=>%v", p, ips)
				}
			}

		})
	}
}
