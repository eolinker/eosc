package drivers

import (
	"errors"
	"sync"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/dlog"
)

var (
	ErrorNeedConfigDriver  = errors.New("need ConfigDriver")
	ErrorNeedCreateHandler = errors.New("need ErrorNeedCreateHandler")
	ErrorNotSuperReset     = errors.New("not supper reset")
)

type CreateHandler func(config interface{}, formatter log.Formatter) (TransporterReset, error)
type TransporterReset interface {
	log.EntryTransporter
	Reset(config interface{}, formatter log.Formatter) error
}
type CacheFactory struct {
	locker        sync.Locker
	createHandler CreateHandler
	configDriver  dlog.ConfigDriver

	cache map[string]TransporterReset

	driver string
}

func (c *CacheFactory) Destroy(name string) {

	c.locker.Lock()
	defer c.locker.Unlock()

	transporter, has := c.cache[name]
	if has {
		transporter.Close()
		delete(c.cache, name)
	}
}

func (c *CacheFactory) Driver() string {
	return c.driver
}

func NewCacheFactory(createHandler CreateHandler, dLogConfigDriver dlog.ConfigDriver) *CacheFactory {
	if dLogConfigDriver == nil {
		panic(ErrorNeedConfigDriver)
	}

	if createHandler == nil {
		panic(ErrorNeedCreateHandler)
	}

	return &CacheFactory{
		driver:        dLogConfigDriver.Name(),
		locker:        &sync.Mutex{},
		createHandler: createHandler,
		configDriver:  dLogConfigDriver,
		cache:         make(map[string]TransporterReset),
	}
}

func (c *CacheFactory) Get(name string, config string, f log.Formatter) (log.EntryTransporter, error) {

	conf, err := c.configDriver.Decode(config)
	if err != nil {
		log.Errorf("decode config name:[%s]  driver:[%s] error:%s", name, c.configDriver.Name(), err)
		return nil, err
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	transporter, has := c.cache[name]
	if has {

		err := transporter.Reset(conf, f)
		if err != nil {
			delete(c.cache, name)
			if err == ErrorNotSuperReset || errors.Unwrap(err) == ErrorNotSuperReset {
				// 不支持reset的的接口返回 ErrorNotSuperReset 会导致新建，并close旧的
				e := transporter.Close()
				log.Errorf("close transporter  err:%s", e)
				transporter = nil

			} else {
				log.Errorf("reset transporter  err:%s", err)
				//其他错误终止
				return nil, err
			}
		} else {
			// reset成功
			return transporter, nil
		}

	}
	// 创建新transporter
	if transporter == nil {
		newTransporter, err := c.createHandler(conf, f)
		if err != nil {
			return nil, err
		}
		c.cache[name] = newTransporter
		transporter = newTransporter
	}
	return transporter, nil
}
