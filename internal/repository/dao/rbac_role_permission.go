package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"github.com/permission-dev/internal/errs"
	"time"
)

/*
（1）唯一索引
- uk_biz_role_permission : 由 BizID (1) + RoleID (2) + PermissionID (3) 组成的复合唯一索引，确保同一业务下角色和权限的关联关系不重复。
（2）普通索引
- idx_biz_role : BizID + RoleID ，优化“查询某业务下某角色的所有权限”场景。
- idx_biz_permission : BizID + PermissionID ，优化“查询某业务下某权限关联的所有角色”场景。
- idx_biz_role_type : BizID + RoleType ，按角色类型过滤权限。
- idx_biz_resource_type : BizID + ResourceType ，按资源类型过滤权限。
- idx_biz_action : BizID + PermissionAction ，按操作类型过滤权限。
- idx_biz_resource_key_action : BizID + ResourceType + ResourceKey + PermissionAction ，支持“按资源和操作精确查询权限”的场景。
*/
type RolePermission struct {
	ID               int64  `gorm:"primaryKey;autoIncrement;comment:'角色权限关联关系ID'"`
	BizID            int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_role_permission,priority:1;index:idx_biz_role,priority:1;index:idx_biz_permission,priority:1;index:idx_biz_role_type,priority:1;index:idx_biz_resource_type,priority:1;index:idx_biz_action,priority:1;index:idx_biz_resource_key_action,priority:1;comment:'业务ID'"`
	RoleID           int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_role_permission,priority:2;index:idx_biz_role,priority:2;comment:'角色ID'"`
	PermissionID     int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_role_permission,priority:3;index:idx_biz_permission,priority:2;comment:'权限ID'"`
	RoleName         string `gorm:"type:VARCHAR(255);NOT NULL;comment:'角色名称（冗余字段，加速查询）'"`
	RoleType         string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_role_type,priority:2;comment:'角色类型（冗余字段，加速查询）'"`
	ResourceType     string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_type,priority:2;index:idx_biz_resource_key_action,priority:2;comment:'资源类型（冗余字段，加速查询）'"`
	ResourceKey      string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_key_action,priority:3;comment:'资源标识符（冗余字段，加速查询）'"`
	PermissionAction string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_action,priority:2;index:idx_biz_resource_key_action,priority:4;comment:'操作类型（冗余字段，加速查询）'"`
	Ctime            int64
	Utime            int64
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type RolePermissionDAO interface {
	Create(ctx context.Context, rp RolePermission) (RolePermission, error)
	FindByBizID(ctx context.Context, bizId int64) ([]RolePermission, error)
	FindByBizIdAndID(ctx context.Context, bizId, id int64) (RolePermission, error)
	FindByBizIdAndResourceType(ctx context.Context, bizId, resourceType string, offset, limit int) ([]RolePermission, error)
	FindByBizIDAndRoleIds(ctx context.Context, bizId int64, roleIds []int64) ([]RolePermission, error)
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error
}
type rolePermissionDAO struct {
	db *egorm.Component
}

func (r *rolePermissionDAO) FindByBizIdAndResourceType(ctx context.Context, bizId, resourceType string, offset, limit int) ([]RolePermission, error) {
	rolePermissions := make([]RolePermission, 0)
	err := r.db.WithContext(ctx).Model(&RolePermission{}).Where("biz_id=? AND resource_type=?", bizId, resourceType).Offset(offset).Limit(limit).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *rolePermissionDAO) Create(ctx context.Context, rp RolePermission) (RolePermission, error) {
	now := time.Now().Unix()
	rp.Utime = now
	rp.Ctime = now
	err := r.db.WithContext(ctx).Model(&RolePermission{}).Create(&rp).Error
	if err != nil {
		if isUniqueConstraintError(err) {
			return RolePermission{}, errs.ErrRolePermissionDuplicate
		}
		return rp, err
	}
	return rp, nil
}

func (r *rolePermissionDAO) FindByBizID(ctx context.Context, bizId int64) ([]RolePermission, error) {
	rolePermissions := make([]RolePermission, 0)
	err := r.db.WithContext(ctx).Model(&RolePermission{}).Where("biz_id=?", bizId).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *rolePermissionDAO) FindByBizIdAndID(ctx context.Context, bizId, id int64) (RolePermission, error) {
	rolePermission := RolePermission{}
	err := r.db.WithContext(ctx).Model(&RolePermission{}).Where("biz_id=? AND id=?", bizId, id).First(&rolePermission).Error
	return rolePermission, err
}

func (r *rolePermissionDAO) FindByBizIDAndRoleIds(ctx context.Context, bizId int64, roleIds []int64) ([]RolePermission, error) {
	rolePermissions := make([]RolePermission, 0)
	err := r.db.WithContext(ctx).Model(&RolePermission{}).Where("biz_id=? AND role_id in (?)", bizId, roleIds).Find(&rolePermissions).Error
	return rolePermissions, err
}

func (r *rolePermissionDAO) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return r.db.WithContext(ctx).Model(&RolePermission{}).Where("biz_id=? AND id=?", bizId, id).Delete(&RolePermission{}).Error
}

func NewRolePermissionDAO(db *egorm.Component) RolePermissionDAO {
	return &rolePermissionDAO{db: db}
}
