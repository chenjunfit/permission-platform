package repository

import (
	"context"
	"github.com/ecodeclub/ekit/mapx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
	"time"
)

var _ UserPermissionRepository = (*userPermissionRepository)(nil)

type UserPermissionRepository interface {
	Create(ctx context.Context, permission domain.UserPermission) (domain.UserPermission, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.UserPermission, error)
	FindByBizIdAndUserID(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error)
	DeleteByBizIdAndID(ctx context.Context, bizId, id int64) error
	//返回用户的个人权限，个人角色以及包含角色的权限
	GetALLUserPermission(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error)
}

type userPermissionRepository struct {
	userRoleDao       dao.UserRoleDAO
	roleInclusionDao  dao.RoleInclusionDAO
	rolePermissionDao dao.RolePermissionDAO
	userPermissionDao dao.UserPermissionDAO
}

func (u *userPermissionRepository) Create(ctx context.Context, permission domain.UserPermission) (domain.UserPermission, error) {
	created, err := u.userPermissionDao.Create(ctx, u.toEntity(permission))
	if err != nil {
		return domain.UserPermission{}, nil
	}
	return u.toDomain(created), nil
}

func (u *userPermissionRepository) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.UserPermission, error) {
	ups, err := u.userPermissionDao.FindByBizID(ctx, bizId, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(ups, func(idx int, src dao.UserPermission) domain.UserPermission {
		return u.toDomain(src)
	}), nil
}

func (u *userPermissionRepository) FindByBizIdAndUserID(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error) {
	ups, err := u.userPermissionDao.FindByBizIdAndUserId(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}
	return slice.Map(ups, func(idx int, src dao.UserPermission) domain.UserPermission {
		return u.toDomain(src)
	}), nil
}

func (u *userPermissionRepository) DeleteByBizIdAndID(ctx context.Context, bizId, id int64) error {
	return u.userPermissionDao.DeleteBizIdAndId(ctx, bizId, id)
}

func (u *userPermissionRepository) GetALLUserPermission(ctx context.Context, bizId, userId int64) ([]domain.UserPermission, error) {
	//获取个人权限
	userPermissions, err := u.userPermissionDao.FindByBizIdAndUserId(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}
	perms := slice.Map(userPermissions, func(idx int, src dao.UserPermission) domain.UserPermission {
		return u.toDomain(src)
	})
	//获取角色以及包含的角色
	roleIds, err := u.GetAllRoleIds(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}
	//获取所有角色的权限
	allRoleUserPermissions, err := u.GetAllRolePermissions(ctx, bizId, userId, roleIds)
	if err != nil {
		return nil, err
	}
	perms = append(perms, allRoleUserPermissions...)
	return perms, nil
}
func (u *userPermissionRepository) GetAllRolePermissions(ctx context.Context, bizId, userId int64, roleIds []int64) ([]domain.UserPermission, error) {
	if len(roleIds) == 0 {
		return []domain.UserPermission{}, nil
	}
	rolePermissions, err := u.rolePermissionDao.FindByBizIDAndRoleIds(ctx, bizId, roleIds)
	if err != nil {
		return []domain.UserPermission{}, err
	}
	return slice.Map(rolePermissions, func(idx int, src dao.RolePermission) domain.UserPermission {
		return domain.UserPermission{
			ID:     0,
			BizID:  bizId,
			UserID: userId,
			Permission: domain.Permission{
				ID:    src.PermissionID,
				BizID: bizId,
				Resource: domain.Resource{
					BizID: bizId,
					Type:  src.ResourceType,
					Key:   src.ResourceKey,
				},
				Action: src.PermissionAction,
			},
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().AddDate(100, 0, 0).Unix(),
			Effect:    domain.EffectAllow,
			Ctime:     src.Ctime,
			Utime:     src.Utime,
		}
	}), nil

}
func (u *userPermissionRepository) GetAllRoleIds(ctx context.Context, bizId, userId int64) ([]int64, error) {
	//直接关联的角色
	directUserRoles, err := u.userRoleDao.FindByBizIDAndUserID(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}
	allRoleIds := make(map[int64]any, len(directUserRoles))
	directUserRoleIds := slice.Map(directUserRoles, func(idx int, src dao.UserRole) int64 {
		allRoleIds[src.RoleID] = struct{}{}
		return src.RoleID
	})
	includeIds := directUserRoleIds
	for {
		roleInclusions, err := u.roleInclusionDao.FindByBizIdAndIncludingIds(ctx, bizId, includeIds)
		if err != nil {
			return nil, err
		}
		if len(roleInclusions) == 0 {
			break
		}
		includeIds = slice.Map(roleInclusions, func(idx int, src dao.RoleInclusion) int64 {
			allRoleIds[src.IncludedRoleID] = struct{}{}
			return src.IncludedRoleID
		})
	}
	return mapx.Keys(allRoleIds), nil
}

func NewUserPermissionRepository(userPermissionDao dao.UserPermissionDAO) UserPermissionRepository {
	return &userPermissionRepository{
		userPermissionDao: userPermissionDao,
	}
}

func (u *userPermissionRepository) toEntity(up domain.UserPermission) dao.UserPermission {
	return dao.UserPermission{
		ID:               up.ID,
		BizID:            up.BizID,
		UserID:           up.UserID,
		PermissionID:     up.Permission.ID,
		PermissionName:   up.Permission.Name,
		ResourceType:     up.Permission.Resource.Type,
		ResourceKey:      up.Permission.Resource.Key,
		PermissionAction: up.Permission.Action,
		StartTime:        up.StartTime,
		EndTime:          up.EndTime,
		Effect:           up.Effect.String(),
		Ctime:            up.Ctime,
		Utime:            up.Utime,
	}
}
func (u *userPermissionRepository) toDomain(up dao.UserPermission) domain.UserPermission {
	return domain.UserPermission{
		ID:     up.ID,
		BizID:  up.BizID,
		UserID: up.UserID,
		Permission: domain.Permission{
			ID:   up.PermissionID,
			Name: up.PermissionName,
			Resource: domain.Resource{
				Type: up.ResourceType,
				Key:  up.ResourceKey,
			},
			Action: up.PermissionAction,
		},
		StartTime: up.StartTime,
		EndTime:   up.EndTime,
		Effect:    domain.Effect(up.Effect),
		Ctime:     up.Ctime,
		Utime:     up.Utime,
	}
}
