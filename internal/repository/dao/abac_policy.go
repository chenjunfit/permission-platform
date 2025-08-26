package dao

import (
	"context"
	"errors"
	"github.com/ego-component/egorm"
	"github.com/permission-dev/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Policy 策略表模型
type Policy struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID       int64  `gorm:"column:biz_id;index:idx_biz_id;comment:业务ID"`
	Name        string `gorm:"column:name;type:varchar(100);not null;uniqueIndex:idx_biz_name;comment:策略名称" json:"name"`
	Description string `gorm:"column:description;type:text;comment:策略描述" json:"description"`
	Status      string `gorm:"column:status;type:enum('active','inactive');not null;default:active;index:idx_status;comment:策略状态" json:"status"`
	ExecuteType string `gorm:"column:execute_type;type:varchar(255);default:logic"`
	Ctime       int64  `gorm:"column:ctime;comment:创建时间"`
	Utime       int64  `gorm:"column:utime;comment:更新时间"`
}

// TableName 指定表名
func (p Policy) TableName() string {
	return "policies"
}

/*
ABAC 模型允许在同一策略下对同一属性定义设置 多个规则 。例如：

- 一个策略可能同时包含 user.level > 5 和 user.level < 10 两条规则（通过 Left 和 Right 字段组合）；
- 不同操作符（ > , < , = , IN 等）搭配同一属性定义可以形成互补的条件判断。
*/
// PolicyRule 策略规则表模型
type PolicyRule struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID     int64  `gorm:"column:biz_id;index:idx_biz_id;comment:业务ID"`
	PolicyID  int64  `gorm:"column:policy_id;not null;index:idx_policy_id;comment:策略ID"`
	AttrDefID int64  `gorm:"column:attr_def_id;not null;index:idx_attr_def_id;comment:属性定义ID"`
	Value     string `gorm:"column:value;type:text;comment:比较值，取决于类型"`
	Left      int64  `gorm:"column:left;comment:左规则ID"`
	Right     int64  `gorm:"column:right;comment:右规则ID"`
	Operator  string `gorm:"column:operator;type:varchar(255);not null;comment:操作符"`
	Ctime     int64  `gorm:"column:ctime;comment:创建时间"`
	Utime     int64  `gorm:"column:utime;comment:更新时间"`
}

// TableName 指定表名
func (r PolicyRule) TableName() string {
	return "policy_rules"
}

// PermissionPolicy 权限策略关联表模型
type PermissionPolicy struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	BizID        int64  `gorm:"column:biz_id;index:idx_biz_id;comment:业务ID;uniqueIndex:idx_permission_policy_bizId"`
	Effect       string `gorm:"column:effect;type:varchar(50)"`
	PermissionID int64  `gorm:"column:permission_id;not null;uniqueIndex:idx_permission_policy_bizId;index:idx_permission_id;comment:权限ID"`
	PolicyID     int64  `gorm:"column:policy_id;not null;uniqueIndex:idx_permission_policy_bizId;index:idx_policy_id;comment:策略ID"`
	Ctime        int64  `gorm:"column:ctime;comment:创建时间"`
	Utime        int64  `gorm:"column:utime;comment:创建时间"`
}

// TableName 指定表名
func (p PermissionPolicy) TableName() string {
	return "permission_policies"
}

type PolicyDAO interface {
	//policy
	SavePolicy(ctx context.Context, policy Policy) (int64, error)
	DeletePolicy(ctx context.Context, bizID, id int64) error
	UpdatePolicyStatus(ctx context.Context, id int64, status string) error
	FindPolicyById(ctx context.Context, bizId, id int64) (Policy, error)
	FindPolicyByIds(ctx context.Context, ids []int64) ([]Policy, error)
	FindPoliciesByBizId(ctx context.Context, bizId int64) ([]Policy, error)
	PolicyList(ctx context.Context, bizId int64, offset, limit int) ([]Policy, error)
	FindPoliciesByBizIds(ctx context.Context, bizIds []int64) (map[int64][]Policy, error)
	PolicyListCount(ctx context.Context, bizId int64) (int64, error)

	//PolicyRule方法
	SavePolicyRule(ctx context.Context, rule PolicyRule) (int64, error)
	DeletePolicyRule(ctx context.Context, bizID, id int64) error
	DeletePolicyRuleCascade(ctx context.Context, bizID, id int64) error
	FindPolicyRule(ctx context.Context, id int64) (PolicyRule, error)
	FindPolicyRulesByPolicyID(ctx context.Context, bizID, policyID int64) ([]PolicyRule, error)
	FindPolicyRulesByPolicyIDs(ctx context.Context, policyIDs []int64) (map[int64][]PolicyRule, error)
	FindPolicyRulesByBiz(ctx context.Context, bizID int64) (map[int64][]PolicyRule, error)
	FindPoliciesRulesByBizIDs(ctx context.Context, bizIDs []int64) (map[int64]map[int64][]PolicyRule, error)

	//permission policy方法
	SavePermissionPolicy(ctx context.Context, permissionPolicy PermissionPolicy) error
	DeletePermissionPolicy(ctx context.Context, bizID int64, permissionID int64, policyID int64) error
	FindPoliciesByPermission(ctx context.Context, bizID int64, permissionIDs []int64) ([]PermissionPolicy, error)
	FindPermissionPolicy(ctx context.Context, bizID int64) (map[int64][]PermissionPolicy, error)
	FindPermissionPolicyByBizIDs(ctx context.Context, bizIDs []int64) (map[int64]map[int64][]PermissionPolicy, error)
}
type policyDao struct {
	db *egorm.Component
}

