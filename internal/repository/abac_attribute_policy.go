package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
	"golang.org/x/sync/errgroup"
)

type AttributePolicyRepository interface {
	//policy相关
	Save(ctx context.Context, policy domain.Policy) (int64, error)
	Delete(ctx context.Context, bizID, id int64) error
	First(ctx context.Context, bizID, id int64) (domain.Policy, error) // 包含规则
	//policy rule相关
	SaveRule(ctx context.Context, bizID, policyId int64, rule domain.PolicyRule) (int64, error)
	DeleteRule(ctx context.Context, bizID, ruleID int64, cascade bool) error
	//policy permission相关
	FindPoliciesByPermissionIDs(ctx context.Context, bizID int64, permissionIDs []int64) ([]domain.Policy, error)
	SavePermissionPolicy(ctx context.Context, bizID, policyID, permissionID int64, effect domain.Effect) error
	FindPolicies(ctx context.Context, bizID int64, offset, limit int) (int64, []domain.Policy, error)
	FindBizPolicies(ctx context.Context, bizID int64) ([]domain.Policy, error)
}

type attributePolicyRepository struct {
	policyDAO dao.PolicyDAO
}

func (a *attributePolicyRepository) SaveRule(ctx context.Context, bizID, policyId int64, rule domain.PolicyRule) (int64, error) {
	ruleDao := dao.PolicyRule{
		ID:        rule.ID,
		BizID:     bizID,
		PolicyID:  policyId,
		AttrDefID: rule.AttrDef.ID,
		Value:     rule.Value,
		Operator:  rule.Operator.String(),
		Ctime:     rule.Ctime,
		Utime:     rule.Utime,
	}
	if rule.RightRule == nil {
		ruleDao.Right = 0
	}
	if rule.LeftRule == nil {
		ruleDao.Left = 0
	}
	ruleId, err := a.policyDAO.SavePolicyRule(ctx, ruleDao)
	return ruleId, err
}

/*
这里直接删除，会产生孤儿规则
1、级联删除 可以做
2、应用层上过滤掉 没想明白
*/

func (a *attributePolicyRepository) DeleteRule(ctx context.Context, bizID, ruleID int64, cascade bool) error {
	if cascade {
		return a.policyDAO.DeletePolicyRuleCascade(ctx, bizID, ruleID)
	}
	return a.policyDAO.DeletePolicyRule(ctx, bizID, ruleID)
}
func (a *attributePolicyRepository) Save(ctx context.Context, policy domain.Policy) (int64, error) {
	// 转换为 DAO 层的 Policy 对象
	policyDAO := dao.Policy{
		ID:          policy.ID,
		BizID:       policy.BizID,
		Name:        policy.Name,
		ExecuteType: string(policy.ExecuteType),
		Description: policy.Description,
		Status:      string(policy.Status),
	}
	// 保存策略
	id, err := a.policyDAO.SavePolicy(ctx, policyDAO)
	return id, err
}

func (a *attributePolicyRepository) Delete(ctx context.Context, bizID, id int64) error {
	return a.policyDAO.DeletePolicy(ctx, bizID, id)
}

