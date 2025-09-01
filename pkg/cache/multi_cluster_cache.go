package cache

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ecache"
	eredis "github.com/ecodeclub/ecache/redis"
	"github.com/ecodeclub/ekit"
	"github.com/redis/go-redis/v9"
	"go.uber.org/multierr"
	"sync"
	"sync/atomic"
	"time"
)

type Cluster struct {
	Name     string
	Instance ecache.Cache
}

func NewCluster(name string, rd redis.Cmdable) *Cluster {
	return &Cluster{
		Name: name,
		Instance: &ecache.NamespaceCache{
			C:         eredis.NewCache(rd),
			Namespace: "permission-platform",
		},
	}
}

type MultiClusterCache struct {
	Clusters []*Cluster
}

func NewMultiClusterCache(clusters []*Cluster) *MultiClusterCache {
	return &MultiClusterCache{clusters}
}

func (mc *MultiClusterCache) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	var err error
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := range mc.Clusters {
		instance := mc.Clusters[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			err1 := instance.Instance.Set(ctx, key, val, exp)
			if err1 != nil {
				mu.Lock()
				err = multierr.Append(err, fmt.Errorf("cluster: %s set failed %w", instance.Name, err1))
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return err
}

func (mc *MultiClusterCache) Get(ctx context.Context, key string) ecache.Value {
	var err error
	var mu sync.Mutex
	var needReturn atomic.Bool
	var valPtr atomic.Pointer[ecache.Value]

	queryCluster := func(index int) {
		val := mc.Clusters[index].Instance.Get(ctx, key)
		if val.Err == nil || val.KeyNotFound() {
			needReturn.Store(true)
			valPtr.Store(&val)
		} else {
			mu.Lock()
			err = multierr.Append(err, fmt.Errorf("集群[%s]: %w", mc.Clusters[index].Name, val.Err))
			mu.Unlock()
		}
	}
	const step = 2

	for i, j := 0, 1; i < len(mc.Clusters); i, j = i+step, j+step {
		var wg sync.WaitGroup
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			queryCluster(index)
		}(i)

		if j < len(mc.Clusters) {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				queryCluster(idx)
			}(j)
		}
		wg.Wait()
		if needReturn.Load() {
			return *valPtr.Load()
		}
	}

	return ecache.Value{ekit.AnyValue{
		Err: err,
	}}

}
