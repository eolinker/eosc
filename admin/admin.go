package admin

import "errors"

var (
	ErrorWorkerNotExist = errors.New("not exits")
)

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
