package cache

import (
	"context"
	"github.com/ecodeclub/ecache"
	eredis "github.com/ecodeclub/ecache/redis"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/pkg/bitring"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"time"
)

type Entry struct {
	Key        string
	Val        any
	Expiration time.Duration
}

type DataLoader func(ctx context.Context) ([]*Entry, error)

type MultiCacheV2 struct {
	redis                ecache.Cache
	local                ecache.Cache
	isRedisAvailable     atomic.Bool
	redisPingTimeOut     time.Duration
	redisHealCheckPeriod time.Duration
	redisCrashDetector   *bitring.BitRing
	mu                   sync.Mutex
	dataLoader           DataLoader

	//定时刷新
	refreshTicker         *time.Ticker
	stopCtx               context.Context
	stopRefreshCancelFunc context.CancelFunc
	localRefreshPeriod    time.Duration
	logger                *elog.Component
}

func NewMultiCacheV2(
	rd redis.Cmdable,
	local ecache.Cache,
	localRefreshPeriod,
	redisPingTimeOut,
	redisHealCheckPeriod time.Duration,
	redisCrashDetector *bitring.BitRing,
	loader DataLoader,
) *MultiCacheV2 {
	mlc := &MultiCacheV2{
		redis: &ecache.NamespaceCache{
			C:         eredis.NewCache(rd),
			Namespace: "permission-platform",
		},
		local:                local,
		redisPingTimeOut:     redisPingTimeOut,
		redisHealCheckPeriod: redisHealCheckPeriod,
		redisCrashDetector:   redisCrashDetector,
		localRefreshPeriod:   localRefreshPeriod,
		dataLoader:           loader,
		logger:               elog.EgoLogger,
	}
	mlc.stopCtx, mlc.stopRefreshCancelFunc = context.WithCancel(context.Background())
	mlc.isRedisAvailable.Store(true)
	go mlc.redisHealthCheck(rd)
	return mlc
}
func (mlc *MultiCacheV2) redisHealthCheck(rd redis.Cmdable) {
	ticker := time.NewTicker(mlc.redisHealCheckPeriod)
	defer ticker.Stop()
	for range ticker.C {
		if !mlc.isRedisAvailable.Load() {
			ctx, cancel := context.WithCancel(context.Background())
			if err := rd.Ping(ctx); err == nil {
				mlc.handleRedisRecoveryEvent(context.Background())
			}
			cancel()
		}
	}
}
func (mlc *MultiCacheV2) handleRedisRecoveryEvent(ctx context.Context) {
	mlc.mu.Lock()
	defer mlc.mu.Unlock()
	if mlc.isRedisAvailable.Load() {
		return
	}
	mlc.isRedisAvailable.Store(true)
	mlc.stopRefreshCancelFunc()
	mlc.redisCrashDetector.Reset()
	if err := mlc.loadFromDBToCache(ctx, mlc.redis); err != nil {
		mlc.logger.Error("从数据库加载数据到Redis失败", elog.FieldErr(err))
	}
}
func (mlc *MultiCacheV2) loadFromDBToCache(ctx context.Context, cache ecache.Cache) error {
	// 从数据库加载数据
	entries, err := mlc.dataLoader(ctx)
	if err != nil {
		return err
	}
	// 保存到缓存
	for i := range entries {
		err = mlc.Set(ctx, entries[i].Key, entries[i].Val, entries[i].Expiration)
		if err != nil {
			return err
		}
	}
	return nil
}
func (mlc *MultiCacheV2) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	if !mlc.isRedisAvailable.Load() {
		// Redis不可用，写入本地缓存
		return mlc.local.Set(ctx, key, val, exp)
	}
	// Redis可用，只写入Redis
	err := mlc.redis.Set(ctx, key, val, exp)
	mlc.redisCrashDetector.Add(err != nil)
	if err != nil && mlc.redisCrashDetector.IsConditionMet() {
		// Redis检测到崩溃，启动降级流程
		mlc.handleRedisCrashEvent(ctx)
	}
	return err
}
func (mlc *MultiCacheV2) handleRedisCrashEvent(ctx context.Context) {
	mlc.mu.Lock()
	defer mlc.mu.Unlock()

	// 已经处于不可用状态，避免重复处理
	if !mlc.isRedisAvailable.Load() {
		return
	}

	// 标记Redis不可用
	mlc.isRedisAvailable.Store(false)

	// 立即从数据库加载数据到本地缓存
	if err := mlc.loadFromDBToCache(ctx, mlc.local); err != nil {
		mlc.logger.Error("从数据库加载数据到本地缓存失败", elog.FieldErr(err))
	}

	mlc.stopCtx, mlc.stopRefreshCancelFunc = context.WithCancel(context.Background())

	// 启动定时从数据库刷新本地缓存的任务
	//nolint:contextcheck // 忽略
	go mlc.refreshLocalCache(mlc.stopCtx)
}
func (m *MultiCacheV2) refreshLocalCache(ctx context.Context) {
	m.refreshTicker = time.NewTicker(m.localRefreshPeriod)
	for {
		select {
		case <-m.refreshTicker.C:
			if err := m.loadFromDBToCache(ctx, m.local); err != nil {
				m.logger.Error("从数据库加载数据到本地缓存失败", elog.FieldErr(err))
			}
		case <-ctx.Done():
			m.refreshTicker.Stop()
			return
		}
	}
}
func (m *MultiCacheV2) Get(ctx context.Context, key string) ecache.Value {
	if !m.isRedisAvailable.Load() {
		// Redis不可用，查本地缓存
		return m.local.Get(ctx, key)
	}
	// Redis可用，从Redis获取
	val := m.redis.Get(ctx, key)
	// 检查Redis是否出错（排除KeyNotFound）
	if val.Err != nil && !val.KeyNotFound() {
		m.redisCrashDetector.Add(true)
		if m.redisCrashDetector.IsConditionMet() {
			// Redis崩溃，切换到使用本地缓存
			m.handleRedisCrashEvent(ctx)
		}
	} else {
		// Redis正常响应
		m.redisCrashDetector.Add(false)
	}
	return val
}
