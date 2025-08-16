package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"time"
)

/*
1）唯一索引
- uk_biz_user_permission : 由 BizID (1) + UserID (2) + PermissionID (3) 组成的复合唯一索引，确保同一业务下用户和权限的关联关系不重复。 （2）普通索引
- idx_biz_user : BizID + UserID ，优化“查询某业务下某用户的所有权限”场景。
- idx_biz_permission : BizID + PermissionID ，优化“查询某业务下某权限关联的所有用户”场景。
- idx_biz_effect : BizID + Effect ，按权限类型（允许/拒绝）过滤。
- idx_biz_resource_type / idx_biz_action : 按资源类型和操作类型过滤权限。
- idx_time_range : BizID + StartTime + EndTime ，优化时间范围内的权限查询。
- idx_current_valid : BizID + Effect + StartTime + EndTime ，专门用于查询 当前有效的权限 （核心索引，提升权限检查性能）。
- idx_biz_resource_key_action : BizID + ResourceType + ResourceKey + PermissionAction ，支持“按资源和操作精确查询用户权限”的场景。
*/
type UserPermission struct {
	ID               int64  `gorm:"primaryKey;autoIncrement;comment:'用户权限关联关系ID'"`
	BizID            int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_permission,priority:1;index:idx_biz_user,priority:1;index:idx_biz_permission,priority:1;index:idx_biz_effect,priority:1;index:idx_biz_resource_type,priority:1;index:idx_biz_action,priority:1;index:idx_time_range,priority:1;index:idx_current_valid,priority:1;index:idx_biz_resource_key_action,priority:1;comment:'业务ID'"`
	UserID           int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_permission,priority:2;index:idx_biz_user,priority:2;comment:'用户ID'"`
	PermissionID     int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_permission,priority:3;index:idx_biz_permission,priority:2;comment:'权限ID'"`
	PermissionName   string `gorm:"type:VARCHAR(255);NOT NULL;comment:'权限名称（冗余字段，加速查询与展示）'"`
	ResourceType     string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_type,priority:2;index:idx_biz_resource_key_action,priority:2;comment:'资源类型（冗余字段，加速查询）'"`
	ResourceKey      string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_resource_key_action,priority:3;comment:'资源标识符（冗余字段，加速查询）'"`
	PermissionAction string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_action,priority:2;index:idx_biz_resource_key_action,priority:4;comment:'操作类型（冗余字段，加速查询）'"`
	StartTime        int64  `gorm:"NOT NULL;index:idx_time_range,priority:2;index:idx_current_valid,priority:3;comment:'权限生效时间'"`
	EndTime          int64  `gorm:"NOT NULL;index:idx_time_range,priority:3;index:idx_current_valid,priority:4;comment:'权限失效时间'"`
	Effect           string `gorm:"type:ENUM('allow', 'deny');NOT NULL;DEFAULT:'allow';index:idx_biz_effect,priority:2;index:idx_current_valid,priority:2;comment:'用于额外授予权限，或者取消权限，理论上不应该出现同时allow和deny，出现了就是deny优先于allow'"`
	Ctime            int64
	Utime            int64
}

func (UserPermission) TableName() string {
	return "user_permissions"
}

type UserPermissionDAO interface {
	Create(ctx context.Context, up UserPermission) (UserPermission, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]UserPermission, error)
	FindByBizIdAndUserId(ctx context.Context, bizId, userId int64) ([]UserPermission, error)
	FindByBizIdAndID(ctx context.Context, bizId, id int64) (UserPermission, error)
	DeleteBizIdAndId(ctx context.Context, bizId, id int64) error
	DeleteBizIdAndUserIdAndPermissionId(ctx context.Context, bizId, userId, permissionId int64) error
}
type userPermissionDao struct {
	db *egorm.Component
}

func (u *userPermissionDao) DeleteBizIdAndUserIdAndPermissionId(ctx context.Context, bizId, userId, permissionId int64) error {
	return u.db.WithContext(ctx).Model(&UserPermission{}).Where("biz_id=? AND user_id=? AND permission_id = ?", bizId, userId, permissionId).Delete(&userPermissionDao{}).Error
}

func (u *userPermissionDao) Create(ctx context.Context, up UserPermission) (UserPermission, error) {
	now := time.Now().Unix()
	up.Ctime = now
	up.Utime = now
	err := u.db.WithContext(ctx).Model(&UserPermission{}).Create(&up).Error
	return up, err
}

func (u *userPermissionDao) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]UserPermission, error) {
	ups := make([]UserPermission, 0)
	err := u.db.WithContext(ctx).Model(&UserPermission{}).Where("biz_id=?", bizId).Offset(offset).Limit(limit).Find(&ups).Error
	return ups, err
}

func (u *userPermissionDao) FindByBizIdAndUserId(ctx context.Context, bizId, userId int64) ([]UserPermission, error) {
	now := time.Now().Unix()
	ups := make([]UserPermission, 0)
	err := u.db.WithContext(ctx).Model(&UserPermission{}).Where("biz_id=? AND user_id=? AND start_time<=? AND end_time>=?", bizId, userId, now, now).Find(&ups).Error
	return ups, err
}

func (u *userPermissionDao) FindByBizIdAndID(ctx context.Context, bizId, id int64) (UserPermission, error) {
	up := UserPermission{}
	err := u.db.WithContext(ctx).Model(&UserPermission{}).Where("biz_id=? AND id=?", bizId, id).First(&up).Error
	return up, err
}

func (u *userPermissionDao) DeleteBizIdAndId(ctx context.Context, bizId, id int64) error {
	return u.db.WithContext(ctx).Model(&UserPermission{}).Where("biz_id=? AND id=?", bizId, id).Delete(&userPermissionDao{}).Error
}

func NewUserPermissionDAO(db *egorm.Component) UserPermissionDAO {
	return &userPermissionDao{db: db}
}
