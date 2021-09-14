/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package eosc

import "fmt"

func (p *Profession) checkConfig(driver string, cdata IData, workers IWorkers) (IProfessionDriver, interface{}, map[RequireId]interface{}, error) {
	d, has := p.getDriver(driver)
	if !has {
		return nil, nil, nil, fmt.Errorf("%s:%w", driver, ErrorDriverNotExist)
	}

	config := newConfig(d.ConfigType())

	err := cdata.UnMarshal(&config)
	if err != nil {
		return nil, nil, nil, err
	}

	requires, err := CheckConfig(config, workers)
	if err != nil {
		return nil, nil, nil, err
	}
	if dc, ok := d.(IProfessionDriverCheckConfig); ok {
		if e := dc.Check(config, requires); e != nil {
			return nil, nil, nil, e
		}
	}
	return d, config, requires, nil
}
func (p *Profession) CheckerConfig(driver string, cdata IData, workers IWorkers) error {

	_, _, _, err := p.checkConfig(driver, cdata, workers)
	return err
}

func (p *Profession) ChangeWorker(driver, id, name string, cdata IData, workers IWorkers) error {
	d, cf, requires, err := p.checkConfig(driver, cdata, workers)
	if err != nil {
		return err
	}

	if w, has := workers.Get(id); has {
		err := w.Reset(cf, requires)
		if err != nil {
			return err
		}

	} else {
		w, err := d.Create(id, name, cf, requires)
		if err != nil {
			return err
		}
		err = w.Start()
		if err != nil {
			return err
		}
		workers.Set(id, w)

	}
	p.setId(name, id)
	return nil
}
