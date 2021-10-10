package states

import (
	"sync"
)

type MemoryStore sync.Map //The current root store of a rule

func newMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) SaveState(_ int64, _ string, _ map[string]interface{}) error {
	//do nothing
	return nil
}

func (s *MemoryStore) SaveCheckpoint(_ int64) error {
	//do nothing
	return nil
}

func (s *MemoryStore) GetOpState(opId string) (*sync.Map, error) {
	return &sync.Map{}, nil
}
