package dao

import (
	"context"
	"fmt"
	"github.com/ego-component/egorm"
	"github.com/permission-dev/internal/errs"
	"time"
)

type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:'角色ID'"`
	BizID       int64  `gorm:"type:BIGINT;NOT NULL;index:idx_biz_id;uniqueIndex:uk_biz_type_name,priority:1;comment:'业务ID'"`
	Type        string `gorm:"type:VARCHAR(255);NOT NULL;index:idx_role_type;uniqueIndex:uk_biz_type_name,priority:2;comment:'角色类（被冗余，创建后不可修改）'"`
	Name        string `gorm:"type:VARCHAR(255);NOT NULL;uniqueIndex:uk_biz_type_name,priority:3;comment:'角色名称（被冗余，创建后不可修改）'"`
	Description string `gorm:"type:TEXT;comment:'角色描述'"`
	Metadata    string `gorm:"type:TEXT;comment:'角色元数据，可扩展字段'"`
	Ctime       int64  `gorm:"DEFAULT NULL"`
	Utime       int64  `gorm:"DEFAULT NULL"`
}

func (Role) TableName() string {
	return "roles"
}

type RoleDAO interface {
	Create(ctx context.Context, role Role) (Role, error)
	FindByBizID(ctx context.Context, bizID int64, offset, limit int) ([]Role, error)
	FindByBizIDAndID(ctx context.Context, bizID, Id int64) (Role, error)
	FindByBizIDAndType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]Role, error)
	UpdateByBizIDAndID(ctx context.Context, role Role) error
	DeleteByBizIDAndID(ctx context.Context, bizID, Id int64) error
}

type roleDao struct {
	db *egorm.Component
}

func NewRoleDao(db *egorm.Component) RoleDAO {
	return &roleDao{db: db}
}

func (r *roleDao) Create(ctx context.Context, role Role) (Role, error) {
	now := time.Now().UnixMilli()
	role.Ctime = now
	role.Utime = now
	err := r.db.WithContext(ctx).Create(&role).Error
	if isUniqueConstraintError(err) {
		return Role{}, fmt.Errorf("%w", errs.ErrRoleDuplicate)
	}
	return role, nil
}

func (r *roleDao) FindByBizID(ctx context.Context, bizID int64, offset, limit int) ([]Role, error) {
	var roles []Role
	err := r.db.WithContext(ctx).Where("biz_id=?", bizID).Offset(offset).Limit(limit).Find(&roles).Error
	return roles, err
}

func (r *roleDao) FindByBizIDAndID(ctx context.Context, bizID, Id int64) (Role, error) {
	var role Role
	err := r.db.WithContext(ctx).Where("biz_id = ? AND id = ?", bizID, Id).First(&role).Error
	return role, err
}

func (r *roleDao) FindByBizIDAndType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]Role, error) {
	var roles []Role
	err := r.db.WithContext(ctx).Where("biz_id=? AND type = ?", bizID, roleType).Offset(offset).Limit(limit).Find(&roles).Error
	return roles, err
}

func (r *roleDao) UpdateByBizIDAndID(ctx context.Context, role Role) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&Role{}).
		Where("biz_id = ? AND id = ?", role.BizID, role.ID).
		Updates(map[string]interface{}{
			"description": role.Description,
			"metadata":    role.Metadata,
			"utime":       now,
		}).Error
}

func (r *roleDao) DeleteByBizIDAndID(ctx context.Context, bizID, Id int64) error {
	return r.db.WithContext(ctx).Where("biz_id = ? AND id = ?", bizID, Id).Delete(&Role{}).Error
}
