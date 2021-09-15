package admin

import "github.com/eolinker/eosc/raft"

type IProfessions interface {

}

type IWorkers interface {

}

type Admin struct {
	professions IProfessions
	workers IWorkers
	raft raft.IRaft
}

func NewAdmin(professions IProfessions, workers IWorkers, raft raft.IRaft) *Admin {
	return &Admin{professions: professions, workers: workers, raft: raft}
}
