package eosc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGzipFile_DecodeData(t *testing.T) {
	type fields struct {
		Name string
		Type string
		Size int
		Data string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test",
			fields: fields{
				Name: "test.txt", Size: 4, Type: "text/plain", Data: "H4sIAAAAAAAAAytJLS4BAAx+f9gEAAAA",
			},
			want: []byte("test"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {

				if err != nil {
					t.Errorf("%s:%w", i[0], err)
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &GzipFile{
				Name: tt.fields.Name,
				Type: tt.fields.Type,
				Size: tt.fields.Size,
				Data: tt.fields.Data,
			}
			got, err := f.DecodeData()
			if !tt.wantErr(t, err, fmt.Errorf("DecodeData()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "DecodeData()")
		})
	}
}
