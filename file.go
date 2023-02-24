package eosc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

type GzipFile struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Data string `json:"data"`
}

type EoFiles []GzipFile

func (f *GzipFile) DecodeData() ([]byte, error) {
	d, err := base64.StdEncoding.DecodeString(f.Data)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(bytes.NewReader(d))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
