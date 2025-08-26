package hybrid

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/abac"
	"github.com/permission-dev/internal/service/rbac"
)

type PermissionService interface {
	Check(ctx context.Context, bizID, userID int64, resource domain.Resource, actions []string, attrs domain.Attributes) (bool, error)
}

type permissionService struct {
	rbacSvc rbac.PermissionService
	abacSvc abac.PermissionSvc
}

func NewPermissionService(
	rbacSvc rbac.PermissionService,
	abacSvc abac.PermissionSvc,
) PermissionService {
	return &permissionService{
		rbacSvc: rbacSvc,
		abacSvc: abacSvc,
	}
}

func (p *permissionService) Check(ctx context.Context, bizID, userID int64, resource domain.Resource, actions []string, attrs domain.Attributes) (bool, error) {
	ok, err := p.rbacSvc.Check(ctx, bizID, userID, resource, actions)
	if err != nil || !ok {
		return false, err
	}
	return p.abacSvc.Check(ctx, bizID, userID, resource, actions, attrs)
}
