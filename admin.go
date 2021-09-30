package eosc

import (
	"net/http"
)

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type IAdmin interface {
	IAdminWorker
	IAdminPermission
}
type IAdminWorker interface {
	ListEmployees(profession string) ([]interface{}, error)
	//ListEmployeeNames(profession string) ([]string, error)
	Update(profession, name, driver string, data []byte) (interface{}, error)
	Delete(profession, name string) (interface{}, error)
	GetEmployee(profession, name string) (interface{}, error)
}

type IAdminPermission interface {
	Drivers(profession string) ([]*DriverInfo, error)
	DriverInfo(profession, driver string) (*DriverDetail, error)
	DriversItem(profession string) ([]*Item, error)
	ListProfessions() []*ProfessionInfo
}

type IAdminHandler interface {
	GenHandler() (http.Handler, error)
}
