package client

import (
	"context"
	"testing"
)

func createClient(t *testing.T) Client {
	client, err := New("127.0.0.1:9400")
	if err != nil {
		t.Fatal(err)
	}
	return client
}
func Test_imlClient_Ping(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ping",
			wantErr: false,
		},
	}
	client := createClient(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := client.Ping(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
