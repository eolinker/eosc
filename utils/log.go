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
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
	"io"
	"os"
)

func InitLogTransport(name string) io.Writer {
	dir := env.LogDir()
	if env.IsDebug() {
		//dir = filepath.Base(".")
		log.InitDebug(true)
	}
	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}
	writer := filelog.NewFileWriteByPeriod()
	period, _ := filelog.ParsePeriod(env.ErrorPeriod())
	writer.Set(dir, fmt.Sprintf("%s.log", env.ErrorName()), period, env.ErrorExpire())
	writer.Open()
	//transport := filelog.CreateTransporter(log.InfoLevel)

	transport := log.NewTransport(writer, env.ErrorLevel())
	transport.SetFormatter(formatter)
	log.Reset(transport)
	log.SetPrefix(fmt.Sprintf("[%s-%d]", name, os.Getpid()))
	return writer
}

func InitStdTransport(name string, level log.Level) {
	if env.IsDebug() {
		log.InitDebug(true)
	}
	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}
	transport := log.NewTransport(os.Stderr, level)
	transport.SetFormatter(formatter)
	log.Reset(transport)
	log.SetPrefix(fmt.Sprintf("[%s-%d]", name, os.Getpid()))
}
