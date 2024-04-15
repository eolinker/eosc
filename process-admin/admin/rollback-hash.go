package admin

import (
	"github.com/eolinker/eosc/utils/hash"
)

type RollBackHashForAdd string

func (r RollBackHashForAdd) RollBack(api iAdminOperator) error {
	api.delHashValue(string(r))
	return nil
}
func NewRollbackForAddHash(key string) RollbackHandler {
	return RollBackHashForAdd(key)
}

type RollBackHashForSet struct {
	key  string
	data hash.Hash[string, string]
}

func NewRollBackForSetHash(key string, data hash.Hash[string, string]) RollbackHandler {
	return &RollBackHashForSet{key: key, data: data}
}
func NewRollBackForDeleteHash(key string, data hash.Hash[string, string]) RollbackHandler {
	return &RollBackHashForSet{key: key, data: data}
}
func (r *RollBackHashForSet) RollBack(api iAdminOperator) error {
	api.setHashValue(r.key, r.data)
	return nil
}
