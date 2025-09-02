package repository

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
	"regexp"
)

type AttributeValueRepository interface {
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

type attributeValueRepository struct {
	resourceAttrDao   dao.ResourceAttributeValueDAO
	envAttrDao        dao.EnvironmentAttributeValueDAO
	subjectAttrDao    dao.SubjectAttributeValueDAO
	attrDefinitionDao dao.AttributeDefinitionDAO
}

func (a *attributeValueRepository) SaveEnvironmentValue(ctx context.Context, bizID int64, val domain.AttributeValue) (int64, error) {
	err := a.checkVal(ctx, bizID, val.AttrDef.ID, val.Value)
	if err != nil {
		return 0, err
	}
	daoVal := dao.EnvironmentAttributeValue{
		ID:        val.ID,
		BizID:     bizID,
		AttrDefID: val.AttrDef.ID,
		Value:     val.Value,
	}
	id, err := a.envAttrDao.Create(ctx, daoVal)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *attributeValueRepository) DeleteEnvironmentValue(ctx context.Context, bizID, id int64) error {
	return a.envAttrDao.DeleteByID(ctx, id)
}

func (a *attributeValueRepository) FindEnvironmentValue(ctx context.Context, bizID int64) (domain.ABACObject, error) {
	//to-do 使用缓存
	return a.FindEnvironmentValue(ctx, bizID)
}

func (a *attributeValueRepository) FindEnvironmentValueWithDefinition(ctx context.Context, bizID int64) (domain.ABACObject, error) {
	values, err := a.envAttrDao.FindByBizID(ctx, bizID)
	if err != nil {
		return domain.ABACObject{}, err
	}
	definitionIds := slice.Map(values, func(_ int, src dao.EnvironmentAttributeValue) int64 {
		return src.AttrDefID
	})
	definitionMap, err := a.attrDefinitionDao.FindByIDs(ctx, definitionIds)
	if err != nil {
		return domain.ABACObject{}, err
	}
	return domain.ABACObject{
		BizId: bizID,
		AttrValues: slice.Map(values, func(_ int, src dao.EnvironmentAttributeValue) domain.AttributeValue {
			return a.ToDomainEnvValue(src, definitionMap[src.AttrDefID])
		}),
	}, nil
}

func (a *attributeValueRepository) SaveResourceValue(ctx context.Context, bizID, resourceID int64, val domain.AttributeValue) (int64, error) {
	err := a.checkVal(ctx, bizID, val.AttrDef.ID, val.Value)
	if err != nil {
		return 0, err
	}
	daoVal := dao.ResourceAttributeValue{
		ID:         val.ID,
		BizID:      bizID,
		ResourceID: resourceID,
		AttrDefID:  val.AttrDef.ID,
		Value:      val.Value,
	}
	id, err := a.resourceAttrDao.Create(ctx, daoVal)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *attributeValueRepository) DeleteResourceValue(ctx context.Context, bizID, id int64) error {
	return a.resourceAttrDao.DeleteByID(ctx, id)
}

func (a *attributeValueRepository) FindResourceValue(ctx context.Context, bizID, resourceID int64) (domain.ABACObject, error) {
	//三级缓存 -todo
	return domain.ABACObject{}, nil
}

func (a *attributeValueRepository) FindResourceValueWithDefinition(ctx context.Context, bizID, resourceID int64) (domain.ABACObject, error) {
	values, err := a.resourceAttrDao.FindByBizIdAndResourceId(ctx, bizID, resourceID)
	if err != nil {
		return domain.ABACObject{}, err
	}
	definitionIds := slice.Map(values, func(_ int, src dao.ResourceAttributeValue) int64 {
		return src.AttrDefID
	})
	definitionMap, err := a.attrDefinitionDao.FindByIDs(ctx, definitionIds)
	if err != nil {
		return domain.ABACObject{}, err
	}
	result := domain.ABACObject{
		ID:    resourceID,
		BizId: bizID,
		AttrValues: slice.Map(values, func(_ int, src dao.ResourceAttributeValue) domain.AttributeValue {
			return a.ToDomainResourceValue(src, definitionMap[src.AttrDefID])
		}),
	}
	return result, nil
}

func (a *attributeValueRepository) matchRex(pattern, input string) error {
	matched, err := regexp.MatchString(pattern, input)
	if err != nil {
		return fmt.Errorf("正则表达式语法错误: ", err)
	}
	if !matched {
		return fmt.Errorf("填写的值不符合规范: ", err)
	}
	return nil
}
func (a *attributeValueRepository) checkVal(ctx context.Context, bizId, definitionId int64, value string) error {
	definitionDao, err := a.attrDefinitionDao.FindByBizIdAndID(ctx, bizId, definitionId)
	if err != nil {
		return err
	}
	return a.matchRex(definitionDao.ValidationRule, value)

}

func (a *attributeValueRepository) SaveSubjectValue(ctx context.Context, bizID, subjectID int64, val domain.AttributeValue) (int64, error) {
	err := a.checkVal(ctx, bizID, val.AttrDef.ID, val.Value)
	if err != nil {
		return 0, err
	}
	id, err := a.subjectAttrDao.Create(ctx, dao.SubjectAttributeValue{
		BizID:     bizID,
		SubjectID: subjectID,
		AttrDefID: val.AttrDef.ID,
		Value:     val.Value,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *attributeValueRepository) DeleteSubjectValue(ctx context.Context, bizID, id int64) error {
	return a.subjectAttrDao.DeleteByID(ctx, id)
}

func (a *attributeValueRepository) FindSubjectValue(ctx context.Context, bizID, subjectID int64) (domain.ABACObject, error) {
	//使用缓存- todo
	return domain.ABACObject{}, nil
}

func (a *attributeValueRepository) ToDomainAttributeValue(value dao.SubjectAttributeValue, definition dao.AttributeDefinition) domain.AttributeValue {
	return domain.AttributeValue{
		ID: value.ID,
		AttrDef: domain.AttributeDefinition{
			ID:             definition.ID,
			Name:           definition.Name,
			Description:    definition.Description,
			DataType:       domain.DataType(definition.DataType),
			EntityType:     domain.EntityType(definition.EntityType),
			ValidationRule: definition.ValidationRule,
			Ctime:          definition.Ctime,
			Utime:          definition.Utime,
		},
		Value: value.Value,
		Ctime: value.Ctime,
		Utime: value.Utime,
	}
}
func (a *attributeValueRepository) ToDomainResourceValue(value dao.ResourceAttributeValue, definition dao.AttributeDefinition) domain.AttributeValue {
	return domain.AttributeValue{
		ID: value.ID,
		AttrDef: domain.AttributeDefinition{
			ID:             definition.ID,
			Name:           definition.Name,
			Description:    definition.Description,
			DataType:       domain.DataType(definition.DataType),
			EntityType:     domain.EntityType(definition.EntityType),
			ValidationRule: definition.ValidationRule,
			Ctime:          definition.Ctime,
			Utime:          definition.Utime,
		},
		Value: value.Value,
		Ctime: value.Ctime,
		Utime: value.Utime,
	}
}
func (a *attributeValueRepository) ToDomainEnvValue(value dao.EnvironmentAttributeValue, definition dao.AttributeDefinition) domain.AttributeValue {
	return domain.AttributeValue{
		ID: value.ID,
		AttrDef: domain.AttributeDefinition{
			ID:             definition.ID,
			Name:           definition.Name,
			Description:    definition.Description,
			DataType:       domain.DataType(definition.DataType),
			EntityType:     domain.EntityType(definition.EntityType),
			ValidationRule: definition.ValidationRule,
			Ctime:          definition.Ctime,
			Utime:          definition.Utime,
		},
		Value: value.Value,
		Ctime: value.Ctime,
		Utime: value.Utime,
	}
}

func (a *attributeValueRepository) FindSubjectValueWithDefinition(ctx context.Context, bizID, subjectID int64) (domain.ABACObject, error) {
	//不使用缓存,直接查数据库
	subjectValueDao, err := a.subjectAttrDao.FindByBizIdAndSubjectID(ctx, bizID, subjectID)
	if err != nil {
		return domain.ABACObject{}, err
	}
	definitionIds := slice.Map(subjectValueDao, func(idx int, src dao.SubjectAttributeValue) int64 {
		return src.AttrDefID
	})
	definitionMap, err := a.attrDefinitionDao.FindByIDs(ctx, definitionIds)
	if err != nil {
		return domain.ABACObject{}, err
	}
	abacObject := domain.ABACObject{
		ID:    subjectID,
		BizId: bizID,
		AttrValues: slice.Map(subjectValueDao, func(idx int, src dao.SubjectAttributeValue) domain.AttributeValue {
			return a.ToDomainAttributeValue(src, definitionMap[src.AttrDefID])
		}),
	}
	return abacObject, nil
}

func NewAttributeValueRepository(
	resourceAttrDao dao.ResourceAttributeValueDAO,
	envAttrDao dao.EnvironmentAttributeValueDAO,
	subjectAttrDao dao.SubjectAttributeValueDAO,
	attrDefinitionDao dao.AttributeDefinitionDAO,
) AttributeValueRepository {
	return &attributeValueRepository{
		resourceAttrDao:   resourceAttrDao,
		envAttrDao:        envAttrDao,
		subjectAttrDao:    subjectAttrDao,
		attrDefinitionDao: attrDefinitionDao,
	}
}
