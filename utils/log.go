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
	"os"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
)

func InitLogTransport(name string) {
	dir := env.LogDir()
	if env.IsDebug() {
		//dir = filepath.Base(".")
		log.InitDebug(true)
	}
	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}
	//writer := filelog.NewFileWriteByPeriod()
	//writer.Set(dir, fmt.Sprintf("%s.log", name), filelog.PeriodDay, 7*24*time.Hour)
	//writer.Open()
	transport := filelog.CreateTransporter(log.InfoLevel)
	transport.Reset(&filelog.Config{
		Dir:    dir,
		File:   fmt.Sprintf("%s.log", name),
		Expire: 7,
		Period: filelog.PeriodDay,
		Level:  log.InfoLevel,
	}, formatter)
	transport.SetFormatter(formatter)
	log.Reset(transport)
	log.SetPrefix(fmt.Sprintf("[%s-%d]", name, os.Getpid()))
}
