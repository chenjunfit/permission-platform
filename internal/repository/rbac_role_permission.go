package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
	"time"
)

var _ RolePermissionRepository = (*rolePermissionRepository)(nil)

type RolePermissionRepository interface {
	Create(ctx context.Context, permission domain.RolePermission) (domain.RolePermission, error)
	FindByBizID(ctx context.Context, bizID int64) ([]domain.RolePermission, error)
	FindByBizIDAndRoleIDs(ctx context.Context, bizID int64, roleIDs []int64) ([]domain.RolePermission, error)

	DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error
}

type rolePermissionRepository struct {
	rolePermissionDao dao.RolePermissionDAO
	logger            *elog.Component
}

func (r *rolePermissionRepository) Create(ctx context.Context, rolePermission domain.RolePermission) (domain.RolePermission, error) {
	now := time.Now().Unix()
	rolePermission.Utime = now
	rolePermission.Ctime = now
	created, err := r.rolePermissionDao.Create(ctx, r.toEntity(rolePermission))
	if err != nil {
		r.logger.Error("为角色添加权限失败",
			elog.FieldErr(err),
			elog.Any("rolePermission", rolePermission),
			elog.Int64("roleId", rolePermission.Role.ID),
			elog.Int64("permissionId", rolePermission.Permission.ID),
			elog.Int64("bizID", rolePermission.BizID),
		)
		return domain.RolePermission{}, err
	} else {
		r.logger.Info("为角色添加权限",
			elog.Any("rolePermission", rolePermission),
			elog.Int64("roleId", rolePermission.Role.ID),
			elog.Int64("permissionId", rolePermission.Permission.ID),
			elog.Int64("bizID", rolePermission.BizID),
		)
	}
	return r.toDomain(created), err
}

func (r *rolePermissionRepository) FindByBizID(ctx context.Context, bizID int64) ([]domain.RolePermission, error) {
	rolePermissions, err := r.rolePermissionDao.FindByBizID(ctx, bizID)
	if err != nil {
		return nil, err
	}
	return slice.Map(rolePermissions, func(idx int, src dao.RolePermission) domain.RolePermission {
		return r.toDomain(src)
	}), nil
}

func (r *rolePermissionRepository) FindByBizIDAndRoleIDs(ctx context.Context, bizID int64, roleIDs []int64) ([]domain.RolePermission, error) {
	rolePermissions, err := r.rolePermissionDao.FindByBizIDAndRoleIds(ctx, bizID, roleIDs)
	if err != nil {
		return nil, err
	}
	return slice.Map(rolePermissions, func(idx int, src dao.RolePermission) domain.RolePermission {
		return r.toDomain(src)
	}), nil
}
func (r *rolePermissionRepository) FindByBizIDAndID(ctx context.Context, bizID, id int64) (domain.RolePermission, error) {
	rp, err := r.rolePermissionDao.FindByBizIdAndID(ctx, bizID, id)
	if err != nil {
		return domain.RolePermission{}, err
	}
	return r.toDomain(rp), nil
}
func (r *rolePermissionRepository) DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error {
	err := r.rolePermissionDao.DeleteByBizIDAndID(ctx, bizID, id)
	if err != nil {
		r.logger.Error("为角色删除权限失败",
			elog.FieldErr(err),
			elog.Int64("biz_id", bizID),
			elog.Int64("role_permission_id", id),
		)
	} else {
		r.logger.Info("为角色删除权限成功",
			elog.Int64("biz_id", bizID),
			elog.Int64("role_permission_id", id),
		)
	}
	return err
}

func NewRolePermissionRepository(rolePermissionDao dao.RolePermissionDAO) RolePermissionRepository {
	return &rolePermissionRepository{
		rolePermissionDao: rolePermissionDao,
		logger:            elog.DefaultLogger,
	}
}

func (r *rolePermissionRepository) toEntity(rp domain.RolePermission) dao.RolePermission {
	return dao.RolePermission{
		ID:               rp.ID,
		BizID:            rp.BizID,
		RoleID:           rp.Role.ID,
		PermissionID:     rp.Permission.ID,
		RoleName:         rp.Role.Name,
		RoleType:         rp.Role.Type,
		ResourceType:     rp.Permission.Resource.Type,
		ResourceKey:      rp.Permission.Resource.Key,
		PermissionAction: rp.Permission.Action,
		Ctime:            rp.Ctime,
		Utime:            rp.Utime,
	}
}
func (r *rolePermissionRepository) toDomain(rp dao.RolePermission) domain.RolePermission {
	return domain.RolePermission{
		ID:    rp.ID,
		BizID: rp.BizID,
		Role: domain.Role{
			ID:   rp.RoleID,
			Type: rp.RoleType,
			Name: rp.RoleName,
		},
		Permission: domain.Permission{
			ID:    rp.PermissionID,
			BizID: rp.BizID,
			Resource: domain.Resource{
				Type: rp.ResourceType,
				Key:  rp.ResourceKey,
			},
			Action: rp.PermissionAction,
		},
		Ctime: rp.Ctime,
		Utime: rp.Utime,
	}
}
