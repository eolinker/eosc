package workers

import (
	"bytes"
	"encoding/gob"
)

// Cmd 用于传输的结构
type Cmd struct {
	Key    string
	Config *baseConfig
}

func (kv *Cmd) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (kv *Cmd) Decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(kv); err != nil {
		return err
	}
	return nil
}
