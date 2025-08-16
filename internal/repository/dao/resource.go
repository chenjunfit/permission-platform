package dao

import (
	"context"
	"fmt"
	"github.com/ego-component/egorm"
	"github.com/permission-dev/internal/errs"
	"time"
)

type Resource struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:资源ID'"`
	BizID       int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_type_key,priority:1;index:idx_biz_type,priority:1;index:idx_biz_key,priority:1;comment:'业务ID'"`
	Type        string `gorm:"type:VARCHAR(100);NOT NULL;uniqueIndex:uk_biz_type_key,priority:2;index:idx_biz_type,priority:2;comment:'资源类型，被冗余，创建后不允许修改'"`
	Key         string `gorm:"type:VARCHAR(255);NOT NULL;uniqueIndex:uk_biz_type_key,priority:3;index:idx_biz_key,priority:2;comment:'资源业务标识符 (如 用户ID, 文档路径)，被冗余，创建后不允许修改'"`
	Name        string `gorm:"type:VARCHAR(255);NOT NULL;comment:'资源名称'"`
	Description string `gorm:"type:TEXT;comment:'资源描述'"`
	Metadata    string `gorm:"type:TEXT;comment:'资源元数据'"`
	Ctime       int64
	Utime       int64
}

func (r Resource) TableName() string {
	return "resources"
}

type ResourceDao interface {
	Create(ctx context.Context, resource Resource) (Resource, error)
	UpdateByBizIDAndID(ctx context.Context, resource Resource) error
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error

	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]Resource, error)
	FindByBizIDAndID(ctx context.Context, bizId, id int64) (Resource, error)
	FindByBizIDAndTypeAndKey(ctx context.Context, bizId int64, resourceType, resourceKey string) (Resource, error)
}

type resourceDao struct {
	db *egorm.Component
}

func NewResourceDao(db *egorm.Component) ResourceDao {
	return &resourceDao{db: db}
}
func (r *resourceDao) Create(ctx context.Context, resource Resource) (Resource, error) {
	now := time.Now().Unix()
	resource.Ctime = now
	resource.Utime = now
	err := r.db.WithContext(ctx).Create(&resource).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return Resource{}, fmt.Errorf("%w", errs.ErrResourceDuplicate)
		}
		return Resource{}, err
	}
	return resource, nil
}

func (r *resourceDao) UpdateByBizIDAndID(ctx context.Context, resource Resource) error {
	now := time.Now().Unix()
	resource.Utime = now
	err := r.db.WithContext(ctx).
		Model(&Resource{}).
		Where("biz_id=? AND id=?", resource.BizID, resource.ID).
		Updates(map[string]interface{}{
			"name":        resource.Name,
			"description": resource.Description,
			"metadata":    resource.Metadata,
			"utime":       resource.Utime,
		}).Error
	return err
}

func (r *resourceDao) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return r.db.WithContext(ctx).Where("biz_id=? AND id=?", bizId, id).Delete(&Resource{}).Error
}

func (r *resourceDao) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]Resource, error) {
	resources := make([]Resource, 0)
	err := r.db.WithContext(ctx).Where("biz_id=?", bizId).Offset(offset).Limit(limit).Find(&resources).Error
	return resources, err
}

func (r *resourceDao) FindByBizIDAndID(ctx context.Context, bizId, id int64) (Resource, error) {
	resource := Resource{}
	err := r.db.WithContext(ctx).Where("biz_id=? AND id=?", bizId, id).First(&resource).Error
	return resource, err
}

func (r resourceDao) FindByBizIDAndTypeAndKey(ctx context.Context, bizId int64, resourceType, resourceKey string) (Resource, error) {
	resource := Resource{}
	err := r.db.WithContext(ctx).Where("biz_id=? AND type = ? AND key=?", bizId, resourceType, resourceKey).First(&resource).Error
	return resource, err
}
