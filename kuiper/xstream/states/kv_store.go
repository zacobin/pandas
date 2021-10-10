package states

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path"
	"sync"

	"github.com/cloustone/pandas/kuiper/kvstore"
	"github.com/cloustone/pandas/kuiper/util"
	"github.com/cloustone/pandas/kuiper/xstream/checkpoints"
)

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(checkpoints.BufferOrEvent{})
}

//The manager for checkpoint storage.
/**
*** mapStore keys
***   { "checkpoint1", "checkpoint2" ... "checkpointn" : The complete or incomplete snapshot
 */
type KVStore struct {
	db          kvstore.KvStore
	mapStore    *sync.Map //The current root store of a rule
	checkpoints []int64
	max         int
}

//Store in path ./data/checkpoint/$ruleId
//Store 2 things:
//"checkpoints":A queue for completed checkpoint id
//"$checkpointId":A map with key of checkpoint id and value of snapshot(gob serialized)
//Assume each operator only has one instance
func getKVStore(ruleId string) (*KVStore, error) {
	dr, _ := util.GetDataLoc()
	db := kvstore.GetKvStore(path.Join(dr, "checkpoints", ruleId))
	s := &KVStore{db: db, max: 3, mapStore: &sync.Map{}}
	//read data from badger db
	if err := s.restore(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *KVStore) restore() error {
	err := s.db.Open()
	if err != nil {
		return err
	}
	defer s.db.Close()
	if b, ok := s.db.Get(CheckpointListKey); ok {
		if cs, err := bytesToSlice(b.([]byte)); err != nil {
			return fmt.Errorf("invalid checkpoint data: %s", err)
		} else {
			s.checkpoints = cs
			for _, c := range cs {
				if b2, ok := s.db.Get(fmt.Sprintf("%d", c)); ok {
					if m, err := bytesToMap(b2.([]byte)); err != nil {
						return fmt.Errorf("invalid checkpoint data: %s", err)
					} else {
						s.mapStore.Store(c, util.MapToSyncMap(m))
					}
				}
			}
		}
	}
	return nil
}

func (s *KVStore) SaveState(checkpointId int64, opId string, state map[string]interface{}) error {
	logger := util.Log
	logger.Debugf("Save state for checkpoint %d, op %s, value %v", checkpointId, opId, state)
	var cstore *sync.Map
	if v, ok := s.mapStore.Load(checkpointId); !ok {
		cstore = &sync.Map{}
		s.mapStore.Store(checkpointId, cstore)
	} else {
		if cstore, ok = v.(*sync.Map); !ok {
			return fmt.Errorf("invalid KVStore for checkpointId %d with value %v: should be *sync.Map type", checkpointId, v)
		}
	}
	cstore.Store(opId, state)
	return nil
}

func (s *KVStore) SaveCheckpoint(checkpointId int64) error {
	if v, ok := s.mapStore.Load(checkpointId); !ok {
		return fmt.Errorf("store for checkpoint %d not found", checkpointId)
	} else {
		if m, ok := v.(*sync.Map); !ok {
			return fmt.Errorf("invalid KVStore for checkpointId %d with value %v: should be *sync.Map type", checkpointId, v)
		} else {
			err := s.db.Open()
			if err != nil {
				return fmt.Errorf("save checkpoint err: %v", err)
			}
			defer s.db.Close()
			b, err := mapToBytes(m)
			if err != nil {
				return fmt.Errorf("save checkpoint err, fail to encode states: %s", err)
			}
			err = s.db.Replace(fmt.Sprintf("%d", checkpointId), b)
			if err != nil {
				return fmt.Errorf("save checkpoint err: %v", err)
			}
			m.Delete(checkpointId)
			s.checkpoints = append(s.checkpoints, checkpointId)
			//TODO is the order promised?
			if len(s.checkpoints) > s.max {
				cp := s.checkpoints[0]
				s.checkpoints = s.checkpoints[1:]
				go func() {
					_ = s.db.Delete(fmt.Sprintf("%d", cp))
				}()
			}
			cs, ok := sliceToBytes(s.checkpoints)
			if !ok {
				return fmt.Errorf("save checkpoint err: fail to encode checkpoint counts")
			}
			err = s.db.Replace(CheckpointListKey, cs)
			if err != nil {
				return fmt.Errorf("save checkpoint err: %v", err)
			}
		}
	}
	return nil
}

//Only run in the initialization
func (s *KVStore) GetOpState(opId string) (*sync.Map, error) {
	if len(s.checkpoints) > 0 {
		if v, ok := s.mapStore.Load(s.checkpoints[len(s.checkpoints)-1]); ok {
			if cstore, ok := v.(*sync.Map); !ok {
				return nil, fmt.Errorf("invalid state %v stored for op %s: data type is not *sync.Map", v, opId)
			} else {
				if sm, ok := cstore.Load(opId); ok {
					switch m := sm.(type) {
					case *sync.Map:
						return m, nil
					case map[string]interface{}:
						return util.MapToSyncMap(m), nil
					default:
						return nil, fmt.Errorf("invalid state %v stored for op %s: data type is not *sync.Map", sm, opId)
					}
				}
			}
		} else {
			return nil, fmt.Errorf("store for checkpoint %d not found", s.checkpoints[len(s.checkpoints)-1])
		}
	}
	return &sync.Map{}, nil
}

func mapToBytes(sm *sync.Map) ([]byte, error) {
	m := util.SyncMapToMap(sm)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bytesToMap(input []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	buf := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func sliceToBytes(s []int64) ([]byte, bool) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		return nil, false
	}
	return buf.Bytes(), true
}

func bytesToSlice(input []byte) ([]int64, error) {
	result := make([]int64, 3)
	buf := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
