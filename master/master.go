/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"github.com/eolinker/eosc/log"
	"os"
	"os/signal"
	"syscall"
)
func Process() {
	master:=NewMasterHandle()
	master.Start()
	master.Wait()

}

type Master struct {

}

func (m *Master) Start() {

}

func (m *Master) Wait() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Wait for a SIGINT or SIGKILL:
	sig := <- sigc
	log.Info("Caught signal %s: shutting down.", sig)

	m.close()
	return nil
}
func (m *Master) close()  {


}

func NewMasterHandle() *Master {
	return &Master{}
}
