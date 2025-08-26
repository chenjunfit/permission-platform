package repository

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

type AttributeDefinitionRepository interface {
	Create(ctx context.Context, bizId int64, definition domain.AttributeDefinition) (int64, error)
	Delete(ctx context.Context, bizId, id int64) error
	//返回bizId下所有的属性定义(env,subject,resource)
	FindByBizID(ctx context.Context, bizId int64) (domain.BizAttrDefinition, error)
	FindByBizIdAndId(ctx context.Context, bizId, id int64) (domain.AttributeDefinition, error)
}

type attributeDefinitionRepository struct {
	dao dao.AttributeDefinitionDAO
}

func (a *attributeDefinitionRepository) Create(ctx context.Context, bizId int64, definition domain.AttributeDefinition) (int64, error) {
	id, err := a.dao.Create(ctx, a.toDao(bizId, definition))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *attributeDefinitionRepository) Delete(ctx context.Context, bizId, id int64) error {
	return a.dao.DeleteByBizIdAndID(ctx, bizId, id)
}

func (a *attributeDefinitionRepository) FindByBizID(ctx context.Context, bizId int64) (domain.BizAttrDefinition, error) {
	daoAttrs, err := a.dao.FindByBizID(ctx, bizId)
	if err != nil {
		return domain.BizAttrDefinition{}, err
	}
	bizDef := domain.BizAttrDefinition{
		BizID:   bizId,
		AllDefs: make(map[int64]domain.AttributeDefinition, len(daoAttrs)),
	}
	for _, daoAttr := range daoAttrs {
		domainAttr := a.toDomain(daoAttr)
		switch daoAttr.EntityType {
		case domain.ResourceTypeEntity.String():
			{
				bizDef.ResourceAttrDefs = append(bizDef.ResourceAttrDefs, domainAttr)
			}
		case domain.EnvironmentTypeEntity.String():
			{
				bizDef.EnvironmentAttrDefs = append(bizDef.EnvironmentAttrDefs, domainAttr)

			}
		case domain.SubjectTypeEntity.String():
			{
				bizDef.SubjectAttrDefs = append(bizDef.SubjectAttrDefs, domainAttr)
			}

		}
		bizDef.AllDefs[domainAttr.ID] = domainAttr
	}
	return bizDef, nil
}

func (a *attributeDefinitionRepository) FindByBizIdAndId(ctx context.Context, bizId, id int64) (domain.AttributeDefinition, error) {
	res, err := a.dao.FindByBizIdAndID(ctx, bizId, id)
	if err != nil {
		return domain.AttributeDefinition{}, err
	}
	return a.toDomain(res), nil
}

func NewAttributeDefinitionRepository(dao dao.AttributeDefinitionDAO) AttributeDefinitionRepository {
	return &attributeDefinitionRepository{dao: dao}
}

func (a *attributeDefinitionRepository) toDomain(definition dao.AttributeDefinition) domain.AttributeDefinition {
	return domain.AttributeDefinition{
		ID:             definition.ID,
		Name:           definition.Name,
		Description:    definition.Description,
		DataType:       domain.DataType(definition.DataType),
		EntityType:     domain.EntityType(definition.EntityType),
		ValidationRule: definition.ValidationRule,
		Ctime:          definition.Ctime,
		Utime:          definition.Utime,
	}
}
func (a *attributeDefinitionRepository) toDao(bizId int64, definition domain.AttributeDefinition) dao.AttributeDefinition {
	return dao.AttributeDefinition{
		BizID:          bizId,
		ID:             definition.ID,
		Name:           definition.Name,
		Description:    definition.Description,
		DataType:       definition.DataType.String(),
		EntityType:     definition.EntityType.String(),
		ValidationRule: definition.ValidationRule,
		Ctime:          definition.Ctime,
		Utime:          definition.Utime,
	}
}
