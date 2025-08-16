package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

type BusinessConfigRepository interface {
	Create(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error)

	Find(ctx context.Context, offset, limit int) ([]domain.BusinessConfig, error)
	FindByID(ctx context.Context, id int64) (domain.BusinessConfig, error)

	UpdateToken(ctx context.Context, id int64, token string) error
	Update(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error)

	Delete(ctx context.Context, id int64) error
}

func NewBusinessConfigRepository(businessConfigDao dao.BusinessConfigDAO) BusinessConfigRepository {
	return &businessConfigRepository{businessConfigDAO: businessConfigDao}
}

type businessConfigRepository struct {
	businessConfigDAO dao.BusinessConfigDAO
}

func (b *businessConfigRepository) Create(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error) {
	created, err := b.businessConfigDAO.Create(ctx, b.toEntity(config))
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	return b.toDomain(created), nil
}

func (b *businessConfigRepository) Find(ctx context.Context, offset, limit int) ([]domain.BusinessConfig, error) {
	list, err := b.businessConfigDAO.Find(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return slice.Map(list, func(_ int, src dao.BusinessConfig) domain.BusinessConfig {
		return b.toDomain(src)
	}), nil
}

func (b *businessConfigRepository) FindByID(ctx context.Context, id int64) (domain.BusinessConfig, error) {
	config, err := b.businessConfigDAO.GetByID(ctx, id)
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	return b.toDomain(config), nil
}

func (b *businessConfigRepository) UpdateToken(ctx context.Context, id int64, token string) error {
	return b.businessConfigDAO.UpdateToken(ctx, id, token)
}

func (b *businessConfigRepository) Update(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error) {
	err := b.businessConfigDAO.Update(ctx, b.toEntity(config))
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	return config, nil
}

func (b *businessConfigRepository) Delete(ctx context.Context, id int64) error {
	return b.businessConfigDAO.Delete(ctx, id)
}
func (b *businessConfigRepository) toEntity(bc domain.BusinessConfig) dao.BusinessConfig {
	return dao.BusinessConfig{
		ID:        bc.ID,
		OwnerID:   bc.OwnerID,
		OwnerType: bc.OwnerType,
		Name:      bc.Name,
		RateLimit: bc.RateLimit,
		Token:     bc.Token,
		Ctime:     bc.Ctime,
		Utime:     bc.Utime,
	}
}

func (b *businessConfigRepository) toDomain(bc dao.BusinessConfig) domain.BusinessConfig {
	return domain.BusinessConfig{
		ID:        bc.ID,
		OwnerID:   bc.OwnerID,
		OwnerType: bc.OwnerType,
		Name:      bc.Name,
		RateLimit: bc.RateLimit,
		Token:     bc.Token,
		Ctime:     bc.Ctime,
		Utime:     bc.Utime,
	}
}
