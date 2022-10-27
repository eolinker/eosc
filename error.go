package eosc

import "errors"

var (
	ErrorDriverNotExist             = errors.New("driver not exist")
	ErrorProfessionNotExist         = errors.New("profession not exist")
	ErrorNotAllowCreateForSingleton = errors.New("not allow create for singleton profession")
	ErrorWorkerNotExits             = errors.New("worker-data not exits")
	ErrorWorkerNotRunning           = errors.New("worker-data not running")
	ErrorRegisterConflict           = errors.New("conflict of register")
	ErrorNotGetSillForRequire       = errors.New("not get skill for require")
	ErrorTargetNotImplementSkill    = errors.New("require of skill not implement")
	ErrorParamsIsNil                = errors.New("params is nil")
	ErrorParamNotExist              = errors.New("not exist")
	ErrorStoreReadOnly              = errors.New("store read only")
	ErrorRequire                    = errors.New("require")
	ErrorProfessionDependencies     = errors.New("profession dependencies not complete")
	ErrorConfigIsNil                = errors.New("config is nil")
	ErrorConfigFieldUnknown         = errors.New("unknown type")
	ErrorConfigType                 = errors.New("error config type")
)
