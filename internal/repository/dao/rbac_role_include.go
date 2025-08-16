package dao

import (
	"context"
	"github.com/ego-component/egorm"
	"time"
)

/*
- 唯一索引“uk_biz_including_included“ : 由“BizID“ ,“IncludingRoleID“ ,“IncludedRoleID“ 组成，确保在同一业务下，角色包含关系唯一
- 普通索引“idx_biz_including_role“ : 由“BizID“ ,“IncludingRoleID“ 组成，用于加速通过业务ID和包含者角色ID查询
- 普通索引“idx_biz_included_role“ : 由“BizID“ ,“IncludedRoleID“ 组成，用于加速通过业务ID和被包含角色ID
*/
type RoleInclusion struct {
	ID                int64  `gorm:"primaryKey;autoIncrement;comment:角色包含关系ID'"`
	BizID             int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_including_included,priority:1;index:idx_biz_including_role,priority:1;index:idx_biz_included_role,priority:1;comment:'业务ID'"`
	IncludingRoleID   int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_including_included,priority:2;index:idx_biz_including_role,priority:2;comment:'包含者角色ID（拥有其他角色权限）'"`
	IncludingRoleType string `gorm:"type:VARCHAR(255);NOT NULL;comment:'包含者角色类型（冗余字段，加速查询）'"`
	IncludingRoleName string `gorm:"type:VARCHAR(255);NOT NULL;comment:'包含者角色名称（冗余字段，加速查询）'"`
	IncludedRoleID    int64  `gorm:"type:BIGINT;NOT NULL;uniqueIndex:uk_biz_including_included,priority:3;index:idx_biz_included_role,priority:2;comment:'被包含角色ID（权限被包含）'"`
	IncludedRoleType  string `gorm:"type:VARCHAR(255);NOT NULL;comment:'被包含角色类型（冗余字段，加速查询）'"`
	IncludedRoleName  string `gorm:"type:VARCHAR(255);NOT NULL;comment:'被包含角色名称（冗余字段，加速查询）'"`
	Ctime             int64
	Utime             int64
}

func (RoleInclusion) TableName() string {
	return "role_inclusions"
}

type RoleInclusionDAO interface {
	Create(ctx context.Context, inclusion RoleInclusion) (RoleInclusion, error)
	FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]RoleInclusion, error)
	FindByBizIdAndIncludingIds(ctx context.Context, bizId int64, IncludingIds []int64) ([]RoleInclusion, error)
	FindByBizIdAndIncludedIds(ctx context.Context, bizId int64, IncludedIds []int64) ([]RoleInclusion, error)
	DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error
	FindByBizIDAndID(ctx context.Context, bizID, id int64) (RoleInclusion, error)
}

type roleInclusionDao struct {
	db *egorm.Component
}

func (r *roleInclusionDao) FindByBizIDAndID(ctx context.Context, bizID, id int64) (RoleInclusion, error) {
	ri := RoleInclusion{}
	err := r.db.WithContext(ctx).Model(&RoleInclusion{}).Where("biz_id=? AND id=?", bizID, id).First(&ri).Error
	return ri, err
}

func (r *roleInclusionDao) Create(ctx context.Context, inclusion RoleInclusion) (RoleInclusion, error) {
	now := time.Now().Unix()
	inclusion.Utime = now
	inclusion.Ctime = now
	err := r.db.WithContext(ctx).Model(&RoleInclusion{}).Create(&inclusion).Error
	return inclusion, err
}

func (r *roleInclusionDao) FindByBizID(ctx context.Context, bizId int64, offset, limit int) ([]RoleInclusion, error) {
	ri := make([]RoleInclusion, 0)
	err := r.db.WithContext(ctx).Model(&RoleInclusion{}).Where("biz_id=?", bizId).Offset(offset).Limit(limit).Find(&ri).Error
	return ri, err
}

func (r *roleInclusionDao) FindByBizIdAndIncludingIds(ctx context.Context, bizId int64, IncludingIds []int64) ([]RoleInclusion, error) {
	ri := make([]RoleInclusion, 0)
	err := r.db.WithContext(ctx).Model(&RoleInclusion{}).Where("biz_id=? AND including_role_id in (?)", bizId, IncludingIds).Find(&ri).Error
	return ri, err
}

func (r *roleInclusionDao) FindByBizIdAndIncludedIds(ctx context.Context, bizId int64, IncludedIds []int64) ([]RoleInclusion, error) {
	ri := make([]RoleInclusion, 0)
	err := r.db.WithContext(ctx).Model(&RoleInclusion{}).Where("biz_id=? AND included_role_id IN ?", bizId, IncludedIds).Find(&ri).Error
	return ri, err
}

func (r *roleInclusionDao) DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error {
	return r.db.WithContext(ctx).Model(&RoleInclusion{}).Where("biz_id=? AND id=?", bizID, id).Delete(&RoleInclusion{}).Error
}

func NewRoleInclusionDAO(db *egorm.Component) RoleInclusionDAO {
	return &roleInclusionDao{db: db}
}
