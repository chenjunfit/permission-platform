package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"gorm.io/gorm/clause"
	"time"
)

// ResourceAttributeValue 资源属性值表模型
type ResourceAttributeValue struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID      int64  `gorm:"column:biz_id;uniqueIndex:idx_biz_resource_attr;comment:biz_id + resource_key + attr_id 唯一索引"`
	ResourceID int64  `gorm:"column:resource_id;not null;uniqueIndex:idx_biz_resource_attr;index:idx_resource_id;comment:资源ID"`
	AttrDefID  int64  `gorm:"column:attr_def_id;not null;uniqueIndex:idx_biz_resource_attr;index:idx_attr_def_id;comment:属性定义ID"`
	Value      string `gorm:"column:value;type:text;not null;comment:属性值，取决于 data_type"`
	Ctime      int64  `gorm:"column:ctime;"`
	Utime      int64  `gorm:"column:utime;"`
}

func (r ResourceAttributeValue) TableName() string {
	return "resource_attribute_values"
}

type ResourceAttributeValueDAO interface {
	Create(ctx context.Context, value ResourceAttributeValue) (int64, error)
	DeleteByID(ctx context.Context, id int64) error
	FindByBizIdAndResourceId(ctx context.Context, bizId, resourceId int64) ([]ResourceAttributeValue, error)
	FindByBizIdAndAttrId(ctx context.Context, bizId, attrId int64) ([]ResourceAttributeValue, error)
	FindByResourceId(ctx context.Context, resourceId []int64) (map[int64][]ResourceAttributeValue, error)
}
type resourceAttributeValueDao struct {
	db *egorm.Component
}

func NewResourceAttributeValueDAO(db *egorm.Component) ResourceAttributeValueDAO {
	return &resourceAttributeValueDao{db: db}
}
func (r *resourceAttributeValueDao) Create(ctx context.Context, value ResourceAttributeValue) (int64, error) {
	now := time.Now().UnixMilli()
	value.Utime = now
	value.Ctime = now
	err := r.db.WithContext(ctx).Model(&ResourceAttributeValue{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "biz_id"}, {Name: "resource_id"}, {Name: "attr_def_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"utime", "value"}),
	}).Create(&value).Error
	return value.ID, err
}

func (r *resourceAttributeValueDao) DeleteByID(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&ResourceAttributeValue{}).Where("id=?", id).Delete(&ResourceAttributeValue{}).Error
}

func (r *resourceAttributeValueDao) FindByBizIdAndResourceId(ctx context.Context, bizId, resourceId int64) ([]ResourceAttributeValue, error) {
	var ravs []ResourceAttributeValue
	err := r.db.WithContext(ctx).Model(&ResourceAttributeValue{}).Where("biz_id=? AND resource_id = ?", bizId, resourceId).Find(&ravs).Error
	return ravs, err
}

func (r *resourceAttributeValueDao) FindByBizIdAndAttrId(ctx context.Context, bizId, attrId int64) ([]ResourceAttributeValue, error) {
	var ravs []ResourceAttributeValue
	err := r.db.WithContext(ctx).Model(&ResourceAttributeValue{}).Where("biz_id=? AND attr_def_id = ? ", bizId, attrId).Find(&ravs).Error
	return ravs, err
}

func (r *resourceAttributeValueDao) FindByResourceId(ctx context.Context, resourceId []int64) (map[int64][]ResourceAttributeValue, error) {
	var ravs []ResourceAttributeValue
	err := r.db.WithContext(ctx).Model(&ResourceAttributeValue{}).Where("resource_id IN ?", resourceId).Find(&ravs).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]ResourceAttributeValue)
	for _, rav := range ravs {
		result[rav.ResourceID] = append(result[rav.ResourceID], rav)
	}
	return result, nil
}
