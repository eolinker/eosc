package eosc

import (
	"net/http"
	"testing"
	"time"
)

type TestEmployee struct {
	WorkerInfo
	skill map[string]bool
}

func (t *TestEmployee) Marshal() ([]byte, error) {
	return nil, nil
}

func (t *TestEmployee) Worker() (IWorker, error) {
	return nil, nil
}

func (t *TestEmployee) CheckSkill(skill string) bool {
	return t.skill[skill]
}

func (t *TestEmployee) Info() WorkerInfo {
	return WorkerInfo{}
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

	globalIEmployees.Set("001", &TestEmployee{WorkerInfo: WorkerInfo{
		Id:     "001",
		Name:   "test",
		Driver: "test",
		Create: time.Now().Format(time.RFC3339),
		Update: time.Now().Format(time.RFC3339),
	}, skill: map[string]bool{
		"net/http.Handler": true,
	},
	})

	globalIEmployees.Set("002", &TestEmployee{WorkerInfo: WorkerInfo{
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
			if err := CheckConfig(tt.args.v); (err != nil) != tt.wantErr {
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
