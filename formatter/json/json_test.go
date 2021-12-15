package json

import (
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/formatter"
	"reflect"
	"testing"
)

type Entry struct {
	data     eosc.IUntyped
	children eosc.IUntyped
}

func initRootEntry() *Entry {
	var entry = genEntry("3")
	entry.children.Set("proxy1", genEntry("1"))
	entry.children.Set("proxy2", genEntry("2"))
	entry.children.Set("proxy3", genEntry("3"))
	return entry
}
func genEntry(index string) *Entry {
	var en = &Entry{
		data:     eosc.NewUntyped(),
		children: eosc.NewUntyped(),
	}
	en.data.Set("id", "123")
	en.data.Set("request_uri", "/path")
	en.data.Set("request_method", "POST")
	en.data.Set("proxy_username", "username"+index)
	en.data.Set("proxy_password", "password"+index)
	en.data.Set("error", "error"+index)
	return en
}

func (e *Entry) Read(pattern string) string {
	v, b := e.data.Get(pattern)
	if b {
		s, ok := v.(string)
		if ok {
			return s
		}
		return ""
	}
	return ""
}

func (e *Entry) Children(name string) []formatter.IEntry {
	res := make([]formatter.IEntry, 0)
	for _, child := range e.children.List() {
		c, _ := child.(formatter.IEntry)
		res = append(res, c)
	}
	return res
}

var config = formatter.Config{
	"fields": []string{
		"$id",
		"@http",
		"@service as t",
		"@tmp",
		"@proxy#errors",
	},
	"http": []string{
		"$request_method",
		"$request_uri",
		"@service",
		"@proxy",
		"@proxy# as proxy2",
	},
	"service": []string{
		"abc as service_name",
	},
	"proxy": []string{
		"$error",
		"$proxy_password",
		"$proxy_username",
	},
	"tmp": []string{
		"123",
		"456 as test",
	},
}

func Test_json_Format(t *testing.T) {

	type args struct {
		entry formatter.IEntry
	}
	tests := []struct {
		name   string
		config formatter.Config
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
		{
			name:   "general",
			config: config,
			args:   args{entry: initRootEntry()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, _ := NewFormatter(tt.config)
			got := j.Format(tt.args.entry)
			t.Log(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
