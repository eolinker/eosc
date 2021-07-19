package internal

import (
	"reflect"
	"testing"
)

func Test_remove(t *testing.T) {
	type args struct {
		src []string
		t   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil",
			args: args{
				src: nil,
				t:   "t",
			},
			want: nil,
		}, {
			name: "none",
			args: args{
				src: []string{"1", "2"},
				t:   "3",
			},
			want: []string{"1", "2"},
		}, {
			name: "ok",
			args: args{
				src: []string{"1", "2", "3"},
				t:   "2",
			},
			want: []string{"1", "3"},
		}, {
			name: "first",
			args: args{
				src: []string{"1", "2", "3"},
				t:   "1",
			},
			want: []string{"2","3"},
		}, {
			name: "last",
			args: args{
				src: []string{"1", "2", "3"},
				t:   "3",
			},
			want: []string{"1","2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := make([]string,len(tt.args.src))
			copy(src,tt.args.src)
			if got := remove(tt.args.src, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("remove() = %v, want %v", got, tt.want)
			}else {
				t.Logf("remove(%v,%s)=%v",src,tt.args.t,got)
			}

		})
	}
}