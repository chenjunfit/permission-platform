package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

type RoleRepository interface {
	Create(ctx context.Context, role domain.Role) (domain.Role, error)
	FindByBizID(ctx context.Context, bizID int64, offset, limit int) ([]domain.Role, error)
	FindByBizIDAndID(ctx context.Context, bizID, Id int64) (domain.Role, error)
	FindByBizIDAndType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]domain.Role, error)
	UpdateByBizIDAndID(ctx context.Context, role domain.Role) (domain.Role, error)
	DeleteByBizIDAndID(ctx context.Context, bizID, Id int64) error
}

type roleRepository struct {
	roleDao dao.RoleDAO
}

func NewRoleRepository(roleDAO dao.RoleDAO) RoleRepository {
	return &roleRepository{
		roleDao: roleDAO,
	}
}
func (r *roleRepository) Create(ctx context.Context, role domain.Role) (domain.Role, error) {
	created, err := r.roleDao.Create(ctx, r.toEntity(role))
	if err != nil {
		return domain.Role{}, err
	}
	return r.toDomain(created), nil
}

func (r *roleRepository) FindByBizID(ctx context.Context, bizID int64, offset, limit int) ([]domain.Role, error) {
	roles, err := r.roleDao.FindByBizID(ctx, bizID, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(roles, func(_ int, src dao.Role) domain.Role {
		return r.toDomain(src)
	}), nil
}

func (r *roleRepository) FindByBizIDAndID(ctx context.Context, bizID, Id int64) (domain.Role, error) {
	role, err := r.roleDao.FindByBizIDAndID(ctx, bizID, Id)
	if err != nil {
		return domain.Role{}, err
	}
	return r.toDomain(role), nil

}

func (r *roleRepository) FindByBizIDAndType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]domain.Role, error) {
	roles, err := r.roleDao.FindByBizIDAndType(ctx, bizID, roleType, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(roles, func(_ int, src dao.Role) domain.Role {
		return r.toDomain(src)
	}), nil
}

func (r *roleRepository) UpdateByBizIDAndID(ctx context.Context, role domain.Role) (domain.Role, error) {
	err := r.roleDao.UpdateByBizIDAndID(ctx, r.toEntity(role))
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) DeleteByBizIDAndID(ctx context.Context, bizID, Id int64) error {
	return r.roleDao.DeleteByBizIDAndID(ctx, bizID, Id)
}

func (r *roleRepository) toEntity(role domain.Role) dao.Role {
	return dao.Role{
		BizID:       role.BizID,
		Type:        role.Type,
		Name:        role.Name,
		Description: role.Description,
		Metadata:    role.Metadata,
		Ctime:       role.Ctime,
		Utime:       role.Utime,
	}
}
func (r *roleRepository) toDomain(role dao.Role) domain.Role {
	return domain.Role{
		BizID:       role.BizID,
		Type:        role.Type,
		Name:        role.Name,
		Description: role.Description,
		Metadata:    role.Metadata,
		Ctime:       role.Ctime,
		Utime:       role.Utime,
	}
}
