package mock

import (
	"context"
	"strings"
	"sync"

	"github.com/murtaza-u/ddos/store"
)

type Store struct {
	sync.RWMutex
	kv map[string][]byte
}

func NewStore() store.Storer {
	return &Store{
		kv:      make(map[string][]byte),
		RWMutex: sync.RWMutex{},
	}
}

func (s *Store) Put(k string, v []byte) (int64, error) {
	s.Lock()
	defer s.Unlock()

	s.kv[k] = v
	return 0, nil
}

func (s *Store) Get(k string) ([]byte, int64, error) {
	s.RLock()
	defer s.RUnlock()

	v := s.kv[k]
	if v == nil {
		return nil, 0, store.ErrKeyNotFound
	}

	return s.kv[k], 0, nil
}

func (s *Store) GetKeysWithPrefix(pre string) ([]string, error) {
	var keys []string

	s.RLock()
	defer s.RUnlock()

	for k := range s.kv {
		if strings.HasPrefix(k, pre) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

func (s *Store) Delete(k string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.kv, k)
	return nil
}

func (s *Store) Watcher(ctx context.Context, k string) store.Watcher {
	ctx, cancel := context.WithCancel(ctx)

	return &watcher{
		ctx:    ctx,
		cancel: cancel,
		evC:    make(chan *store.Event),
		errC:   make(chan error),
	}
}
