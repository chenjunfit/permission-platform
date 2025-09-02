package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

var _ UserRoleRepository = (*userRoleRepository)(nil)

type UserRoleRepository interface {
	Create(ctx context.Context, userRole domain.UserRole) (domain.UserRole, error)
	FindByBizID(ctx context.Context, bizId int64) ([]domain.UserRole, error)
	FindByBizIDAndUserID(ctx context.Context, bizId, userId int64) ([]domain.UserRole, error)
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error
}

type userRoleRepository struct {
	userRoleDao dao.UserRoleDAO
	logger      *elog.Component
}

func (r *userRoleRepository) FindByBizIDAndRoleIDs(ctx context.Context, bizID int64, roleIDs []int64) ([]domain.UserRole, error) {
	userRoles, err := r.userRoleDao.FindByBizIDAndRoleIDs(ctx, bizID, roleIDs)
	if err != nil {
		return nil, err
	}
	return slice.Map(userRoles, func(_ int, src dao.UserRole) domain.UserRole {
		return r.toDomain(src)
	}), nil
}
func (u *userRoleRepository) FindByBizIDAndID(ctx context.Context, bizID, id int64) (domain.UserRole, error) {
	ur, err := u.userRoleDao.FindByBizIDAndID(ctx, bizID, id)
	if err != nil {
		return domain.UserRole{}, err
	}
	return u.toDomain(ur), nil
}
func (ur *userRoleRepository) Create(ctx context.Context, userRole domain.UserRole) (domain.UserRole, error) {
	created, err := ur.userRoleDao.Create(ctx, ur.toEntity(userRole))
	if err != nil {
		ur.logger.Error("授予用户角色权限失败",
			elog.Int64("biz_id:", userRole.BizID),
			elog.Int64("role_id:", userRole.Role.ID),
			elog.Int64("user_id:", userRole.UserID),
			elog.String("role_name", userRole.Role.Name),
			elog.FieldErr(err),
		)
		return domain.UserRole{}, err
	} else {
		ur.logger.Info("授予用户角色权限",
			elog.Int64("biz_id:", userRole.BizID),
			elog.Int64("role_id:", userRole.Role.ID),
			elog.Int64("user_id:", userRole.UserID),
			elog.String("role_name", userRole.Role.Name),
			elog.Any("created", created),
		)
	}
	return ur.toDomain(created), nil
}

func (ur *userRoleRepository) FindByBizID(ctx context.Context, bizId int64) ([]domain.UserRole, error) {

	userRoles, err := ur.userRoleDao.FindByBizID(ctx, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(userRoles, func(idx int, src dao.UserRole) domain.UserRole {
		return ur.toDomain(src)
	}), nil
}

func (ur *userRoleRepository) FindByBizIDAndUserID(ctx context.Context, bizId, userId int64) ([]domain.UserRole, error) {
	userRoles, err := ur.userRoleDao.FindByBizIDAndUserID(ctx, bizId, userId)
	if err != nil {
		return nil, err
	}

	return slice.Map(userRoles, func(_ int, src dao.UserRole) domain.UserRole {
		return ur.toDomain(src)
	}), nil
}

func (ur *userRoleRepository) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	err := ur.userRoleDao.DeleteByBizIDAndID(ctx, bizId, id)
	if err != nil {
		ur.logger.Error("撤销用户角色权限失败",
			elog.Int64("biz_id:", bizId),
			elog.Int64("user_role_id:", id),
			elog.FieldErr(err),
		)
	} else {
		ur.logger.Error("撤销用户角色权限",
			elog.Int64("biz_id:", bizId),
			elog.Int64("user_role_id:", id),
		)
	}
	return err
}

func NewUserRoleRepository(userDao dao.UserRoleDAO) UserRoleRepository {
	return &userRoleRepository{
		userRoleDao: userDao,
		logger:      elog.DefaultLogger,
	}
}

func (ur *userRoleRepository) toEntity(userRole domain.UserRole) dao.UserRole {
	return dao.UserRole{
		ID:        userRole.ID,
		BizID:     userRole.BizID,
		UserID:    userRole.UserID,
		RoleID:    userRole.Role.ID,
		RoleName:  userRole.Role.Name,
		RoleType:  userRole.Role.Type,
		StartTime: userRole.StartTime,
		EndTime:   userRole.EndTime,
		Ctime:     userRole.Ctime,
		Utime:     userRole.Utime,
	}
}
func (ur *userRoleRepository) toDomain(userRole dao.UserRole) domain.UserRole {
	return domain.UserRole{
		ID:     userRole.ID,
		BizID:  userRole.BizID,
		UserID: userRole.UserID,
		Role: domain.Role{
			ID:    userRole.RoleID,
			BizID: userRole.BizID,
			Type:  userRole.RoleType,
			Name:  userRole.RoleName,
		},
		StartTime: userRole.StartTime,
		EndTime:   userRole.EndTime,
		Ctime:     userRole.Ctime,
		Utime:     userRole.Utime,
	}
}
