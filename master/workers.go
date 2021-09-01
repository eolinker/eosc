/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"context"
	"sync"
)

type workerManager struct {
	cmd string
	max int

	once sync.Once

	ctx context.Context
	cancelFunc context.CancelFunc


}

func newWorkerManager(ctx context.Context,cmd string,max int) *workerManager {
	c,cf:=context.WithCancel(ctx)
	return &workerManager{
		cmd: cmd,
		max: max,
		ctx: c,
		cancelFunc: cf,
	}
}

func (m *workerManager) Start() {
	m.once.Do(func() {
		go m.loop()
	})
}
func (m *workerManager) loop() {
	for{

	}
}

func (m *workerManager) forkWorker() {

}