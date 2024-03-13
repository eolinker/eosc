package json

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/eolinker/eosc"
)

type Entry struct {
	data     eosc.Untyped[string, string]
	children eosc.Untyped[string, eosc.IEntry]
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
		data:     eosc.BuildUntyped[string, string](),
		children: eosc.BuildUntyped[string, eosc.IEntry](),
	}
	en.data.Set("id", "123")
	en.data.Set("request_uri", "/path")
	en.data.Set("request_method", "POST")
	en.data.Set("proxy_username", "username"+index)
	en.data.Set("proxy_password", "password"+index)
	en.data.Set("error", "error"+index)
	return en
}

func (e *Entry) Read(pattern string) interface{} {
	s, b := e.data.Get(pattern)
	if b {
		return s
	}
	return ""
}

func (e *Entry) Children(name string) []eosc.IEntry {
	return e.children.List()
}

var config = eosc.FormatterConfig{
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
		entry eosc.IEntry
	}
	tests := []struct {
		name   string
		config eosc.FormatterConfig
		args   args
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
		{
			name:   "general",
			config: config,
			args:   args{entry: initRootEntry()},
			want: map[string]interface{}{
				"id": "123",
				"http": map[string]interface{}{
					"request_method": "POST",
					"request_uri":    "/path",
					"service": map[string]interface{}{
						"service_name": "abc",
					},
					"proxy": map[string]interface{}{
						"error":          "error3",
						"proxy_password": "password3",
						"proxy_username": "username3",
					},
					"proxy2": []interface{}{
						map[string]interface{}{
							"error":          "error1",
							"proxy_password": "password1",
							"proxy_username": "username1",
						},
						map[string]interface{}{
							"error":          "error2",
							"proxy_password": "password2",
							"proxy_username": "username2",
						},
						map[string]interface{}{
							"error":          "error3",
							"proxy_password": "password3",
							"proxy_username": "username3",
						},
					},
				},
				"t": map[string]interface{}{
					"service_name": "abc",
				},
				"tmp": map[string]interface{}{
					"123":  "123",
					"test": "456",
				},
				"proxy": []interface{}{
					map[string]interface{}{
						"error":          "error1",
						"proxy_password": "password1",
						"proxy_username": "username1",
					},
					map[string]interface{}{
						"error":          "error2",
						"proxy_password": "password2",
						"proxy_username": "username2",
					},
					map[string]interface{}{
						"error":          "error3",
						"proxy_password": "password3",
						"proxy_username": "username3",
					},
				},
			},
		},
	}
	bodyConfig := []contentResize{
		{
			Size:   1024,
			Suffix: "kb",
		},
		{
			Size:   1024 * 1024,
			Suffix: "mb",
		},
		{
			Size:   1024 * 1024 * 1024,
			Suffix: "gb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantData, _ := json.Marshal(tt.want)

			j, _ := NewFormatter(tt.config, bodyConfig)
			got := j.Format(tt.args.entry)
			gotObj := map[string]interface{}{}
			err := json.Unmarshal(got, &gotObj)
			if err != nil {
				t.Errorf("Format() = %s error:%v", string(got), err)
				return
			}
			t.Log(string(got))
			if !reflect.DeepEqual(gotObj, tt.want) {
				t.Errorf("Format() = %s, want %s", string(got), string(wantData))
			}
		})
	}
}
