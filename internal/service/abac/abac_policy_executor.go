package abac

import (
	"github.com/ecodeclub/ekit/mapx"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/abac/evaluator"
)

type PolicyExecutor interface {
	Check(policy domain.Policy, subject domain.ABACObject, resource domain.ABACObject, enviroment domain.ABACObject) bool
}

// 基于逻辑运算符的方法
type logicOperatorExecutor struct {
	selector evaluator.Selector
	logger   *elog.Component
}

func (l *logicOperatorExecutor) Check(policy domain.Policy, subject domain.ABACObject, resource domain.ABACObject, enviroment domain.ABACObject) bool {
	subjectMap := subject.ValuesMap()
	resourceMap := resource.ValuesMap()
	enviromentMap := enviroment.ValuesMap()
	allAttributeMap := mapx.Merge(subjectMap, resourceMap, enviromentMap)
	res := true
	for index := range policy.Rules {
		rule := policy.Rules[index]
		res = res && l.checkOneRule(rule, allAttributeMap)
	}
	return res
}

func NewPolicyExecutor(selector evaluator.Selector) PolicyExecutor {
	return &logicOperatorExecutor{
		selector: selector,
		logger:   elog.DefaultLogger,
	}
}

func (l *logicOperatorExecutor) checkOneRule(rule domain.PolicyRule, values map[int64]domain.AttributeValue) bool {
	if rule.LeftRule == nil && rule.RightRule == nil {
		val := values[rule.AttrDef.ID]
		actualVal := val.Value
		checker, err := l.selector.Select(val.AttrDef.DataType)
		if err != nil {
			return false
		}
		ok, err := checker.Evaluator(rule.Value, actualVal, rule.Operator)
		if err != nil {
			return false
		}
		return ok
	}
	left, right := true, true
	if rule.LeftRule != nil {
		left = l.checkOneRule(*rule.LeftRule, values)
	}
	if rule.RightRule != nil {
		right = l.checkOneRule(*rule.RightRule, values)
	}
	switch rule.Operator {
	case domain.AND:
		return left && right
	case domain.OR:
		return left || right
	case domain.NOT:
		return !right
	default:
		return false
	}
}
