package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"time"
)

/*
- 唯一约束 ： uk_biz_user_role (BizID, UserID, RoleID) 确保数据一致性
- 多维度索引 ：
- idx_biz_user (BizID, UserID)：优化按业务和用户查询角色的场景
- idx_biz_role (BizID, RoleID)：优化按业务和角色查询用户的场景
- idx_biz_user_role_validity (BizID, UserID, RoleType, StartTime, EndTime)：优化有效期相关的查询
*/
type UserRole struct {
	ID        int64  `json:"id,string" gorm:"primaryKey;autoIncrement;comment:用户角色关联关系主键'"`
	BizID     int64  `json:"biz_id,string" gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_role,priority:1;index:idx_biz_user,priority:1;index:idx_biz_role,priority:1;index:idx_biz_user_role_validity,priority:1;comment:'业务ID'"`
	UserID    int64  `json:"user_id,string" gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_role,priority:2;index:idx_biz_user,priority:2;index:idx_biz_user_role_validity,priority:2;comment:'用户ID'"`
	RoleID    int64  `json:"role_id,string" gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_user_role,priority:3;index:idx_biz_role,priority:2;comment:'角色ID'"`
	RoleName  string `json:"role_name" gorm:"type:VARCHAR(255);NOT NULL;comment:'角色名称（冗余字段，加速查询）'"`
	RoleType  string `json:"role_type" gorm:"type:VARCHAR(255);NOT NULL;index:idx_biz_user_role_validity,priority:3;comment:'角色类型（冗余字段，加速查询）'"`
	StartTime int64  `json:"start_time,string" gorm:"NOT NULL;index:idx_biz_user_role_validity,priority:4;comment:'授予角色生效时间'"`
	EndTime   int64  `json:"end_time,string" gorm:"NOT NULL;index:idx_biz_user_role_validity,priority:5;comment:'授予角色失效时间'"`
	Ctime     int64  `json:"ctime,string"`
	Utime     int64  `json:"utime,string"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type UserRoleDAO interface {
	Create(ctx context.Context, role UserRole) (UserRole, error)
	FindByBizID(ctx context.Context, bizId int64) ([]UserRole, error)
	FindByBizIDAndID(ctx context.Context, bizId, id int64) (UserRole, error)
	FindByBizIDAndUserID(ctx context.Context, bizId, userId int64) ([]UserRole, error)
	FindByBizIDAndRoleIDs(ctx context.Context, bizID, roleIds []int64) ([]UserRole, error)
	DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error
}

type userRoleDao struct {
	db *egorm.Component
}

func (u *userRoleDao) Create(ctx context.Context, role UserRole) (UserRole, error) {
	now := time.Now().Unix()
	role.Utime = now
	role.Ctime = now
	err := u.db.WithContext(ctx).Model(&UserRole{}).Create(&role).Error
	return role, err
}

func (u *userRoleDao) FindByBizID(ctx context.Context, bizId int64) ([]UserRole, error) {
	userRoles := make([]UserRole, 0)
	err := u.db.WithContext(ctx).Model(&UserRole{}).Where("biz_id=?", bizId).Find(&userRoles).Error
	return userRoles, err
}

func (u *userRoleDao) FindByBizIDAndID(ctx context.Context, bizId, id int64) (UserRole, error) {
	userRole := UserRole{}
	err := u.db.WithContext(ctx).Model(&UserRole{}).Where("biz_id=? AND id=?", bizId, id).First(&userRole).Error
	return userRole, err
}

func (u *userRoleDao) FindByBizIDAndUserID(ctx context.Context, bizId, userId int64) ([]UserRole, error) {
	userRoles := make([]UserRole, 0)
	err := u.db.WithContext(ctx).Model(&UserRole{}).Where("biz_id=? AND user_id=?", bizId, userId).Find(&userRoles).Error
	return userRoles, err
}

func (u *userRoleDao) FindByBizIDAndRoleIDs(ctx context.Context, bizID, roleIds []int64) ([]UserRole, error) {
	userRoles := make([]UserRole, 0)
	now := time.Now().Unix()
	err := u.db.WithContext(ctx).Model(&UserRole{}).Where("biz_id=? AND role_id in (?) AND start_time<=? AND end_time>=?", bizID, roleIds, now, now).Find(&userRoles).Error
	return userRoles, err
}

func (u *userRoleDao) DeleteByBizIDAndID(ctx context.Context, bizId, id int64) error {
	return u.db.WithContext(ctx).Model(&UserRole{}).Where("biz_id=? AND id=?", bizId, id).Delete(&UserRole{}).Error
}

func NewUserDaoDAO(db *egorm.Component) UserRoleDAO {
	return &userRoleDao{
		db,
	}
}
