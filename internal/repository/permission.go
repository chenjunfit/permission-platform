package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	FindPermissions(ctx context.Context, bizId int64, resourceType, resourceKey string, action []string) ([]domain.Permission, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.Permission, error)
	FindByBizIDANdID(ctx context.Context, bizId, id int64) (domain.Permission, error)
	UpdateByBizIDAndID(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error
}

type permissionRepository struct {
	permissionDao dao.PermissionDAO
}

func NewPermissionRepository(permissionDao dao.PermissionDAO) PermissionRepository {
	return &permissionRepository{
		permissionDao: permissionDao,
	}
}
func (p *permissionRepository) FindByBizIDANdID(ctx context.Context, bizId, id int64) (domain.Permission, error) {
	permission, err := p.permissionDao.FindByBizIDAndID(ctx, bizId, id)
	if err != nil {
		return domain.Permission{}, err
	}
	return p.toDomain(permission), nil
}
func (p *permissionRepository) Create(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	created, err := p.permissionDao.Create(ctx, p.toEntity(permission))
	if err != nil {
		return domain.Permission{}, err
	}
	return p.toDomain(created), nil
}

func (p *permissionRepository) FindPermissions(ctx context.Context, bizId int64, resourceType, resourceKey string, action []string) ([]domain.Permission, error) {
	permissions, err := p.permissionDao.FindPermissions(ctx, bizId, resourceType, resourceKey, action)
	if err != nil {
		return nil, err
	}
	list := slice.Map(permissions, func(idx int, src dao.Permission) domain.Permission {
		return p.toDomain(src)
	})
	return list, nil
}

func (p *permissionRepository) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.Permission, error) {
	permissions, err := p.permissionDao.FindByBizID(ctx, bizId, offset, limit)
	if err != nil {
		return nil, err
	}
	list := slice.Map(permissions, func(idx int, src dao.Permission) domain.Permission {
		return p.toDomain(src)
	})
	return list, nil

}

func (p *permissionRepository) UpdateByBizIDAndID(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	err := p.permissionDao.UpdateByBizIDAndID(ctx, p.toEntity(permission))
	if err != nil {
		return domain.Permission{}, err
	}
	return permission, nil
}

func (p *permissionRepository) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return p.permissionDao.DeleteByBizIDAndID(ctx, bizId, id)
}
func (r *permissionRepository) toEntity(p domain.Permission) dao.Permission {
	return dao.Permission{
		ID:           p.ID,
		BizID:        p.BizID,
		Name:         p.Name,
		Description:  p.Description,
		ResourceID:   p.Resource.ID,
		ResourceType: p.Resource.Type,
		ResourceKey:  p.Resource.Key,
		Action:       p.Action,
		Metadata:     p.Metadata,
		Ctime:        p.Ctime,
		Utime:        p.Utime,
	}
}

func (r *permissionRepository) toDomain(p dao.Permission) domain.Permission {
	return domain.Permission{
		ID:          p.ID,
		BizID:       p.BizID,
		Name:        p.Name,
		Description: p.Description,
		Resource: domain.Resource{
			ID:   p.ResourceID,
			Type: p.ResourceType,
			Key:  p.ResourceKey,
		},
		Action:   p.Action,
		Metadata: p.Metadata,
		Ctime:    p.Ctime,
		Utime:    p.Utime,
	}
}
