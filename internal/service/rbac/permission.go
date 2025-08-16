package rbac

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
)

type PermissionService interface {
	Check(ctx context.Context, bizId, userId int64, resource domain.Resource, actions []string) (bool, error)
}

type permissionService struct {
	userPermissionRepo repository.UserPermissionRepository
}

func (p *permissionService) Check(ctx context.Context, bizId, userId int64, resource domain.Resource, actions []string) (bool, error) {
	allUserPermissions, err := p.userPermissionRepo.GetALLUserPermission(ctx, bizId, userId)
	if err != nil {
		return false, err
	}
	var res bool
	for _, p := range allUserPermissions {
		pr := p.Permission.Resource
		if resource.Type == pr.Type && resource.Key == pr.Key && slice.Contains(actions, p.Permission.Action) {
			if p.Effect.IsDeny() {
				return false, nil
			}
			res = true
		}
	}
	return res, nil

}

func NewPermissionService(userPermissionRepo repository.UserPermissionRepository) PermissionService {
	return &permissionService{userPermissionRepo: userPermissionRepo}
}
