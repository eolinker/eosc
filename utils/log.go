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
	"io"
	"net/http"
	"os"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
)

func InitMasterLog() (io.Writer, func(prefix string) http.Handler) {
	dir := env.LogDir()

	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}
	fileWriter := filelog.NewFileWriteByPeriod(filelog.Config{
		Dir:    dir,
		File:   fmt.Sprintf("%s.log", env.ErrorName()),
		Expire: env.ErrorExpire(),
		Period: filelog.ParsePeriod(env.ErrorPeriod()),
	})

	level := env.ErrorLevel()
	writer := ToCopyToIoWriter(os.Stdout, fileWriter)

	transport := log.NewTransport(writer, level)
	transport.SetFormatter(formatter)
	transport.Output()
	log.Reset(transport)
	log.SetPrefix(fmt.Sprintf("[%s-%d]", eosc.ProcessMaster, os.Getpid()))

	return writer, fileWriter.ServeHTTP
}

type writes []io.Writer

func ToCopyToIoWriter(ws ...io.Writer) io.Writer {
	return writes(ws)
}
func (ws writes) Write(p []byte) (n int, err error) {
	for _, w := range ws {
		n, err = w.Write(p)
	}
	return
}

func InitStdTransport(name ...string) {
	level := env.ErrorLevel()
	if env.IsDebug() {
		level = log.DebugLevel
	}
	transport := log.NewTransport(os.Stderr, level)
	transport.SetFormatter(&log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	})
	log.Reset(transport)
	if len(name) > 0 {
		log.SetPrefix(fmt.Sprintf("[%s-%d]", name[0], os.Getpid()))
	}
}
