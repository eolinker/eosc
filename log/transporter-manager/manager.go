package transporter_manager

import (
	"github.com/eolinker/eosc/log"
)

var (
	defaultNameSpaceManager INameSpaceManager
)

type ResetHandler func(transports ...log.EntryTransporter)

type INameSpaceManager interface {
	RegisterTransporterManager(nameSpace string, rh ResetHandler) error
	GetTransporterManager(nameSpace string) ITransporterManager
}

type ITransporterManager interface {
	Set(workerID string, transporter log.EntryTransporter) error
	Del(workerID string) error
}

func init() {
	defaultNameSpaceManager = NewNameSpaceManager()
	defaultNameSpaceManager.RegisterTransporterManager("default", log.Reset)
}

func GetTransporterManager(nameSpace string) ITransporterManager {
	return defaultNameSpaceManager.GetTransporterManager(nameSpace)
}

func RegisterTransporterManager(nameSpace string, rh ResetHandler) error {
	return defaultNameSpaceManager.RegisterTransporterManager(nameSpace, rh)
}