func (p *policyDao) SavePolicy(ctx context.Context, policy Policy) (int64, error) {
	now := time.Now().UnixMilli()
	if policy.ID == 0 {
		policy.Ctime = now
	}
	policy.Utime = now
	err := p.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"description", "utime"}),
		}).Create(&policy).Error
	return policy.ID, err

}

func (p *policyDao) DeletePolicy(ctx context.Context, bizID, id int64) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 删除策略与权限的关联关系
		if err := tx.Where("biz_id = ? AND policy_id = ?", bizID, id).Delete(&PermissionPolicy{}).Error; err != nil {
			return err
		}

		// 2. 删除策略关联的规则
		if err := tx.Where("biz_id = ? AND policy_id = ?", bizID, id).Delete(&PolicyRule{}).Error; err != nil {
			return err
		}

		// 3. 删除策略本身
		return tx.Where("id = ? AND biz_id = ?", id, bizID).Delete(&Policy{}).Error
	})
}

func (p *policyDao) UpdatePolicyStatus(ctx context.Context, id int64, status string) error {
	now := time.Now().UnixMilli()
	res := p.db.WithContext(ctx).Model(&Policy{}).Where("id=?", id).Updates(map[string]any{"status": status, "utime": now})
	if res.RowsAffected == 0 {
		return errors.New("更新数据库失败")
	}
	return res.Error
}

func (p *policyDao) FindPolicyById(ctx context.Context, bizId, id int64) (Policy, error) {
	var policy Policy
	err := p.db.WithContext(ctx).
		Where("id = ? AND biz_id = ?", id, bizId).
		First(&policy).Error
	return policy, err
}

func (p *policyDao) FindPolicyByIds(ctx context.Context, ids []int64) ([]Policy, error) {
	var policies []Policy
	err := p.db.WithContext(ctx).
		Where(" id IN ? AND status = ?", ids, domain.PolicyStatusActive).
		Find(&policies).Error
	return policies, err
}

func (p *policyDao) FindPoliciesByBizId(ctx context.Context, bizId int64) ([]Policy, error) {
	var policies []Policy
	err := p.db.WithContext(ctx).
		Where("biz_id = ?", bizId).
		Find(&policies).Error
	return policies, err
}

func (p *policyDao) PolicyList(ctx context.Context, bizId int64, offset, limit int) ([]Policy, error) {
	var list []Policy
	err := p.db.WithContext(ctx).Where("biz_id = ?", bizId).
		Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (p *policyDao) FindPoliciesByBizIds(ctx context.Context, bizIds []int64) (map[int64][]Policy, error) {
	var policies []Policy
	err := p.db.WithContext(ctx).Model(&Policy{}).Where("biz_id IN ?", bizIds).Find(&policies).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]Policy, len(policies))
	for _, policy := range policies {
		result[policy.BizID] = append(result[policy.BizID], policy)
	}
	return result, nil
}

func (p *policyDao) PolicyListCount(ctx context.Context, bizId int64) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Model(&Policy{}).
		Where("biz_id = ?", bizId).Count(&count).Error
	return count, err
}

func (p *policyDao) SavePolicyRule(ctx context.Context, rule PolicyRule) (int64, error) {
	now := time.Now().UnixMilli()
	if rule.ID == 0 {
		rule.Ctime = now
	}
	rule.Utime = now

	err := p.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "left", "right", "operator", "utime"}),
		}).Create(&rule).Error
	return rule.ID, err
}

func (p *policyDao) DeletePolicyRule(ctx context.Context, bizID, id int64) error {
	return p.db.WithContext(ctx).
		Where("id = ? AND biz_id = ?", id, bizID).
		Delete(&PolicyRule{}).Error
}
func (p *policyDao) DeletePolicyRuleCascade(ctx context.Context, bizID, id int64) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		// 1. 递归删除所有子规则
		if err := p.deleteChildRules(tx, bizID, id); err != nil {
			return err
		}
		// 2. 删除当前规则
		return tx.Where("id = ? AND biz_id = ?", id, bizID).Delete(&PolicyRule{}).Error
	})
}
func (p *policyDao) deleteChildRules(tx *gorm.DB, bizID, parentRuleID int64) error {
	// 查询当前规则的子规则
	var childRules []PolicyRule
	err := tx.Where("biz_id = ? AND (left = ? OR right = ?)", bizID, parentRuleID, parentRuleID).Find(&childRules).Error
	if err != nil {
		return err
	}
	// 递归删除每个子规则
	for _, childRule := range childRules {
		if err := p.deleteChildRules(tx, bizID, childRule.ID); err != nil {
			return err
		}
		// 删除子规则
		if err := tx.Where("id = ? AND biz_id = ?", childRule.ID, bizID).Delete(&PolicyRule{}).Error; err != nil {
			return err
		}
	}

	return nil
}
func (p *policyDao) FindPolicyRule(ctx context.Context, id int64) (PolicyRule, error) {
	var rule PolicyRule
	err := p.db.WithContext(ctx).
		Where("id = ?", id).
		First(&rule).Error
	return rule, err
}

