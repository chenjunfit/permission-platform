package dao

import (
	"context"
	"fmt"
	"github.com/ego-component/egorm"
	"github.com/permission-dev/internal/errs"
	"time"
)

/*
### 索引设计说明
- uk_biz_resource_action : 复合唯一索引 (BizID, ResourceID, Action)，确保同一业务下对同一资源的同一操作不会重复定义权限
- idx_biz_action : 复合索引 (BizID, Action)，加速按业务和操作类型查询
- idx_biz_resource_type : 复合索引 (BizID, ResourceType)，加速按业务和资源类型查询
- idx_biz_resource_key : 复合索引 (BizID, ResourceKey)，加速按业务和资源标识符查询
- idx_resource_id : 单字段索引，加速按资源 ID 查询
*/
type Permission struct {
	ID           int64  `gorm:"primaryKey;autoIncrement;comment:'权限ID'"`
	BizID        int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_resource_action,priority:1;index:idx_biz_action,priority:1;index:idx_biz_resource_type,priority:1;index:idx_biz_resource_key,priority:1;comment:'业务ID'"`
	Name         string `gorm:"type:VARCHAR(255);NOT NULL;comment:'权限名称'"`
	Description  string `gorm:"type:TEXT;comment:'权限描述'"`
	ResourceID   int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_resource_action,priority:2;index:idx_resource_id;comment:'关联的资源ID，创建后不可修改'"`
	ResourceType string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_type,priority:2;comment:'资源类型，冗余字段，加速查询'"`
	ResourceKey  string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_key,priority:2;comment:'资源业务标识符 (如 用户ID, 文档路径)，冗余字段，加速查询'"`
	Action       string `gorm:"type:VARCHAR(255);NOT NULL;NOT NULL;uniqueIndex:uk_biz_resource_action,priority:3;index:idx_biz_action,priority:2;comment:'操作类型'"`
	Metadata     string `gorm:"type:TEXT;comment:'权限元数据，可扩展字段'"`
	Ctime        int64
	Utime        int64
}

func (Permission) TableName() string {
	return "permissions"
}

type PermissionDAO interface {
	Create(ctx context.Context, permission Permission) (Permission, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]Permission, error)
	FindByBizIDAndID(ctx context.Context, bizId, id int64) (Permission, error)
	UpdateByBizIDAndID(ctx context.Context, permission Permission) error
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error
	FindPermissions(ctx context.Context, bizId int64, resourceType, resourceKey string, action []string) ([]Permission, error)
}
type permissionDao struct {
	db *egorm.Component
}

func NewPermissionDAO(db *egorm.Component) PermissionDAO {
	return &permissionDao{
		db: db,
	}
}
func (p *permissionDao) Create(ctx context.Context, permission Permission) (Permission, error) {
	now := time.Now().Unix()
	permission.Utime = now
	permission.Ctime = now
	err := p.db.WithContext(ctx).Create(&permission).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return Permission{}, fmt.Errorf("%w", errs.ErrRoleDuplicate)
		}
		return Permission{}, err
	}
	return permission, nil
}

func (p *permissionDao) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]Permission, error) {
	permissions := make([]Permission, 0)
	err := p.db.WithContext(ctx).Model(&Permission{}).Where("biz_id=?", bizId).Offset(offset).Limit(limit).Find(&permissions).Error
	return permissions, err
}

func (p *permissionDao) FindByBizIDAndID(ctx context.Context, bizId, id int64) (Permission, error) {
	permission := Permission{}
	err := p.db.WithContext(ctx).Model(&Permission{}).Where("biz_id=? AND id=?", bizId, id).First(&permission).Error
	return permission, err
}

func (p *permissionDao) UpdateByBizIDAndID(ctx context.Context, permission Permission) error {
	now := time.Now().Unix()
	permission.Utime = now
	return p.db.Where(ctx).Model(&Permission{}).Where("biz_id=? AND id=?", permission.BizID, permission.ID).
		Updates(map[string]interface{}{
			"description": permission.Description,
			"metadata":    permission.Metadata,
			"name":        permission.Name,
			"utime":       permission.Utime,
			"action":      permission.Action,
		}).Error

}

func (p *permissionDao) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return p.db.WithContext(ctx).Model(&Permission{}).Where("biz_id=? AND id=?", bizId, id).Delete(&Permission{}).Error
}

func (p *permissionDao) FindPermissions(ctx context.Context, bizId int64, resourceType, resourceKey string, action []string) ([]Permission, error) {
	permissions := make([]Permission, 0)
	err := p.db.WithContext(ctx).
		Model(&Permission{}).
		Where("biz_id=? AND resource_type=? AND resource_key=? AND action IN ?", bizId, resourceType, resourceKey, action).
		Find(&permissions).
		Error
	return permissions, err
}
