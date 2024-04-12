package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils/set"
	"reflect"
	"strings"
)

func (oe *imlAdminApi) SetSetting(ctx context.Context, name string, data marshal.IData) error {
	driver, has := oe.settings.GetDriver(name)
	if !has {
		return ErrorNotExist
	}
	if driver.Mode() == eosc.SettingModeReadonly {
		return ErrorReadOnly
	}

	inputData, err := data.Encode()
	if err != nil {

		return err
	}
	configType := driver.ConfigType()
	if driver.Mode() == eosc.SettingModeSingleton {
		oldConfig, hasOld := oe.settings.GetConfigBody(name)
		if !hasOld {
			ocObj := oe.settings.GetConfig(name)
			if ocObj != nil {
				od, err := json.Marshal(ocObj)
				if err != nil {
					oldConfig = od
					hasOld = true
				}
			}
		}
		err := oe.settings.SettingWorker(name, inputData, oe.variable)
		if err != nil {
			return err
		}
		wc := &eosc.WorkerConfig{
			Id:          fmt.Sprintf("%s@setting", name),
			Profession:  Setting,
			Name:        name,
			Driver:      name,
			Create:      eosc.Now(),
			Update:      eosc.Now(),
			Body:        inputData,
			Description: "",
		}
		eventData, _ := json.Marshal(wc)
		if hasOld {
			oe.actions = append(oe.actions, newRollbackForSettingSet(name, oldConfig))
		} else {
			oe.actions = append(oe.actions, newRollbackForSettingSet(name, nil))
		}
		oe.events = append(oe.events, &service.Event{
			Command:   eosc.EventSet,
			Namespace: eosc.NamespaceWorker,
			Key:       wc.Id,
			Data:      eventData,
		})

	} else {
		err = oe.batchSetWorker(ctx, inputData, driver, configType)
		if err != nil {
			log.Debug("batch set:", name, ":", string(inputData))
			log.Info("batch set:", name, ":", err)
			return err
		}

	}
	return nil
}
func (oe *imlAdminApi) batchSetWorker(ctx context.Context, inputData []byte, driver eosc.ISetting, configType reflect.Type) error {

	inputList := marshal.SplitConfig(inputData)
	cfgs := make(map[string]*eosc.WorkerConfig, len(inputList))
	allWorkers := set.NewSet(driver.AllWorkers()...)

	version := GenVersion()
	for _, inp := range inputList {
		configData, _ := inp.Encode()
		cfg, _, err2 := oe.variable.Unmarshal(configData, configType)
		if err2 != nil {

			return err2
		}
		profession, workerName, driverName, desc, errCk := driver.Check(cfg)
		if errCk != nil {

			return errCk
		}
		id, _ := eosc.ToWorkerId(workerName, profession)
		if allWorkers.Contains(id) {
			allWorkers.Remove(id)
		}
		cfgs[id] = &eosc.WorkerConfig{
			Id:          id,
			Profession:  profession,
			Name:        workerName,
			Driver:      driverName,
			Description: desc,
			Body:        configData,
			Version:     version,
			Update:      eosc.Now(),
		}
	}
	idToDelete := allWorkers.List()

	cannotDelete := oe.CheckDelete(idToDelete...)
	if len(cannotDelete) > 0 {
		return fmt.Errorf("should not delete:%s", strings.Join(cannotDelete, ","))
	}

	for id, cfg := range cfgs {
		_, errSet := oe.SetWorker(ctx, cfg)
		if errSet != nil {
			log.Warnf("bath set  %s fail :%v", id, errSet)
			return fmt.Errorf("bath set  %s fail :%v", id, errSet)
		}

	}

	for _, id := range idToDelete {
		_, err := oe.DeleteWorker(ctx, id)
		if err != nil {
			return fmt.Errorf("delete worker %s %w", id, err)
		}

	}
	return nil

}
