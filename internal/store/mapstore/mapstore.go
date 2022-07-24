package mapstore

import (
	"fmt"
	"sync"

	"github.com/ineverbee/wbl0/internal/store"
)

type MapStore struct {
	sync.RWMutex
	m map[int]*store.Model
}

func NewMapStore(mp map[int]*store.Model) *MapStore {
	return &MapStore{m: mp}
}

func (ms *MapStore) Get(id int) (*store.Model, error) {
	defer ms.RUnlock()
	ms.RLock()
	if model, ok := ms.m[id]; ok {
		return model, nil
	}
	return nil, fmt.Errorf("error: no rows with id '%d'", id)
}

func (ms *MapStore) Set(id *int, model *store.Model) error {
	defer ms.Unlock()
	ms.Lock()
	ms.m[*id] = model
	return nil
}
