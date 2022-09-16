package setting

import (
	"reflect"
	"testing"
)

func Test_splitConfig(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		{
			name: "test",
			args: args{data: []byte(`[{"a":1,"b":2},{"a":2,"b":3}]`)},
			want: [][]byte{[]byte(`{"a":1,"b":2}`), []byte(`{"a":2,"b":3}`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitConfig(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
