// +build !windows

package syslog

import (
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/drivers"
	"github.com/eolinker/eosc/log/drivers/syslog/config"
)

type Transporter struct {
	*log.Transporter
	writer *_SysWriter
}

func (t *Transporter) Reset(c interface{}, formatter log.Formatter) error {
	conf, err := config.ToConfig(c)
	if err != nil {
		return err
	}
	t.Transporter.SetFormatter(formatter)
	return t.reset(conf)
}
func (t *Transporter) reset(c *config.Config) error {
	t.SetOutput(t.writer)
	t.SetLevel(c.Level)

	return nil
}

var createHandler drivers.CreateHandler = func(c interface{}, f log.Formatter) (reset drivers.TransporterReset, err error) {
	conf, err := config.ToConfig(c)
	if err != nil {
		return nil, err
	}

	sysWriter, err := newSysWriter(conf.Network, conf.RAddr, conf.Level, "")
	if err != nil {
		return nil, err
	}

	transport := &Transporter{
		Transporter: log.NewTransport(sysWriter, conf.Level, f),
		writer:      sysWriter,
	}
	e := transport.Reset(conf, f)
	if e != nil {
		return nil, e
	}
	return transport, nil
	return transport, nil
}

func NewFactory() drivers.TFactory {
	sysConfigDriver := config.NewSysLogConfigDriver()
	return drivers.NewCacheFactory(createHandler, sysConfigDriver)
}
