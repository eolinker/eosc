package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Example() {
	fmt.Println(getIps())
	//output:
}

func TestGetListens(t *testing.T) {
	type args struct {
		ucs UrlConfig
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test",
			args: args{
				ucs: UrlConfig{

					ListenUrl: ListenUrl{ListenUrls: []string{"http://0.0.0.0:8088", "http://0.0.0.0", "https://0.0.0.0", "http://192.168.0.5", "https://192.168.0.5"}},
				},
			},
			want: []string{":8088", ":80", ":443", "192.168.0.5:80", "192.168.0.5:443"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addrs := GetListens(tt.args.ucs.ListenUrl)
			sort.Strings(addrs)
			sort.Strings(tt.want)
			assert.Equalf(t, tt.want, addrs, "GetListens(%v)", tt.name)
		})
	}
}
