/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_worker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/listener"

	"github.com/eolinker/eosc/traffic"
)

func Process() {

	worker := NewWorker()
	listener.SetTraffic(worker.tf)
	loadPluginEnv()

	worker.wait()
}

type Worker struct {
	tf traffic.ITraffic
}

func (w *Worker) wait() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		sig := <-sigc
		log.Infof("Caught signal pid:%d ppid:%d signal %s: .\n", os.Getpid(), os.Getppid(), sig.String())
		fmt.Println(os.Interrupt.String(), sig.String(), sig == os.Interrupt)
		switch sig {
		case os.Interrupt, os.Kill:
			{
				w.close()
				return nil
			}
		case syscall.SIGQUIT:
			{
				w.close()
				return nil
			}
		case syscall.SIGUSR1:
			{

			}
		default:
			continue
		}
	}

}
func NewWorker() *Worker {
	w := &Worker{}
	tf := traffic.NewTraffic()
	tf.Read(os.Stdin)
	w.tf = tf
	_, err := utils.ReadFrame(os.Stdin)
	if err != nil {
		return nil
	}
	return w
}

func (w *Worker) close() {

	w.tf.Close()
}
