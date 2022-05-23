package extender

import (
	"fmt"
	"time"
)

type Item struct {
	Group      string
	Project    string
	Version    string
	Status     int
	NextTime   time.Time
	RetryCount int
}

func (e *Item) Key() string {
	return fmt.Sprint(e.Group, ":", e.Project, ":", e.Version)
}
func NewItem(group, project, version string) *Item {
	return &Item{
		Group:      group,
		Project:    project,
		Version:    version,
		Status:     StatusInit,
		NextTime:   time.Now(),
		RetryCount: 0,
	}
}

func (e *Item) ToStatus() *Status {
	return &Status{
		Group:   e.Group,
		Project: e.Project,
		Version: e.Version,
		Status:  e.Status,
	}
}

func (e *Item) Reset(version string) {
	if e.Version != version {
		e.Status = StatusInit
	}
}
