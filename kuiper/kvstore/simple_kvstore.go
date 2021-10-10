package kvstore

import (
	"fmt"
	"os"
	"sync"

	"github.com/cloustone/pandas/kuiper/util"

	cache "github.com/patrickmn/go-cache"
)

type SyncKVMap struct {
	sync.RWMutex
	internal map[string]*SimpleKVStore
}

func (sm *SyncKVMap) Load(path string) (result *SimpleKVStore) {
	sm.Lock()
	defer sm.Unlock()
	if s, ok := sm.internal[path]; ok {
		result = s
	} else {
		c := cache.New(cache.NoExpiration, 0)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
		result = NewSimpleKVStore(path+"/stores.data", c)
		sm.internal[path] = result
	}

	return
}

var simpleStore = &SyncKVMap{
	internal: make(map[string]*SimpleKVStore),
}

type CtrlType int

const (
	OPEN CtrlType = iota
	SAVE
	CLOSE
)

type SimpleKVStore struct {
	path string
	c    *cache.Cache
	/* These 2 channels must be mapping one by one*/
	ctrlCh chan CtrlType
	errCh  chan error
}

func NewSimpleKVStore(path string, c *cache.Cache) *SimpleKVStore {
	r := &SimpleKVStore{
		path:   path,
		c:      c,
		ctrlCh: make(chan CtrlType),
		errCh:  make(chan error),
	}
	go r.run()
	return r
}

func (m *SimpleKVStore) run() {
	count := 0
	opened := false
	for c := range m.ctrlCh {
		switch c {
		case OPEN:
			count++
			if !opened {
				if _, err := os.Stat(m.path); os.IsNotExist(err) {
					m.errCh <- nil
					break
				}
				if e := m.c.LoadFile(m.path); e != nil {
					m.errCh <- e
					break
				}
			}
			m.errCh <- nil
			opened = true
		case CLOSE:
			count--
			if count == 0 {
				opened = false
				err := m.doClose()
				if err != nil {
					// Log.Error(err) TODO
					m.errCh <- err
					break
				}
			}
			m.errCh <- nil
		case SAVE:
			//swallow duplicate requests
			if len(m.ctrlCh) > 0 {
				m.errCh <- nil
				break
			}
			if e := m.c.SaveFile(m.path); e != nil {
				// Log.Error(e) TODO
				m.errCh <- e
				break
			}
			m.errCh <- nil
		}
	}
}

func (m *SimpleKVStore) Open() error {
	m.ctrlCh <- OPEN
	return <-m.errCh
}

func (m *SimpleKVStore) Close() error {
	m.ctrlCh <- CLOSE
	return <-m.errCh
}

func (m *SimpleKVStore) doClose() error {
	//e := m.c.SaveFile(m.path)
	m.c.Flush() //Delete all of the values from memory.
	return nil
}

func (m *SimpleKVStore) saveToFile() error {
	m.ctrlCh <- SAVE
	return <-m.errCh
}

func (m *SimpleKVStore) Set(key string, value interface{}) error {
	if m.c == nil {
		return fmt.Errorf("cache %s has not been initialized yet.", m.path)
	}
	if err := m.c.Add(key, value, cache.NoExpiration); err != nil {
		return err
	}
	return m.saveToFile()
}

func (m *SimpleKVStore) Replace(key string, value interface{}) error {
	if m.c == nil {
		return fmt.Errorf("cache %s has not been initialized yet.", m.path)
	}
	m.c.Set(key, value, cache.NoExpiration)
	return m.saveToFile()
}

func (m *SimpleKVStore) Get(key string) (interface{}, bool) {
	return m.c.Get(key)
}

func (m *SimpleKVStore) Delete(key string) error {
	if m.c == nil {
		return fmt.Errorf("cache %s has not been initialized yet.", m.path)
	}
	if _, found := m.c.Get(key); found {
		m.c.Delete(key)
	} else {
		return util.NewErrorWithCode(util.NOT_FOUND, fmt.Sprintf("%s is not found", key))
	}
	return m.saveToFile()
}

func (m *SimpleKVStore) Keys() (keys []string, err error) {
	if m.c == nil {
		return nil, fmt.Errorf("cache %s has not been initialized yet.", m.path)
	}
	its := m.c.Items()
	keys = make([]string, 0, len(its))
	for k := range its {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *SimpleKVStore) Clean() error {
	return os.RemoveAll(m.path)
}
