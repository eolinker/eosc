package admin

import (
	"github.com/eolinker/eosc/admin"
	"github.com/eolinker/eosc/raft"
)

type WorkerInfo map[string]interface{}

type Admin struct {
	professions admin.IProfessions
	workers     admin.IWorkers
	raft        raft.IRaftSender
}

func NewAdmin(professions admin.IProfessions, workers admin.IWorkers, raft raft.IRaft) *Admin {
	return &Admin{professions: professions, workers: workers, raft: raft}
}
