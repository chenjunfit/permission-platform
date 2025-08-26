package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"gorm.io/gorm/clause"
	"time"
)

// EnvironmentAttributeValue 环境属性表模型
type EnvironmentAttributeValue struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID     int64  `gorm:"column:biz_id;uniqueIndex:idx_biz_attribute;comment:业务ID"`
	AttrDefID int64  `gorm:"column:attr_def_id;not null;uniqueIndex:idx_biz_attribute;comment:属性定义ID"`
	Value     string `gorm:"column:value;type:text;comment:属性值，取决于 data_type"`
	Ctime     int64  `gorm:"column:ctime;comment:创建时间"`
	Utime     int64  `gorm:"column:utime;comment:更新时间"`
}

type EnvironmentAttributeValueDAO interface {
	Create(ctx context.Context, value EnvironmentAttributeValue) (int64, error)
	FindByID(ctx context.Context, id int64) (EnvironmentAttributeValue, error)
	FindByBizIDAndAttrId(ctx context.Context, bizId, attrId int64) (EnvironmentAttributeValue, error)
	FindByBizID(ctx context.Context, bizId int64) ([]EnvironmentAttributeValue, error)
	DeleteByID(ctx context.Context, id int64) error
}

type environmentAttributeValueDao struct {
	db *egorm.Component
}

func (e *environmentAttributeValueDao) Create(ctx context.Context, value EnvironmentAttributeValue) (int64, error) {
	now := time.Now().UnixMilli()
	value.Ctime = now
	value.Utime = now
	err := e.db.WithContext(ctx).Model(&EnvironmentAttributeValue{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "biz_id"}, {Name: "attr_def_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "utime"}),
	}).Create(&value).Error
	return value.ID, err
}

func (e *environmentAttributeValueDao) FindByID(ctx context.Context, id int64) (EnvironmentAttributeValue, error) {
	var environmentAttributeValue EnvironmentAttributeValue
	err := e.db.WithContext(ctx).Model(&EnvironmentAttributeValue{}).Where("id=?", id).First(&environmentAttributeValue).Error
	return environmentAttributeValue, err
}

func (e *environmentAttributeValueDao) FindByBizIDAndAttrId(ctx context.Context, bizId, attrId int64) (EnvironmentAttributeValue, error) {
	var environmentAttributeValue EnvironmentAttributeValue
	err := e.db.WithContext(ctx).Model(&EnvironmentAttributeValue{}).Where("biz_id=? AND attr_def_id=?", bizId, attrId).First(&environmentAttributeValue).Error
	return environmentAttributeValue, err
}

func (e *environmentAttributeValueDao) FindByBizID(ctx context.Context, bizId int64) ([]EnvironmentAttributeValue, error) {
	var attrs []EnvironmentAttributeValue
	err := e.db.WithContext(ctx).Model(&EnvironmentAttributeValue{}).Where("biz_id=?", bizId).Find(&attrs).Error
	return attrs, err
}

func (e *environmentAttributeValueDao) DeleteByID(ctx context.Context, id int64) error {
	return e.db.WithContext(ctx).
		Model(&EnvironmentAttributeValue{}).
		Where("id = ?", id).
		Delete(&EnvironmentAttributeValue{}).Error
}

func NewEnvironmentAttributeValueDAO(db *egorm.Component) EnvironmentAttributeValueDAO {
	return &environmentAttributeValueDao{db: db}
}
