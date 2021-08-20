package filelog

import (
	"github.com/eolinker/goku-standard/common/log"
	"github.com/eolinker/goku-standard/common/log/drivers"
	"github.com/eolinker/goku-standard/common/log/drivers/filelog/config"
	"time"
)

type Transporter struct {
	*log.Transporter
	writer *FileWriterByPeriod
}
func (t *Transporter) Close() error {
	t.writer.Close()
	return nil
}
func (t *Transporter) Reset(c interface{},f log.Formatter) error {
	conf, err := config.ToConfig(c)
	if err!=nil{
		return err
	}
	t.Transporter.SetFormatter(f)
	return t.reset(conf)
}

func (t *Transporter) reset(c *config.Config) error {
	t.SetOutput(t.writer)
	t.SetLevel(c.Level)

	t.writer.Set(
		c.Dir,
		c.File,
		c.Period,
		time.Duration(c.Expire)*time.Hour*24,
		)
	t.writer.Open()
	return nil
}
var createHandler drivers.CreateHandler = func(c interface{},formatter log.Formatter) (reset drivers.TransporterReset, err error) {
	conf, err := config.ToConfig(c)
	if err!=nil{
		return  nil,err
	}

	fileWriterByPeriod := NewFileWriteBytePeriod()

	transport := &Transporter{
		Transporter: log.NewTransport(fileWriterByPeriod,conf.Level,formatter),
		writer:       fileWriterByPeriod,
	}

	e:=transport.Reset(conf,formatter)
	if e!= nil{
		return nil,e
	}
	return transport,nil
}
func NewFactory() drivers.TFactory {
	fileConfigDriver := config.NewFileConfigDriver()
	return drivers.NewCacheFactory(createHandler,fileConfigDriver)
}



