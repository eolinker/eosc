/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package master

import (
	"fmt"
	"time"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
	"github.com/eolinker/eosc/process"
)

func (m *Master) InitLogTransport() {
	writer := filelog.NewFileWriteByPeriod()
	writer.Set(fmt.Sprintf("/var/log/%s", process.AppName()), "error.log", filelog.PeriodDay, 7*24*time.Hour)
	writer.Open()
	transport := log.NewTransport(writer, log.InfoLevel)
	formater := &log.LineFormatter{
		TimestampFormat:  "[2006-01-02 15:04:05]",
		CallerPrettyfier: nil,
	}
	transport.SetFormatter(formater)
	log.NewStdTransport(formater)
	log.Reset(transport, log.NewStdTransport(formater))
}
