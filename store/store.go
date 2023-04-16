package store

import (
	"context"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var ErrKeyNotFound = errors.New("key not found")

// Timeout defines the default database dial timeout.
const Timeout = time.Second * 30

// Storer defines behaviour all storage implementations must have.
type Storer interface {
	Put(string, []byte) (int64, error)
	Get(string) ([]byte, int64, error)
	GetKeysWithPrefix(string) ([]string, error)
	Delete(string) error
	Watcher(context.Context, string) Watcher
}

type store struct {
	client *clientv3.Client
}

// New returns a new storage implememtation.
func New(ends []string) (Storer, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   ends,
		DialTimeout: Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &store{c}, nil
}

// Put creates/updates a key-value pair.
func (s store) Put(k string, v []byte) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	res, err := s.client.Put(ctx, k, string(v), clientv3.WithPrevKV())
	if err != nil {
		return 0, err
	}

	if res.PrevKv == nil {
		return 1, nil
	}

	return res.PrevKv.Version + 1, nil
}

// Get returns the value associated with the key. Returns an error if
// the key does not exist.
func (s store) Get(k string) ([]byte, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	res, err := s.client.Get(ctx, k)
	if err != nil {
		return nil, 0, err
	}

	if len(res.Kvs) == 0 {
		return nil, 0, ErrKeyNotFound
	}

	kv := res.Kvs[0]

	return kv.Value, kv.Version, nil
}

func (s store) Delete(k string) error {
	return s.Delete(k)
}

// Watcher returns a new watcher instance that can be started to watch
// over the given key.
func (s store) Watcher(ctx context.Context, k string) Watcher {
	ctx, cancel := context.WithCancel(ctx)
	wc := s.client.Watch(ctx, k)
	return &watch{
		wc:     wc,
		ctx:    ctx,
		cancel: cancel,
		evC:    make(chan *Event, 100),
		errC:   make(chan error),
	}
}

// GetKeysWithPrefix returns keys starting with the given prefix.
func (s store) GetKeysWithPrefix(pre string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	res, err := s.client.Get(
		ctx, pre,
		clientv3.WithPrefix(), clientv3.WithKeysOnly(),
	)
	if err != nil {
		return nil, err
	}

	if len(res.Kvs) == 0 {
		return nil, ErrKeyNotFound
	}

	var keys []string
	for _, v := range res.Kvs {
		keys = append(keys, string(v.Key))
	}

	return keys, nil
}
