package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/pkg/cache"
	"time"
)

const (
	day               = 24 * time.Hour
	defaultExpiration = 36500 * day
)

type UserPermissionCache interface {
	Get(ctx context.Context, bizID, userID int64) ([]domain.UserPermission, error)
	Set(ctx context.Context, permissions []domain.UserPermission) error
}
type userPermissionCache struct {
	c            cache.Cache
	cacheKeyFunc func(bizID, userID int64) string
}

func (u *userPermissionCache) Get(ctx context.Context, bizID, userID int64) ([]domain.UserPermission, error) {
	val := u.c.Get(ctx, u.cacheKeyFunc(bizID, userID))
	if val.Err != nil {
		if val.KeyNotFound() {
			return nil, fmt.Errorf("%w", errors.New("Key not Found"))
		}
		return nil, val.Err
	}
	var res []domain.UserPermission
	err := val.JSONScan(&res)
	return res, err
}

func (u *userPermissionCache) Set(ctx context.Context, permissions []domain.UserPermission) error {
	if len(permissions) == 0 {
		return nil
	}
	value, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	bizID, userID := permissions[0].BizID, permissions[0].UserID
	return u.c.Set(ctx, u.cacheKeyFunc(bizID, userID), value, defaultExpiration)
}

func NewUserPermissionCache(c cache.Cache, cacheKeyFunc func(bizID, userID int64) string) UserPermissionCache {
	return &userPermissionCache{
		c:            c,
		cacheKeyFunc: cacheKeyFunc,
	}
}
