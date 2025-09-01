package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ecodeclub/ecache"
	"github.com/ego-component/eetcd"
	"github.com/gotomicro/ego/core/econf"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
	"github.com/permission-dev/pkg/bitring"
	cache2 "github.com/permission-dev/pkg/cache"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync/atomic"
	"time"
)

func InitCacheKeyFunc() func(bizID, userID int64) string {
	return func(bizID, userID int64) string {
		return fmt.Sprintf("permission:bizID:%d:userID:%d", bizID, userID)
	}
}

func InitMultiLevelCache(
	r redis.Cmdable,
	local ecache.Cache,
	repo repository.UserPermissionRepository,
	etcdClient *eetcd.Component,
	cacheKeyFunc func(bizID, userID int64) string,
) cache2.Cache {
	type ErrorEventConfig struct {
		BitRingSize      int     `yaml:"bitRingSize"`
		RateThreshold    float64 `yaml:"rateThreshold"`
		ConsecutiveCount int     `yaml:"consecutiveCount"`
	}
	type Config struct {
		EtcdKey                 string           `yaml:"etcdKey"`
		LocalCacheRefreshPeriod time.Duration    `yaml:"localCacheRefreshPeriod"`
		RedisPingTimeout        time.Duration    `yaml:"redisPingTimeout"`
		RedisHealthCheckPeriod  time.Duration    `yaml:"redisHealthCheckPeriod"`
		ErrorEvents             ErrorEventConfig `yaml:"errorEvents"`
	}
	var cfg Config
	err := econf.UnmarshalKey("cache.multilevel", &cfg)
	if err != nil {
		panic(err)
	}

	hotUserLoader := NewHotUserLoader([]domain.User{}, repo, cacheKeyFunc)
	go func() {
		watchChan := etcdClient.Watch(context.Background(), cfg.EtcdKey)
		for watchChanResp := range watchChan {
			for _, event := range watchChanResp.Events {
				if event.Type == clientv3.EventTypePut {
					_ = hotUserLoader.UpdateUsers(event.Kv.Value)
				}
			}
		}
	}()
	return cache2.NewMultiCacheV2(
		r,
		local,
		cfg.LocalCacheRefreshPeriod,
		cfg.RedisPingTimeout,
		cfg.RedisHealthCheckPeriod,
		bitring.NewBitRing(
			cfg.ErrorEvents.BitRingSize,
			cfg.ErrorEvents.ConsecutiveCount,
			cfg.ErrorEvents.RateThreshold,
		),
		hotUserLoader.LoadUserPermissionFromDB,
	)
}

type HotUserLoader struct {
	hotUsersPtr  atomic.Pointer[[]domain.User]
	repo         repository.UserPermissionRepository
	cacheKeyFunc func(bizID, userID int64) string
}

func NewHotUserLoader(hotUsers []domain.User, repo repository.UserPermissionRepository, cacheKeyFunc func(bizID, userID int64) string) *HotUserLoader {
	h := &HotUserLoader{
		repo:         repo,
		cacheKeyFunc: cacheKeyFunc,
	}
	h.hotUsersPtr.Store(&hotUsers)
	return h
}

func (h *HotUserLoader) LoadUserPermissionFromDB(ctx context.Context) ([]*cache2.Entry, error) {
	const day = 24 * time.Hour
	const defaultExpiration = 36500 * day
	hotUsers := *h.hotUsersPtr.Load()
	entries := make([]*cache2.Entry, 0, len(hotUsers))
	for index := range hotUsers {
		perms, err := h.repo.FindByBizIdAndUserID(ctx, hotUsers[index].BizID, hotUsers[index].ID)
		if err == nil {
			val, _ := json.Marshal(perms)
			entries = append(entries, &cache2.Entry{
				Key:        h.cacheKeyFunc(hotUsers[index].BizID, hotUsers[index].ID),
				Val:        val,
				Expiration: defaultExpiration,
			})
		}
	}
	return entries, nil
}

func (h *HotUserLoader) UpdateUsers(value []byte) error {
	var users []domain.User

	err := json.Unmarshal(value, &users)
	if err == nil {
		h.hotUsersPtr.Store(&users)
	}
	return err
}
