package config

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func Test_createAdvertiseUrls(t *testing.T) {
	type args struct {
		listenUrls []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "local",
			args: args{
				listenUrls: []string{"http://0.0.0.0"},
			},
			want: []string{
				"http://192.168.3.110:80",
				"http://192.168.3.114:80",
				"http://10.8.0.32:80",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			urls := createAdvertiseUrls(tt.args.listenUrls)
			sort.Strings(urls)
			sort.Strings(tt.want)
			assert.Equalf(t, tt.want, urls, "createAdvertiseUrls(%v)", tt.args.listenUrls)
		})
	}
}
