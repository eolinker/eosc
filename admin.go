package eosc

import (
	"net/http"
)

type ProfessionInfo struct {
	Name         string   `json:"name"`
	LocalName    string   `json:"local_name"`
	Desc         string   `json:"desc"`
	Dependencies []string `json:"dependencies"`
	AppendLabels []string `json:"labels"`
	Drivers      []DriverInfo
}
type ProfessionItem struct {
}
type WorkerInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Create string `json:"create"`
	Update string `json:"update"`
}

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type IAdmin interface {
	ListProfessions() []ProfessionInfo
	ListEmployees(profession string) ([]interface{}, error)
	Update(profession, name, driver string, data IData) (*WorkerInfo, error)
	Delete(profession, name string) error
	Get(profession, name string) (interface{}, error)
	Render(profession, driver string) (*Render, error)
	Renders(profession string) (map[string]*Render, error)
	Drivers(profession string) ([]DriverInfo, error)
	DriverInfo(profession, driver string) (DriverDetail, error)
	DriversItem(profession string) ([]Item, error)
	SearchBySkill(profession string, skill []string) ([]WorkerInfo, error)
}

type IAdminHandler interface {
	GenHandler() (http.Handler, error)
}
