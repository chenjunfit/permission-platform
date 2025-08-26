package abac

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
)

type PolicySvc interface {
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

type policySvc struct {
	repository.AttributePolicyRepository
}

func NewPolicySvc(repo repository.AttributePolicyRepository) PolicySvc {
	return &policySvc{repo}
}
