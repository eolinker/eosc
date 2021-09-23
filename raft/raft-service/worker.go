package raft_service

import "github.com/eolinker/eosc"

type IWorkers interface {
	eosc.IWorkers
	IRaftServiceHandler
}
