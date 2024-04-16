package metrics

import (
	"reflect"
	"testing"
)

type LabelReaderTest map[string]string

func (m LabelReaderTest) GetLabel(name string) string {
	return m[name]
}

func TestParse(t *testing.T) {

	testreader := LabelReaderTest{
		"name":  "testReader",
		"value": "testValue",
	}

	type args = string

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: "readtest:${name}-${value}",
			want: "readtest:testReader-testValue",
		}, {
			name: "test2",
			args: "${name}-${value}xxx",
			want: "testReader-testValuexxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Parse(tt.args)
			got := m.Metrics(testreader)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
