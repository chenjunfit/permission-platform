package repository

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository/dao"
)

func GenDomainPolicyRules(rules []dao.PolicyRule) []domain.PolicyRule {
	//dao policyRule只有id关联的关系，需要将关系转换成domain PolicyRule
	ruleMap := make(map[int64]dao.PolicyRule, len(rules))
	for _, rule := range rules {
		ruleMap[rule.ID] = rule
	}
	//获取所有的根结点的ID
	//根rule left==nil right==nil
	rootRules := findRootRules(rules, ruleMap)

	for index := range rootRules {
		rootRules[index] = genRule(rootRules[index], ruleMap)
	}
	return rootRules
}

func findRootRules(rules []dao.PolicyRule, ruleMap map[int64]dao.PolicyRule) []domain.PolicyRule {
	childMap := make(map[int64]struct{})
	rootRules := make([]domain.PolicyRule, 0)
	for _, rule := range rules {
		if rule.Left > 0 {
			//将左节点存起来
			childMap[rule.Left] = struct{}{}
		}
		if rule.Right > 0 {
			//将左节点存起来
			childMap[rule.Right] = struct{}{}
		}
	}
	//childMap将所有的左/右节点存储,不在其中的就是根节点
	for _, rule := range rules {
		if _, ok := childMap[rule.ID]; !ok {
			domainPolicyRule := domain.PolicyRule{
				ID: rule.ID,
			}
			rootRules = append(rootRules, domainPolicyRule)
		}
	}
	return rootRules
}
func genRule(rule domain.PolicyRule, ruleMap map[int64]dao.PolicyRule) domain.PolicyRule {
	ruleDao, ok := ruleMap[rule.ID]
	if !ok {
		return domain.PolicyRule{}
	}
	rule = domain.PolicyRule{
		ID: ruleDao.ID,
		AttrDef: domain.AttributeDefinition{
			ID: ruleDao.AttrDefID,
		},
		Value:    ruleDao.Value,
		Operator: domain.RuleOperator(ruleDao.Operator),
		Ctime:    ruleDao.Ctime,
		Utime:    ruleDao.Utime,
	}
	if ruleDao.Left > 0 {
		left := genRule(domain.PolicyRule{ID: ruleDao.Left}, ruleMap)
		rule.LeftRule = &left
	}
	if ruleDao.Right > 0 {
		right := genRule(domain.PolicyRule{ID: ruleDao.Right}, ruleMap)
		rule.RightRule = &right
	}
	return rule
}
