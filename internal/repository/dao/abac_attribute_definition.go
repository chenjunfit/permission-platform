package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"gorm.io/gorm/clause"
	"time"
)

type AttributeDefinition struct {
	ID             int64  `gorm:"column:id;primaryKey;;autoIncrement;"`
	BizID          int64  `gorm:"column:biz_id;uniqueIndex:idx_biz_id_name;comment:和name组成唯一索引，比如说代表订单组的biz_id"`
	Name           string `gorm:"column:name;size:100;not null;type:varchar(255);uniqueIndex:idx_biz_id_name;comment:属性名称"`
	Description    string `gorm:"column:description;type:text;comment:属性描述"`
	DataType       string `gorm:"column:data_type;type:varchar(255);not null;comment:属性数据类型"`
	EntityType     string `gorm:"column:entity_type;type:enum('subject','resource','environment');not null;comment:属性所属实体类型;index:idx_entity_type"`
	ValidationRule string `gorm:"column:validation_rule;comment:验证规则，正则表达式"`
	Ctime          int64  `gorm:"column:ctime;comment:创建时间"` // 使用毫秒级时间戳
	Utime          int64  `gorm:"column:utime;comment:更新时间"` // 使用毫秒级时间戳
}

func (AttributeDefinition) TableName() string {
	return "attribute_definitions"
}

type AttributeDefinitionDAO interface {
	Create(ctx context.Context, definition AttributeDefinition) (int64, error)
	FindByBizIdAndID(ctx context.Context, bizId, id int64) (AttributeDefinition, error)
	DeleteByBizIdAndID(ctx context.Context, bizId, id int64) error
	FindByBizID(ctx context.Context, bizId int64) ([]AttributeDefinition, error)
	FindByIDs(ctx context.Context, ids []int64) (map[int64]AttributeDefinition, error)
}

type attributeDefinitionDAO struct {
	db *egorm.Component
}

func NewAttributeDefinitionDAO(db *egorm.Component) AttributeDefinitionDAO {
	return &attributeDefinitionDAO{db: db}
}

func (a *attributeDefinitionDAO) Create(ctx context.Context, definition AttributeDefinition) (int64, error) {
	now := time.Now().UnixMilli()
	definition.Utime = now
	definition.Ctime = now
	err := a.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "biz_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "data_type", "entity_type", "validation_rule"}),
	}).Create(&definition).Error
	return definition.ID, err
}

func (a *attributeDefinitionDAO) FindByBizIdAndID(ctx context.Context, bizId, id int64) (AttributeDefinition, error) {
	var definition AttributeDefinition
	err := a.db.WithContext(ctx).Model(&AttributeDefinition{}).Where("biz_id=? AND id=?").First(&definition).Error
	return definition, err
}

func (a *attributeDefinitionDAO) DeleteByBizIdAndID(ctx context.Context, bizId, id int64) error {
	return a.db.WithContext(ctx).Model(&AttributeDefinition{}).Where("biz_id=? AND id=?", bizId, id).Delete(&AttributeDefinition{}).Error
}

func (a *attributeDefinitionDAO) FindByBizID(ctx context.Context, bizId int64) ([]AttributeDefinition, error) {
	var definitions []AttributeDefinition
	err := a.db.WithContext(ctx).Model(&AttributeDefinition{}).Where("biz_id=?", bizId).Find(&definitions).Error
	return definitions, err
}

func (a *attributeDefinitionDAO) FindByIDs(ctx context.Context, ids []int64) (map[int64]AttributeDefinition, error) {
	//TODO implement me
	panic("implement me")
}
