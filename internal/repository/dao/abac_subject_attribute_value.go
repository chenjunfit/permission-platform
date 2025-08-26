package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"gorm.io/gorm/clause"
	"time"
)

/*
- 复合唯一索引 idx_biz_subject_attr ：由 biz_id 、 subject_id 和 attr_def_id 组成，确保每个主体在特定业务下对同一属性只有一个值
- 单字段索引 idx_subject_id ：加速按主体 ID 查询
- 单字段索引 idx_attr_def_id ：加速按属性定义 ID 查询
*/
type SubjectAttributeValue struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID     int64  `gorm:"column:biz_id;uniqueIndex:idx_biz_subject_attr;comment:biz_id + subject_id + attr_id 唯一索引"`
	SubjectID int64  `gorm:"column:subject_id;not null;uniqueIndex:idx_biz_subject_attr;index:idx_subject_id;comment:主体ID，通常是用户ID"`
	AttrDefID int64  `gorm:"column:attr_def_id;not null;uniqueIndex:idx_biz_subject_attr;index:idx_attr_def_id;comment:属性定义ID"`
	Value     string `gorm:"column:value;type:text;not null;comment:属性值，取决于 data_type"`
	Ctime     int64  `gorm:"column:ctime;"`
	Utime     int64  `gorm:"column:utime;"`
}

func (s SubjectAttributeValue) TableName() string {
	return "subject_attribute_values"
}

type SubjectAttributeValueDAO interface {
	Create(ctx context.Context, value SubjectAttributeValue) (int64, error)
	FindByID(ctx context.Context, id int64) (SubjectAttributeValue, error)
	DeleteByID(ctx context.Context, id int64) error
	FindByBizIdAndSubjectID(ctx context.Context, bizId, subjectId int64) ([]SubjectAttributeValue, error)
}

type subjectAttributeValueDao struct {
	db *egorm.Component
}

func NewSubjectAttributeValueDAO(db *egorm.Component) SubjectAttributeValueDAO {
	return &subjectAttributeValueDao{db: db}
}
func (s *subjectAttributeValueDao) Create(ctx context.Context, value SubjectAttributeValue) (int64, error) {
	now := time.Now().UnixMilli()
	value.Ctime = now
	value.Utime = now
	err := s.db.WithContext(ctx).Model(&SubjectAttributeValue{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "biz_id"}, {Name: "subject_id"}, {Name: "attr_def_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "utime"}),
	}).Create(&value).Error
	return value.ID, err
}

func (s *subjectAttributeValueDao) FindByID(ctx context.Context, id int64) (SubjectAttributeValue, error) {
	var subjectAttributeValue SubjectAttributeValue
	err := s.db.WithContext(ctx).Model(&SubjectAttributeValue{}).Where("id = ?", id).First(&subjectAttributeValue).Error
	return subjectAttributeValue, err
}

func (s *subjectAttributeValueDao) DeleteByID(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Model(&SubjectAttributeValue{}).Where("id = ?", id).Delete(&SubjectAttributeValue{}).Error
}

func (s *subjectAttributeValueDao) FindByBizIdAndSubjectID(ctx context.Context, bizId, subjectId int64) ([]SubjectAttributeValue, error) {
	var subjectAttributeValues []SubjectAttributeValue
	err := s.db.WithContext(ctx).Model(&SubjectAttributeValue{}).Where("biz_id=? AND subject_id=?", bizId, subjectId).Find(&subjectAttributeValues).Error
	return subjectAttributeValues, err
}
