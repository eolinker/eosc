package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"testing"
)

func Test_imlClient_List(t *testing.T) {

	client := createClient(t)
	tests := []struct {
		name       string
		profession string

		wantErr bool
	}{
		{
			name:       "router",
			profession: "router",
		}, {
			name:       "service",
			profession: "service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.List(context.Background(), tt.profession)
			if (err != nil && !errors.Is(err, proto.Nil)) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("list [%s] get %d ", tt.profession, len(got))
			for i, v := range got {
				t.Logf("== %d ==> %s", i, string([]byte(v)))
			}
		})
	}
}

func Test_imlClient_Get(t *testing.T) {
	client := createClient(t)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "service",
			id:      "test22@service",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.Get(context.Background(), tt.id)
			if (err != nil && !errors.Is(err, proto.Nil)) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Get %s ==>%s", tt.id, string(got))
		})
	}
}

func Test_imlClient_Set(t *testing.T) {
	client := createClient(t)
	type args struct {
		ctx   context.Context
		id    string
		value any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set",
			args: args{
				ctx:   context.Background(),
				id:    "test22@service",
				value: "test22@service",
			},
			wantErr: false,
		}, {
			name: "set",
			args: args{
				ctx:   context.Background(),
				id:    "test22@service",
				value: `{"balance":"round-robin","create":"2024-04-08 18:17:31","description":"","driver":"http","id":"test22@service","name":"test22","nodes":["demo.apinto.com:8280 weight=10"],"pass_host":"node","profession":"service","retry":0,"scheme":"HTTP","timeout":1000,"update":"2024-04-08 18:17:31","version":"20230605160409"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := client.Set(tt.args.ctx, tt.args.id, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_imlClient_PList(t *testing.T) {
	client := createClient(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "list profession",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.PList(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("PList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("list profession got:%d", len(got))
			for i, v := range got {
				body, _ := json.Marshal(v)
				t.Logf("== %d %s==> %s\n", i, v.Name, string(body))
			}
			t.Logf("profession list done")

		})
	}
}

func Test_imlClient_PGet(t *testing.T) {
	client := createClient(t)
	tests := []struct {
		name string

		wantErr bool
	}{
		{
			name:    "service",
			wantErr: false,
		}, {
			name:    "router",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.PGet(context.Background(), tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			data, _ := json.Marshal(got)
			t.Logf("%s got =>%s", tt.name, string(data))
		})
	}
}