func (p *policyDao) FindPolicyRulesByPolicyID(ctx context.Context, bizID, policyID int64) ([]PolicyRule, error) {
	var rules []PolicyRule
	err := p.db.WithContext(ctx).
		Where(" policy_id = ? AND biz_id = ?", policyID, bizID).
		Find(&rules).Error
	return rules, err
}

func (p *policyDao) FindPolicyRulesByPolicyIDs(ctx context.Context, policyIDs []int64) (map[int64][]PolicyRule, error) {
	var rules []PolicyRule
	err := p.db.WithContext(ctx).
		Where(" policy_id IN ?", policyIDs).
		Find(&rules).Error
	if err != nil {
		return nil, err
	}

	// 将规则按策略ID分组
	result := make(map[int64][]PolicyRule)
	for _, rule := range rules {
		result[rule.PolicyID] = append(result[rule.PolicyID], rule)
	}

	return result, nil
}

func (p *policyDao) FindPolicyRulesByBiz(ctx context.Context, bizID int64) (map[int64][]PolicyRule, error) {
	var rules []PolicyRule
	err := p.db.WithContext(ctx).
		Where(" biz_id = ?", bizID).
		Find(&rules).Error
	if err != nil {
		return nil, err
	}
	// 将规则按策略ID分组
	result := make(map[int64][]PolicyRule)
	for _, rule := range rules {
		result[rule.PolicyID] = append(result[rule.PolicyID], rule)
	}
	return result, nil
}

func (p *policyDao) FindPoliciesRulesByBizIDs(ctx context.Context, bizIDs []int64) (map[int64]map[int64][]PolicyRule, error) {
	var rules []PolicyRule
	err := p.db.WithContext(ctx).
		Where("biz_id IN ?", bizIDs).
		Find(&rules).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int64]map[int64][]PolicyRule)
	for _, rule := range rules {
		if _, exists := result[rule.BizID]; !exists {
			result[rule.BizID] = make(map[int64][]PolicyRule)
		}
		result[rule.BizID][rule.PolicyID] = append(result[rule.BizID][rule.PolicyID], rule)
	}

	return result, nil
}

func (p *policyDao) SavePermissionPolicy(ctx context.Context, permissionPolicy PermissionPolicy) error {
	permissionPolicy.Ctime = time.Now().UnixMilli()
	permissionPolicy.Utime = time.Now().UnixMilli()
	err := p.db.WithContext(ctx).Create(&permissionPolicy).Error
	return err
}

func (p *policyDao) DeletePermissionPolicy(ctx context.Context, bizID int64, permissionID int64, policyID int64) error {
	return p.db.WithContext(ctx).
		Where("biz_id = ? AND permission_id = ? AND policy_id = ?", bizID, permissionID, policyID).
		Delete(&PermissionPolicy{}).Error
}

func (p *policyDao) FindPoliciesByPermission(ctx context.Context, bizID int64, permissionIDs []int64) ([]PermissionPolicy, error) {
	var relations []PermissionPolicy
	err := p.db.WithContext(ctx).
		Where("biz_id = ? AND permission_id in ?", bizID, permissionIDs).
		Find(&relations).Error
	return relations, err
}

func (p *policyDao) FindPermissionPolicy(ctx context.Context, bizID int64) (map[int64][]PermissionPolicy, error) {
	var list []PermissionPolicy
	err := p.db.WithContext(ctx).
		Model(&PermissionPolicy{}).
		Where("biz_id = ?", bizID).
		Find(&list).Error
	result := make(map[int64][]PermissionPolicy)
	for idx := range list {
		permissionPolicy := list[idx]
		result[permissionPolicy.PolicyID] = append(result[permissionPolicy.PolicyID], permissionPolicy)
	}
	return result, err
}

func (p *policyDao) FindPermissionPolicyByBizIDs(ctx context.Context, bizIDs []int64) (map[int64]map[int64][]PermissionPolicy, error) {
	var policies []PermissionPolicy
	err := p.db.WithContext(ctx).
		Where("biz_id in ?", bizIDs).
		Find(&policies).Error
	if err != nil {
		return nil, err
	}

	// Create a nested map: bizID -> policyID -> []PermissionPolicy
	result := make(map[int64]map[int64][]PermissionPolicy)
	for _, policy := range policies {
		if _, exists := result[policy.BizID]; !exists {
			result[policy.BizID] = make(map[int64][]PermissionPolicy)
		}
		result[policy.BizID][policy.PolicyID] = append(result[policy.BizID][policy.PolicyID], policy)
	}

	return result, nil
}

func NewPolicyDAO(db *egorm.Component) PolicyDAO {
	return &policyDao{db: db}
}
