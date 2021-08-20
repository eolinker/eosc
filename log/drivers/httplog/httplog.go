package httplog

import (
	"github.com/eolinker/goku-standard/common/log"
	"github.com/eolinker/goku-standard/common/log/drivers"
	"github.com/eolinker/goku-standard/common/log/drivers/httplog/config"
)

type Transporter struct {
	*log.Transporter
	writer *_HttpWriter
}


func (t *Transporter) Reset(c interface{},formatter log.Formatter) error {
	conf, err := config.ToConfig(c)
	if err!=nil{
		return err
	}
	t.Transporter.SetFormatter(formatter)
	return t.reset(conf)
}
func (t *Transporter) Close() error {
	t.Transporter.Close()
	return t.writer.Close()
}
func (t *Transporter) reset(c *config.Config) error {
	t.SetOutput(t.writer)
	t.SetLevel(c.Level)

	t.writer.reset(c)
	t.Transporter.SetOutput(t.writer)
	return nil
}
var createHandler drivers.CreateHandler = func(c interface{},f log.Formatter) (reset drivers.TransporterReset, err error) {
	conf, err := config.ToConfig(c)
	if err!=nil{
		return  nil,err
	}

	httpWriter := newHttpWriter()

	transport := &Transporter{
		Transporter: log.NewTransport(httpWriter,conf.Level,f),
		writer:       httpWriter,
	}
	e:=transport.Reset(conf,f)
	if e!= nil{
		return nil,e
	}
	return transport,nil
}
func NewFactory() drivers.TFactory {
	fileConfigDriver := config.NewHttpLogConfigDriver()
	return drivers.NewCacheFactory(createHandler,fileConfigDriver)
}


