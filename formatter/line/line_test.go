package line

import (
	"reflect"
	"testing"

	"github.com/eolinker/eosc"
)

var (
	_ eosc.IEntry = (*myEntry)(nil)
)

type myEntry struct {
	data map[string]string
}

func (m *myEntry) ReadLabel(pattern string) string {
	return eosc.String(m.Read(pattern))
}

func (m *myEntry) Read(pattern string) interface{} {
	return m.data[pattern]
}

func (m *myEntry) Children(child string) []eosc.IEntry {
	//默认返回proxies
	return childEntries
}

var (
	childEntries []eosc.IEntry
)

func init() {
	childEntries = []eosc.IEntry{
		&myEntry{data: map[string]string{"proxy_username": "user1", "proxy_password": "pwd1"}},
		&myEntry{data: map[string]string{"proxy_username": "user2", "proxy_password": "pwd2"}},
	}
}

func TestLine_Format(t *testing.T) {

	type args struct {
		conf  eosc.FormatterConfig
		entry eosc.IEntry
	}
	tests := []struct {
		name   string
		fields *Line
		args   args
		want   []byte
	}{
		{
			"示例1",
			nil,
			args{
				eosc.FormatterConfig{
					"fields":  {"$id as", "@http", "@service as t", "@proxy", "@proxy#"},
					"http":    {"$request_method", "$request_uri", "@service", "@proxy", "@proxy# as proxy2"},
					"service": {"abc as service_name"},
					"proxy":   {"$proxy_username", "$proxy_password", "@abc"},
					"abc":     {"123"},
				},
				&myEntry{data: map[string]string{"id": "123", "request_method": "POST", "request_uri": "/path?a=1", "proxy_username": "user1", "proxy_password": "pwd1"}},
			},
			[]byte("123\t\"POST /path?a=1 [abc] [user1,pwd1,<123>] [<user1|pwd1|->,<user2|pwd2|->]\"\t\"abc\"\t\"user1 pwd1 [123]\"\t\"[user1,pwd1,<123>] [user2,pwd2,<123>]\""),
		}, {
			"示例2  超过第四层的不显示， 备注，在第四层能显示的只有$变量和常量，若是object或arr则显示为空字符串，依旧是用分隔符隔开",
			nil,
			args{
				eosc.FormatterConfig{
					"fields": {"@layer1"},
					"layer1": {"@layer2"},
					"layer2": {"@layer3"},
					"layer3": {"$id", "456", "@tmp", "@proxy", "@proxy#"},
					"tmp":    {"abc"},
					"proxy":  {"$proxy_username", "$proxy_password"},
				},
				&myEntry{data: map[string]string{"id": "123", "proxy_username": "user1", "proxy_password": "pwd1"}},
			},
			[]byte("\"[<123|456|-|-|->]\""),
		}, {
			"示例3  service pattern不存在的情况",
			nil,
			args{
				eosc.FormatterConfig{
					"fields": {"$id as", "@http"},
					"http":   {"$request_method", "@service", "@tmp", "@proxy# as proxy2"},
					"tmp":    {"@service", "@proxy#"},
					"proxy":  {"$proxy_username", "$proxy_password", "@abc"},
					"abc":    {"123"},
				},
				&myEntry{data: map[string]string{"id": "123", "request_method": "POST", "proxy_username": "user1", "proxy_password": "pwd1"}},
			},
			[]byte("123\t\"POST - [-,<-|->] [<user1|pwd1|->,<user2|pwd2|->]\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields, _ = NewLine(tt.args.conf)

			if got := tt.fields.Format(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() = %v \n want %v", string(got), string(tt.want))
			}

		})
	}
}
