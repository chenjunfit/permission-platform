package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

var _ RoleIncludeRepository = (*roleIncludeRepository)(nil)

type RoleIncludeRepository interface {
	Create(ctx context.Context, inclusion domain.RoleInclusion) (domain.RoleInclusion, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.RoleInclusion, error)
	FindByBizIdAndIncludingIds(ctx context.Context, bizId int64, IncludingIds []int64) ([]domain.RoleInclusion, error)
	FindByBizIdAndIncludedIds(ctx context.Context, bizId int64, IncludedIds []int64) ([]domain.RoleInclusion, error)
	DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error
	FindByBizIDAndID(ctx context.Context, bizID, id int64) (domain.RoleInclusion, error)
}

type roleIncludeRepository struct {
	roleInclusionDao dao.RoleInclusionDAO
}

func (r *roleIncludeRepository) FindByBizIDAndID(ctx context.Context, bizID, id int64) (domain.RoleInclusion, error) {
	ri, err := r.roleInclusionDao.FindByBizIDAndID(ctx, bizID, id)
	if err != nil {
		return domain.RoleInclusion{}, err
	}
	return r.toDomain(ri), err
}

func (r *roleIncludeRepository) Create(ctx context.Context, inclusion domain.RoleInclusion) (domain.RoleInclusion, error) {
	create, err := r.roleInclusionDao.Create(ctx, r.toEntity(inclusion))
	if err != nil {
		return domain.RoleInclusion{}, err
	}
	return r.toDomain(create), nil
}

func (r *roleIncludeRepository) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.RoleInclusion, error) {
	roleInclusions, err := r.roleInclusionDao.FindByBizID(ctx, bizId, offset, limit)
	if err != nil {
		return nil, err
	}

	return slice.Map(roleInclusions, func(_ int, src dao.RoleInclusion) domain.RoleInclusion {
		return r.toDomain(src)
	}), nil
}

func (r *roleIncludeRepository) FindByBizIdAndIncludingIds(ctx context.Context, bizId int64, IncludingIds []int64) ([]domain.RoleInclusion, error) {
	roleInclusions, err := r.roleInclusionDao.FindByBizIdAndIncludingIds(ctx, bizId, IncludingIds)
	if err != nil {
		return nil, err
	}

	return slice.Map(roleInclusions, func(_ int, src dao.RoleInclusion) domain.RoleInclusion {
		return r.toDomain(src)
	}), nil
}

func (r *roleIncludeRepository) FindByBizIdAndIncludedIds(ctx context.Context, bizId int64, IncludedIds []int64) ([]domain.RoleInclusion, error) {
	roleInclusions, err := r.roleInclusionDao.FindByBizIdAndIncludedIds(ctx, bizId, IncludedIds)
	if err != nil {
		return nil, err
	}

	return slice.Map(roleInclusions, func(_ int, src dao.RoleInclusion) domain.RoleInclusion {
		return r.toDomain(src)
	}), nil
}

func (r *roleIncludeRepository) DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error {
	return r.roleInclusionDao.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *roleIncludeRepository) toEntity(inclusion domain.RoleInclusion) dao.RoleInclusion {
	return dao.RoleInclusion{
		ID:                inclusion.ID,
		BizID:             inclusion.BizID,
		IncludingRoleID:   inclusion.IncludingRole.ID,
		IncludingRoleType: inclusion.IncludingRole.Type,
		IncludingRoleName: inclusion.IncludingRole.Name,
		IncludedRoleID:    inclusion.IncludedRole.ID,
		IncludedRoleType:  inclusion.IncludedRole.Type,
		IncludedRoleName:  inclusion.IncludedRole.Name,
		Ctime:             inclusion.Ctime,
		Utime:             inclusion.Utime,
	}
}
func (r *roleIncludeRepository) toDomain(ri dao.RoleInclusion) domain.RoleInclusion {
	return domain.RoleInclusion{
		ID:    ri.ID,
		BizID: ri.BizID,
		IncludingRole: domain.Role{
			ID:   ri.IncludingRoleID,
			Type: ri.IncludingRoleType,
			Name: ri.IncludingRoleName,
		},
		IncludedRole: domain.Role{
			ID:   ri.IncludedRoleID,
			Type: ri.IncludedRoleType,
			Name: ri.IncludedRoleName,
		},
		Ctime: ri.Ctime,
		Utime: ri.Utime,
	}
}

func NewRoleIncludeRepository(roleInclusionDao dao.RoleInclusionDAO) RoleIncludeRepository {
	return &roleIncludeRepository{
		roleInclusionDao: roleInclusionDao,
	}
}
