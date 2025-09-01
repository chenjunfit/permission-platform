package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/event/permission"
	"github.com/permission-dev/internal/repository/cache"
)

var (
	_ UserPermissionRepository    = (*UserPermissionCachedRepository)(nil)
	_ UserPermissionCacheReloader = (*UserPermissionCachedRepository)(nil)
)

type UserPermissionCacheReloader interface {
	Reload(ctx context.Context, user []domain.User) error
}
type UserPermissionCachedRepository struct {
	repo     UserPermissionRepository
	cache    cache.UserPermissionCache
	producer permission.UserPermissionEventProducer
	logger   *elog.Component
}

func (u *UserPermissionCachedRepository) FindByBizIDAndID(ctx context.Context, bizId, id int64) (domain.UserPermission, error) {
	return domain.UserPermission{}, nil
}

func (u *UserPermissionCachedRepository) Reload(ctx context.Context, user []domain.User) error {
	var evt permission.UserPermissionEvent
	evt.Permissions = make(map[int64]permission.UserPermission)
	for index := range user {
		perms, err := u.repo.GetALLUserPermission(ctx, user[index].BizID, user[index].ID)
		if err != nil {
			return err
		}
		err = u.cache.Set(ctx, perms)
		if err != nil {
			elog.Error("重新加载全部用户权限失败",
				elog.FieldErr(err),
				elog.Any("bizID:", user[index].BizID),
				elog.Any("userID:", user[index].ID),
			)
		} else {
			evt.Permissions[user[index].ID] = permission.UserPermission{
				UserID: user[index].ID,
				BizID:  user[index].BizID,
				Permissions: slice.Map(perms, func(idx int, src domain.UserPermission) permission.PermissionV1 {
					return permission.PermissionV1{
						Resource: permission.Resource{
							Key:  src.Permission.Resource.Key,
							Type: src.Permission.Resource.Type,
						},
						Action: src.Permission.Action,
						Effect: src.Effect.String(),
					}
				}),
			}
		}

	}
	if len(evt.Permissions) > 0 {
		if err := u.producer.Produce(ctx, evt); err != nil {
			u.logger.Warn("发送用户权限事件失败",
				elog.FieldErr(err),
				elog.Any("evt", evt),
			)
		}
	}
	return nil
}
func (u *UserPermissionCachedRepository) Create(ctx context.Context, permission domain.UserPermission) (domain.UserPermission, error) {
	created, err := u.repo.Create(ctx, permission)
	if err != nil {
		return domain.UserPermission{}, err
	}
	if err1 := u.Reload(ctx, []domain.User{{ID: created.UserID, BizID: created.BizID}}); err1 != nil {
		u.logger.Warn("创建用户权限成功后，重新加载缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", created.BizID),
			elog.Any("userID", created.UserID),
		)
	}
	return created, err
}

func (u *UserPermissionCachedRepository) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.UserPermission, error) {
	return u.repo.FindByBizID(ctx, bizId, offset, limit)
}

func (u *UserPermissionCachedRepository) FindByBizIdAndUserID(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error) {
	perms, err := u.cache.Get(ctx, bizId, userId)
	if err == nil {
		return perms, nil
	}
	perms, err = u.repo.FindByBizIdAndUserID(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}
	if err1 := u.cache.Set(ctx, perms); err1 != nil {
		u.logger.Warn("查找用户权限成功后，重新设置缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", bizId),
			elog.Any("userID", userId),
		)
	}
	return perms, nil
}

func (u *UserPermissionCachedRepository) DeleteByBizIdAndID(ctx context.Context, bizId, id int64) error {
	deleted, err := u.repo.FindByBizIDAndID(ctx, bizId, id)
	if err != nil {
		return err
	}
	err = u.repo.DeleteByBizIdAndID(ctx, bizId, id)
	if err != nil {
		return err
	}
	if err1 := u.Reload(ctx, []domain.User{{ID: deleted.UserID, BizID: deleted.BizID}}); err1 != nil {
		u.logger.Warn("删除用户权限成功后，重新加载缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", bizId),
			elog.Any("userID", deleted.UserID),
		)
	}
	return nil
}

func (u *UserPermissionCachedRepository) GetALLUserPermission(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error) {
	perms, err := u.cache.Get(ctx, bizId, userId)
	if err == nil {
		return perms, nil
	}

	perms, err = u.repo.GetALLUserPermission(ctx, bizId, userId)
	if err != nil {
		u.logger.Error("从数据库中查找用户全部权限失败",
			elog.FieldErr(err),
			elog.Any("bizID", bizId),
			elog.Any("userID", userId),
		)
		return nil, err
	}

	if err1 := u.cache.Set(ctx, perms); err1 != nil {
		u.logger.Warn("存储用户全部权限到缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", bizId),
			elog.Any("userID", userId),
		)
	}
	return perms, nil
}

func NewUserPermissionCachedRepository(
	repo UserPermissionRepository,
	cache cache.UserPermissionCache,
	producer permission.UserPermissionEventProducer,
) *UserPermissionCachedRepository {
	return &UserPermissionCachedRepository{
		repo:     repo,
		cache:    cache,
		producer: producer,
		logger:   elog.DefaultLogger.With(elog.FieldName("UserPermissionCachedRepository")),
	}
}
