package eosc

import (
	"net/http"
	"testing"
	"time"
)

type TestWorker struct {
	WorkerInfo
	skill map[string]bool
}

func (t *TestWorker) Id() string {
	return t.WorkerInfo.Id
}

func (t *TestWorker) Start() error {
	return nil
}

func (t *TestWorker) Reset(conf interface{}, workers map[RequireId]interface{}) error {
	return nil
}

func (t *TestWorker) Stop() error {
	return nil
}

func (t *TestWorker) CheckSkill(skill string) bool {
	return t.skill[skill]
}

func (t *TestWorker) Info() WorkerInfo {
	return WorkerInfo{}
}

type TestWorkers struct {
	data IUntyped
}

func (t *TestWorkers) Get(id string) (IWorker, bool) {
	d, has := t.data.Get(id)
	if has {
		w, ok := d.(*TestWorker)
		if ok {
			return w, true
		}
	}
	return nil, false
}
func (t *TestWorkers) Set(id string, w *TestWorker) {
	t.data.Set(id, w)
}
func NewTestWorkers() *TestWorkers {
	return &TestWorkers{data: NewUntyped()}
}
func TestCheckConfig(t *testing.T) {
	type args struct {
		v interface{}
	}
	type TestConfig struct {
		Id        string    `json:"id"`
		name      string    `json:"name"`
		Discovery RequireId `json:"target" skill:"net/http.Handler"`
		app       http.Handler
		//Next *TestConfig `json:"next"`
	}
	workers := NewTestWorkers()
	workers.Set("001", &TestWorker{WorkerInfo: WorkerInfo{
		Id:     "001",
		Name:   "test",
		Driver: "test",
		Create: time.Now().Format(time.RFC3339),
		Update: time.Now().Format(time.RFC3339),
	}, skill: map[string]bool{
		"net/http.Handler": true,
	},
	})

	workers.Set("002", &TestWorker{WorkerInfo: WorkerInfo{
		Id:     "002",
		Name:   "test",
		Driver: "test",
		Create: time.Now().Format(time.RFC3339),
		Update: time.Now().Format(time.RFC3339),
	}, skill: map[string]bool{},
	})
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				v: &TestConfig{
					//Target: "002",
					Id:   "0001",
					name: "Test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CheckConfig(tt.args.v, workers); (err != nil) != tt.wantErr {
				t.Errorf("CheckConfig() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("CheckConfig() error = %v", err)
			}
		})
	}
}

func TestTypeNameOf(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				v: (*http.Handler)(nil),
			},
			want: "net/http.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeNameOf(tt.args.v); got != tt.want {
				t.Errorf("TypeNameOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