func (a *attributePolicyRepository) First(ctx context.Context, bizID, id int64) (domain.Policy, error) {
	var policy dao.Policy
	var policyRules []dao.PolicyRule
	var eg errgroup.Group
	eg.Go(func() error {
		var err error
		policy, err = a.policyDAO.FindPolicyById(ctx, bizID, id)
		return err
	})
	eg.Go(func() error {
		var err error
		policyRules, err = a.policyDAO.FindPolicyRulesByPolicyID(ctx, bizID, id)
		return err
	})
	if err := eg.Wait(); err != nil {
		return domain.Policy{}, err
	}
	//toPolicyDomain
	//将所有的policy--policyrule转换成domain.policyRule
	return a.toPolicyDomain(policy, policyRules, map[int64][]dao.PermissionPolicy{}), nil
}
func (p *attributePolicyRepository) getPolicies(ctx context.Context, bizID int64) ([]domain.Policy, error) {
	var (
		eg                   errgroup.Group
		daoPolicies          []dao.Policy
		daoPolicyRules       map[int64][]dao.PolicyRule
		daoPolicyPermissions map[int64][]dao.PermissionPolicy
	)
	eg.Go(func() error {
		var eerr error
		daoPolicyPermissions, eerr = p.policyDAO.FindPermissionPolicy(ctx, bizID)
		return eerr
	})
	eg.Go(func() error {
		var eerr error
		daoPolicies, eerr = p.policyDAO.FindPoliciesByBizId(ctx, bizID)
		return eerr
	})
	eg.Go(func() error {
		var eerr error
		daoPolicyRules, eerr = p.policyDAO.FindPolicyRulesByBiz(ctx, bizID)
		return eerr
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	res := make([]domain.Policy, 0, len(daoPolicies))
	for index := range daoPolicies {
		daoPolicy := daoPolicies[index]
		rules := daoPolicyRules[daoPolicy.ID]
		res = append(res, p.toPolicyDomain(daoPolicy, rules, daoPolicyPermissions))
	}
	return res, nil
}
func (p *attributePolicyRepository) FindPoliciesByPermissionIDs(ctx context.Context, bizID int64, permissionIDs []int64) ([]domain.Policy, error) {
	//1、获取daoPolicys
	//2、获取daoPolicyRlues
	//3、获取daoPermissions
	//4、查询数据，转换数据
	policies, err := p.getPolicies(ctx, bizID)
	if err != nil {
		return nil, err
	}
	return p.getPolicyByPermissionID(policies, permissionIDs), nil

}
func (p *attributePolicyRepository) getPolicyByPermissionID(policies []domain.Policy, permissionIDs []int64) []domain.Policy {
	res := make([]domain.Policy, 0, len(policies))
	for idx := range policies {
		policy := policies[idx]
		if policy.ContainsAnyPermissions(permissionIDs) {
			res = append(res, policy)
		}
	}
	return res
}
func (p *attributePolicyRepository) toPolicyDomain(policy dao.Policy, rules []dao.PolicyRule, permissionPolicyMap map[int64][]dao.PermissionPolicy) domain.Policy {
	domainPolicy := domain.Policy{
		ID:          policy.ID,
		BizID:       policy.BizID,
		Name:        policy.Name,
		ExecuteType: domain.ExecuteType(policy.ExecuteType),
		Description: policy.Description,
		Status:      domain.PolicyStatusType(policy.Status),
		Rules:       GenDomainPolicyRules(rules),
	}

	if permissionPolicies, ok := permissionPolicyMap[policy.ID]; ok {
		for idx := range permissionPolicies {
			permissionPolicy := permissionPolicies[idx]
			domainPolicy.Permissions = append(domainPolicy.Permissions, domain.UserPermission{
				BizID: permissionPolicy.BizID,
				Permission: domain.Permission{
					ID: permissionPolicy.PermissionID,
				},
				Effect: domain.Effect(permissionPolicy.Effect),
			})
		}
	}
	return domainPolicy
}
func (p *attributePolicyRepository) SavePermissionPolicy(ctx context.Context, bizID, policyID, permissionID int64, effect domain.Effect) error {
	return p.policyDAO.SavePermissionPolicy(ctx, dao.PermissionPolicy{
		BizID:        bizID,
		Effect:       effect.String(),
		PermissionID: permissionID,
		PolicyID:     policyID,
	})
}
func (p *attributePolicyRepository) FindPolicies(ctx context.Context, bizID int64, offset, limit int) (int64, []domain.Policy, error) {
	var (
		eg       errgroup.Group
		res      []domain.Policy
		policies []dao.Policy
		count    int64
	)
	eg.Go(func() error {
		var eerr error
		policies, eerr = p.policyDAO.PolicyList(ctx, bizID, offset, limit)
		res = slice.Map(policies, func(idx int, src dao.Policy) domain.Policy {
			return p.toPolicyDomain(src, []dao.PolicyRule{}, map[int64][]dao.PermissionPolicy{})
		})
		return eerr
	})
	eg.Go(func() error {
		var eerr error
		count, eerr = p.policyDAO.PolicyListCount(ctx, bizID)
		return eerr

	})
	err := eg.Wait()
	return count, res, err

}
func (p *attributePolicyRepository) FindBizPolicies(ctx context.Context, bizID int64) ([]domain.Policy, error) {
	return p.getPolicies(ctx, bizID)
}

func NewAttributePolicyRepository(dao dao.PolicyDAO) AttributePolicyRepository {
	return &attributePolicyRepository{policyDAO: dao}
}
