package inmemory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/k3s-io/kine/pkg/client"
	"github.com/vishal-chdhry/policy-reports-extension-api/server/pkg/common"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type inMemoryDb struct {
	sync.Mutex
	db map[string]client.Value
}

func New() client.Client {
	return &inMemoryDb{}
}

func (i *inMemoryDb) List(ctx context.Context, prefix string, rev int) ([]client.Value, error) {
	i.Lock()
	defer i.Unlock()

	res := make([]client.Value, 0)

	for k, v := range i.db {
		if strings.HasPrefix(k, prefix) {
			res = append(res, v)
		}
	}

	return res, nil
}

func (i *inMemoryDb) Get(ctx context.Context, key string) (client.Value, error) {
	i.Lock()
	defer i.Unlock()

	if val, ok := i.db[key]; ok {
		return val, nil
	} else {
		return client.Value{}, errors.NewNotFound(schema.GroupResource{Group: common.GroupVersion, Resource: ""}, key)
	}
}

func (i *inMemoryDb) Put(ctx context.Context, key string, value []byte) error {
	i.Lock()
	defer i.Unlock()

	i.db[key] = client.Value{
		Key:      []byte(key),
		Data:     value,
		Modified: time.Now().Unix(),
	}

	return nil
}

func (i *inMemoryDb) Create(ctx context.Context, key string, value []byte) error {
	i.Lock()
	defer i.Unlock()

	if _, found := i.db[key]; found {
		return errors.NewAlreadyExists(schema.GroupResource{Group: common.GroupVersion, Resource: ""}, key)
	} else {
		i.db[key] = client.Value{
			Key:      []byte(key),
			Data:     value,
			Modified: time.Now().Unix(),
		}
		return nil
	}
}

func (i *inMemoryDb) Update(ctx context.Context, key string, revision int64, value []byte) error {
	i.Lock()
	defer i.Unlock()

	if _, found := i.db[key]; !found {
		return errors.NewNotFound(schema.GroupResource{Group: common.GroupVersion, Resource: ""}, key)
	} else {
		i.db[key] = client.Value{
			Key:      []byte(key),
			Data:     value,
			Modified: time.Now().Unix(),
		}
		return nil
	}
}

func (i *inMemoryDb) Delete(ctx context.Context, key string, revision int64) error {
	i.Lock()
	defer i.Unlock()

	if _, found := i.db[key]; !found {
		return errors.NewNotFound(schema.GroupResource{Group: common.GroupVersion, Resource: ""}, key)
	} else {
		delete(i.db, key)
		return nil
	}
}

func (i *inMemoryDb) Close() error {
	i.db = nil
	return nil
}
