package abac

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
)

type AttributeDefinitionSvc interface {
	Create(ctx context.Context, bizId int64, definition domain.AttributeDefinition) (int64, error)
	Delete(ctx context.Context, bizId, id int64) error
	//返回bizId下所有的属性定义(env,subject,resource)
	FindByBizID(ctx context.Context, bizId int64) (domain.BizAttrDefinition, error)
	FindByBizIdAndId(ctx context.Context, bizId, id int64) (domain.AttributeDefinition, error)
}

type attributeDefinitionSvc struct {
	repo repository.AttributeDefinitionRepository
}

func (a *attributeDefinitionSvc) Create(ctx context.Context, bizId int64, definition domain.AttributeDefinition) (int64, error) {
	return a.Create(ctx, bizId, definition)
}

func (a *attributeDefinitionSvc) Delete(ctx context.Context, bizId, id int64) error {
	return a.Delete(ctx, bizId, id)
}

func (a *attributeDefinitionSvc) FindByBizID(ctx context.Context, bizId int64) (domain.BizAttrDefinition, error) {
	return a.FindByBizID(ctx, bizId)
}

func (a *attributeDefinitionSvc) FindByBizIdAndId(ctx context.Context, bizId, id int64) (domain.AttributeDefinition, error) {
	return a.repo.FindByBizIdAndId(ctx, bizId, id)
}

func NewAttributeDefinitionSvc(repo repository.AttributeDefinitionRepository) AttributeDefinitionSvc {
	return &attributeDefinitionSvc{repo: repo}
}
