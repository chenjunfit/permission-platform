package abac

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
)

type AttributeValueSvc interface {
	SaveSubjectValue(ctx context.Context, bizID, subjectID int64, val domain.AttributeValue) (int64, error)
	DeleteSubjectValue(ctx context.Context, bizID, id int64) error
	FindSubjectValue(ctx context.Context, bizID, subjectID int64) (domain.ABACObject, error)
	FindSubjectValueWithDefinition(ctx context.Context, bizID, subjectID int64) (domain.ABACObject, error)

	SaveResourceValue(ctx context.Context, bizID, resourceID int64, val domain.AttributeValue) (int64, error)
	DeleteResourceValue(ctx context.Context, bizID, id int64) error
	FindResourceValue(ctx context.Context, bizID, resourceID int64) (domain.ABACObject, error)
	FindResourceValueWithDefinition(ctx context.Context, bizID, resourceID int64) (domain.ABACObject, error)

	SaveEnvironmentValue(ctx context.Context, bizID int64, val domain.AttributeValue) (int64, error)
	DeleteEnvironmentValue(ctx context.Context, bizID, id int64) error
	FindEnvironmentValue(ctx context.Context, bizID int64) (domain.ABACObject, error)
	FindEnvironmentValueWithDefinition(ctx context.Context, bizID int64) (domain.ABACObject, error)
}
type attributeValueSvc struct {
	repository.AttributeValueRepository
}

func NewAttributeValueSvc(repository repository.AttributeValueRepository) AttributeValueSvc {
	return &attributeValueSvc{repository}
}
