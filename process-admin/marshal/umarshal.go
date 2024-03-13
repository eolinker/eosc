/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package marshal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/workers"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
)

var (
	ErrorUnknownContentType = errors.New("unknown content type")
)

type JsonData []byte

func (j JsonData) String() string {
	return string(j)
}

func (j JsonData) Encode() ([]byte, error) {
	return j, nil
}

func (j JsonData) UnMarshal(v interface{}) error {
	return json.Unmarshal(j, &v)
}

func (j JsonData) Marshal() ([]byte, error) {
	return j, nil
}

type YamlData []byte

func (y YamlData) String() string {
	return string(y)
}

func (y YamlData) Encode() ([]byte, error) {
	v := make(map[string]interface{})
	err := y.UnMarshal(&v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (y YamlData) UnMarshal(v interface{}) error {
	return yaml.Unmarshal(y, v)
}

func (y YamlData) Marshal() ([]byte, error) {
	return y, nil
}

func GetData(req *http.Request) (workers.IData, error) {
	mediaType, _, err := mime.ParseMediaType(req.Header.Get("content-type"))
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(mediaType) {
	case "application/json":
		data, e := io.ReadAll(req.Body)
		if e != nil {
			return nil, e
		}
		req.Body.Close()
		log.Debug("GetData:JsonData:", string(data))
		return JsonData(data), nil
	case "application/yaml":
		data, e := io.ReadAll(req.Body)
		if e != nil {
			return nil, e
		}
		req.Body.Close()
		log.Debug("GetData:YamlData:", string(data))

		return YamlData(data), nil

	}

	return nil, fmt.Errorf("%s:%w", mediaType, ErrorUnknownContentType)
}
