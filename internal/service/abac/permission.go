package abac

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/repository"
	"golang.org/x/sync/errgroup"
)

type PermissionSvc interface {
	Check(ctx context.Context, bizId, uid int64, resource domain.Resource, action []string, attrs domain.Attributes) (bool, error)
}

type permissionSvc struct {
	permissionRepo repository.PermissionRepository
	resourceRepo   repository.ResourceRepository
	policyRepo     repository.AttributePolicyRepository
	valRepo        repository.AttributeValueRepository
	attrRepo       repository.AttributeDefinitionRepository
	parser         PolicyExecutor
}

func NewPermissionSvc(
	permissionRepo repository.PermissionRepository,
	resourceRepo repository.ResourceRepository,
	policyRepo repository.AttributePolicyRepository,
	valRepo repository.AttributeValueRepository,
	attrRepo repository.AttributeDefinitionRepository,
	parser PolicyExecutor,
) PermissionSvc {
	return &permissionSvc{
		permissionRepo: permissionRepo,
		resourceRepo:   resourceRepo,
		policyRepo:     policyRepo,
		valRepo:        valRepo,
		attrRepo:       attrRepo,
		parser:         parser,
	}
}
func (p *permissionSvc) Check(ctx context.Context, bizId, uid int64, resource domain.Resource, action []string, attrs domain.Attributes) (bool, error) {
	permissions, res, bizDefinition, err := p.getPermissionAndRes(ctx, bizId, resource, action)
	if err != nil {
		return false, err
	}
	permissionIds := slice.Map(permissions, func(idx int, src domain.Permission) int64 {
		return src.ID
	})
	resource.ID = res.ID

	var (
		eg       errgroup.Group
		subObj   domain.ABACObject
		resObj   domain.ABACObject
		envObj   domain.ABACObject
		policies []domain.Policy
	)
	eg.Go(func() error {
		var err error
		subObj, err = p.valRepo.FindSubjectValue(ctx, bizId, uid)
		subObj.FillDefinitions(bizDefinition.SubjectAttrDefs)
		return err
	})
	eg.Go(func() error {
		var err error
		resObj, err = p.valRepo.FindResourceValue(ctx, bizId, resource.ID)
		resObj.FillDefinitions(bizDefinition.ResourceAttrDefs)
		return err
	})
	eg.Go(func() error {
		var err error
		envObj, err = p.valRepo.FindEnvironmentValue(ctx, bizId)
		envObj.FillDefinitions(bizDefinition.EnvironmentAttrDefs)
		return err
	})
	eg.Go(func() error {
		var err error
		policies, err = p.policyRepo.FindPoliciesByPermissionIDs(ctx, bizId, permissionIds)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return false, err
	}
	subObj.MergeRealTimeAttr(bizDefinition.SubjectAttrDefs, attrs.Subject)
	resObj.MergeRealTimeAttr(bizDefinition.ResourceAttrDefs, attrs.Resource)
	envObj.MergeRealTimeAttr(bizDefinition.EnvironmentAttrDefs, attrs.Environment)
	var hasPermit bool
	var hasDeny bool
	if len(policies) == 0 {
		return false, nil
	}
	for index := range policies {
		policy := policies[index]
		if p.parser.Check(policy, subObj, resObj, envObj) {
			for index := range policy.Permissions {
				perm := policy.Permissions[index]
				if perm.Effect == domain.EffectAllow {
					hasPermit = true
				}
				if perm.Effect == domain.EffectDeny {
					hasDeny = true
				}
			}
		}
	}
	if hasDeny {
		return false, nil
	}
	if hasPermit {
		return true, nil
	}
	return false, nil
}
func (p *permissionSvc) getPermissionAndRes(ctx context.Context, bizId int64, resource domain.Resource, action []string) ([]domain.Permission, domain.Resource, domain.BizAttrDefinition, error) {
	var (
		eg          errgroup.Group
		permissions []domain.Permission
		res         domain.Resource
		bizDef      domain.BizAttrDefinition
	)
	eg.Go(func() error {
		var err error
		permissions, err = p.permissionRepo.FindPermissions(ctx, bizId, resource.Type, resource.Key, action)
		return err
	})
	eg.Go(func() error {
		var err error
		res, err = p.resourceRepo.FindByBizIDAndTypeAndKey(ctx, bizId, resource.Type, resource.Key)
		return err
	})
	eg.Go(func() error {
		var err error
		bizDef, err = p.attrRepo.FindByBizID(ctx, bizId)
		return err
	})
	err := eg.Wait()
	return permissions, res, bizDef, err
}
