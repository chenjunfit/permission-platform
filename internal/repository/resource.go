package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

type ResourceRepository interface {
	Create(ctx context.Context, resource domain.Resource) (domain.Resource, error)
	UpdateByBizIDAndID(ctx context.Context, resource domain.Resource) (domain.Resource, error)
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error

	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.Resource, error)
	FindByBizIDAndID(ctx context.Context, bizId, id int64) (domain.Resource, error)
	FindByBizIDAndTypeAndKey(ctx context.Context, bizId int64, resourceType, resourceKey string) (domain.Resource, error)
}

type resourceRepository struct {
	resourceDao dao.ResourceDao
}

func NewResourceRepository(resourceDao dao.ResourceDao) ResourceRepository {
	return &resourceRepository{resourceDao: resourceDao}
}

func (r *resourceRepository) Create(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	created, err := r.resourceDao.Create(ctx, r.toEntity(resource))
	if err != nil {
		return domain.Resource{}, err
	}
	return r.toDomain(created), nil
}

func (r *resourceRepository) UpdateByBizIDAndID(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	err := r.resourceDao.UpdateByBizIDAndID(ctx, r.toEntity(resource))
	if err != nil {
		return domain.Resource{}, err
	}
	return resource, nil
}

func (r *resourceRepository) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return r.resourceDao.DeleteByBizIDAndID(ctx, bizId, id)
}

func (r *resourceRepository) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]domain.Resource, error) {
	rescources, err := r.resourceDao.FindByBizID(ctx, bizId, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(rescources, func(idx int, src dao.Resource) domain.Resource {
		return r.toDomain(src)
	}), nil
}

func (r *resourceRepository) FindByBizIDAndID(ctx context.Context, bizId, id int64) (domain.Resource, error) {
	resource, err := r.resourceDao.FindByBizIDAndID(ctx, bizId, id)
	if err != nil {
		return domain.Resource{}, err
	}
	return r.toDomain(resource), nil
}

func (r *resourceRepository) FindByBizIDAndTypeAndKey(ctx context.Context, bizId int64, resourceType, resourceKey string) (domain.Resource, error) {
	resource, err := r.resourceDao.FindByBizIDAndTypeAndKey(ctx, bizId, resourceType, resourceKey)
	if err != nil {
		return domain.Resource{}, err
	}
	return r.toDomain(resource), nil
}
func (r *resourceRepository) toEntity(res domain.Resource) dao.Resource {
	return dao.Resource{
		ID:          res.ID,
		BizID:       res.BizID,
		Type:        res.Type,
		Key:         res.Key,
		Name:        res.Name,
		Description: res.Description,
		Metadata:    res.Metadata,
		Ctime:       res.Ctime,
		Utime:       res.Utime,
	}
}

func (r *resourceRepository) toDomain(res dao.Resource) domain.Resource {
	return domain.Resource{
		ID:          res.ID,
		BizID:       res.BizID,
		Type:        res.Type,
		Key:         res.Key,
		Name:        res.Name,
		Description: res.Description,
		Metadata:    res.Metadata,
		Ctime:       res.Ctime,
		Utime:       res.Utime,
	}
}
