/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
)

func InitLogTransport(name string) {
	dir := fmt.Sprintf("/var/log/%s", env.AppName())
	if env.IsDebug() {
		dir = filepath.Base(".")
		log.InitDebug(true)
	}
	writer := filelog.NewFileWriteByPeriod()
	writer.Set(dir, fmt.Sprintf("%s.log", name), filelog.PeriodDay, 7*24*time.Hour)
	writer.Open()
	transport := log.NewTransport(writer, log.InfoLevel)
	formater := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}
	transport.SetFormatter(formater)
	log.Reset(transport)
}
