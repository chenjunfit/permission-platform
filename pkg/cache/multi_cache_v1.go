package cache

import (
	"context"
	"github.com/ecodeclub/ecache"
	eredis "github.com/ecodeclub/ecache/redis"
	"github.com/permission-dev/pkg/bitring"
	"github.com/redis/go-redis/v9"
	"sync/atomic"
	"time"
)

type MultiCacheV1 struct {
	redis                ecache.Cache
	local                ecache.Cache
	isRedisAvailable     atomic.Bool
	redisPingTimeOut     time.Duration
	redisHealCheckPeriod time.Duration
	redisCrashDetector   *bitring.BitRing
}

func NewMultiCacheV1(
	local ecache.Cache,
	redisClient redis.Cmdable,
	redisPingTimeOut time.Duration,
	redisHealCheckPeriod time.Duration,
	redisCrashDetector *bitring.BitRing,
) *MultiCacheV1 {
	cacheV1 := &MultiCacheV1{
		redis: &ecache.NamespaceCache{
			C:         eredis.NewCache(redisClient),
			Namespace: "permission-platform",
		},
		local:                local,
		redisPingTimeOut:     redisPingTimeOut,
		redisHealCheckPeriod: redisHealCheckPeriod,
		redisCrashDetector:   redisCrashDetector,
	}
	cacheV1.isRedisAvailable.Store(true)
	go cacheV1.redisHealthCheck(redisClient)
	return cacheV1
}

func (m *MultiCacheV1) redisHealthCheck(redis redis.Cmdable) {
	ticker := time.NewTicker(m.redisHealCheckPeriod)
	defer ticker.Stop()

	for range ticker.C {
		if !m.isRedisAvailable.Load() {
			ctx, cancel := context.WithTimeout(context.Background(), m.redisPingTimeOut)
			if err := redis.Ping(ctx); err == nil {
				m.handleRedisRecoveryEvent()
			}
			cancel()
		}
	}
}
func (m *MultiCacheV1) handleRedisRecoveryEvent() {
	// 标记Redis已恢复
	m.isRedisAvailable.CompareAndSwap(false, true)
	m.redisCrashDetector.Reset()
}
func (m *MultiCacheV1) handleRedisCrashEvent() {
	// 标记Redis不可用
	m.isRedisAvailable.CompareAndSwap(true, false)
}

func (m *MultiCacheV1) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	_, err := m.local.SetNX(ctx, key, value, exp)
	if err != nil {
		return err
	}
	if !m.isRedisAvailable.Load() {
		return nil
	}
	err = m.redis.Set(ctx, key, value, exp)
	m.redisCrashDetector.Add(err != nil)
	if err != nil && m.redisCrashDetector.IsConditionMet() {
		// Redis检测到崩溃，启动降级流程
		m.handleRedisCrashEvent()
	}
	return err
}
func (m *MultiCacheV1) Get(ctx context.Context, key string) ecache.Value {
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
			m.handleRedisCrashEvent()
		}
	} else {
		// Redis正常响应
		m.redisCrashDetector.Add(false)
	}
	return val
}
