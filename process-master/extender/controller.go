package extender

import (
	"context"
)

type Manager struct {
	*Check
	ctx    context.Context
	cancel context.CancelFunc
}

func NewManager(ctx context.Context, callbackFunc ICallback) *Manager {
	c, cancel := context.WithCancel(ctx)
	controller := &Manager{
		ctx:    c,
		cancel: cancel,
		Check:  NewCheck(c, callbackFunc),
	}
	return controller
}

func (e *Manager) Set(key string, ver string) error {
	e.Check.Set(key, ver)
	e.Check.Scan()
	return nil
}

func (e *Manager) Del(key string) error {
	e.Check.Del(key)
	e.Check.Scan()
	return nil
}

func (e *Manager) Reset(data map[string][]byte) error {
	e.Check.Reset(data)
	e.Check.Scan()
	return nil
}
