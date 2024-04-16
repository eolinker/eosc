package client

import (
	"context"
	"testing"
)

func Test_imlClient_SGet(t *testing.T) {
	client := createClient(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "plugins",
			wantErr: false,
		}, {
			name:    "plugin",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.SGet(context.Background(), tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//data, _ := json.Marshal(got)
			t.Logf("sget %s=>%s", tt.name, string(got))
		})
	}
}

func Test_imlClient_SSet(t *testing.T) {
	client := createClient(t)
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{
			name: "plugin",
			value: `{"plugins":[
        {
            "id": "eolinker.com:apinto:access_log",
            "name": "access_log",
            "status": "global"
        },
        {
            "id": "eolinker.com:apinto:proxy_rewrite_v2",
            "name": "proxy_rewrite",
            "status": "enable"
        },
        {
            "id": "eolinker.com:apinto:extra_params",
            "name": "extra_params",
            "status": "enable"
        },
        {
            "id": "eolinker.com:apinto:plugin_app",
            "name": "app",
            "status": "global"
        },
        {
            "id": "eolinker.com:apinto:strategy-plugin-visit",
            "name": "strategy_visit",
            "status": "global"
        },
        {
            "id": "eolinker.com:apinto:strategy-plugin-grey",
            "name": "strategy_grey",
            "status": "global"
        },
        {
            "config": {
                "cache": "redis@output"
            },
            "id": "eolinker.com:apinto:strategy-plugin-limiting",
            "name": "strategy_limiting",
            "status": "global"
        },
        {
            "config": {
                "cache": "redis@output"
            },
            "id": "eolinker.com:apinto:strategy-plugin-fuse",
            "name": "strategy_fuse",
            "status": "global"
        },
        {
            "config": {
                "cache": "redis@output"
            },
            "id": "eolinker.com:apinto:strategy-plugin-cache",
            "name": "strategy_cache",
            "status": "global"
        }
    ]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := client.SSet(context.Background(), tt.name, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := client.SGet(context.Background(), tt.name)
			if err != nil {
				t.Errorf("SGet error = %v", err)
				return
			}
			t.Logf("SGet %s=> %s", tt.name, string(got))

		})
	}
}
