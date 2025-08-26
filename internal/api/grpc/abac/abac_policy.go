package abac

import (
	"context"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/service/abac"
)

type ABACPolicyServer struct {
	baseServer
	permissionv1.UnsafePolicyServiceServer
	svc abac.PolicySvc
}

func (a *ABACPolicyServer) Save(ctx context.Context, request *permissionv1.PolicyServiceSaveRequest) (*permissionv1.PolicyServiceSaveResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	request.Policy.BizId = bizId
	id, err := a.svc.Save(ctx, a.convertToDomainPolicy(request.Policy))
	return &permissionv1.PolicyServiceSaveResponse{Id: id}, nil
}

func (a *ABACPolicyServer) Delete(ctx context.Context, request *permissionv1.PolicyServiceDeleteRequest) (*permissionv1.PolicyServiceDeleteResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.Delete(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceDeleteResponse{}, nil
}

func (a *ABACPolicyServer) First(ctx context.Context, request *permissionv1.PolicyServiceFirstRequest) (*permissionv1.PolicyServiceFirstResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	policy, err := a.svc.First(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceFirstResponse{
		Policy: a.convertToProtoPolicy(policy),
	}, nil
}

func (a *ABACPolicyServer) SaveRule(ctx context.Context, request *permissionv1.PolicyServiceSaveRuleRequest) (*permissionv1.PolicyServiceSaveRuleResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	rule := a.convertToDomainPolicyRule(request.Rule)
	id, err := a.svc.SaveRule(ctx, bizId, request.PolicyId, rule) // Dereference the pointer
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceSaveRuleResponse{
		Id: id,
	}, nil
}

func (a *ABACPolicyServer) DeleteRule(ctx context.Context, request *permissionv1.PolicyServiceDeleteRuleRequest) (*permissionv1.PolicyServiceDeleteRuleResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.DeleteRule(ctx, bizId, request.RuleId, false)
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceDeleteRuleResponse{}, nil
}

func (a *ABACPolicyServer) SavePermissionPolicy(ctx context.Context, request *permissionv1.PolicyServiceSavePermissionPolicyRequest) (*permissionv1.PolicyServiceSavePermissionPolicyResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.SavePermissionPolicy(ctx, bizId, request.PolicyId, request.PermissionId, a.convertToDomainEffect(request.Effect))
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceSavePermissionPolicyResponse{}, nil

}

func (a *ABACPolicyServer) FindPolicies(ctx context.Context, request *permissionv1.PolicyServiceFindPoliciesRequest) (*permissionv1.PolicyServiceFindPoliciesResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	total, policies, err := a.svc.FindPolicies(ctx, bizId, int(request.Offset), int(request.Limit))
	if err != nil {
		return nil, err
	}
	return &permissionv1.PolicyServiceFindPoliciesResponse{
		Total:    total,
		Policies: a.convertToProtoPolicies(policies),
	}, nil
}

func NewABACPolicyServer(svc abac.PolicySvc) *ABACPolicyServer {
	return &ABACPolicyServer{svc: svc}
}
