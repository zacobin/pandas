package kuiper

import (
	"sync"

	"github.com/cloustone/pandas/kuiper/xstream"
)

type ruleState struct {
	name      string
	topology  *xstream.TopologyNew
	triggered bool
}
type ruleRegistry struct {
	sync.RWMutex
	internal map[string]*ruleState
}

func (rr *ruleRegistry) store(key string, value *ruleState) {
	rr.Lock()
	rr.internal[key] = value
	rr.Unlock()
}

func (rr *ruleRegistry) load(key string) (value *ruleState, ok bool) {
	rr.RLock()
	result, ok := rr.internal[key]
	rr.RUnlock()
	return result, ok
}

func (rr *ruleRegistry) delete(key string) {
	rr.Lock()
	delete(rr.internal, key)
	rr.Unlock()
}
