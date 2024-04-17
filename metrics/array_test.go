package metrics

import (
	"reflect"
	"testing"
)

type LabelReaderTest map[string]string

func (m LabelReaderTest) ReadLabel(name string) string {
	return m[name]
}

func TestParse(t *testing.T) {
	ctx := LabelReaderTest{
		"name": "test",
	}
	type args struct {
		metrics []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				metrics: []string{"name", "{name}"},
			},
			want: "name-test",
		},
		{
			name: "skip1",
			args: args{
				metrics: []string{"name", "{name}", ""},
			},
			want: "name-test",
		},
		{
			name: "skip2",
			args: args{
				metrics: []string{"name", "{name}", "{}"},
			},
			want: "name-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := ParseArray(tt.args.metrics, "-")
			got := metrics.Metrics(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
