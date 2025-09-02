package repository

import (
	"context"
	"github.com/ecodeclub/ekit/mapx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/domain"
)

type RoleInclusionReloadCacheRepository struct {
	repo          *roleIncludeRepository
	userRoleRepo  *userRoleRepository
	cacheReloader UserPermissionCacheReloader
	logger        *elog.Component
}

func NewRoleInclusionReloadCacheRepository(repo *roleIncludeRepository, userRepo *userRoleRepository, cacheReloader UserPermissionCacheReloader) *RoleInclusionReloadCacheRepository {
	return &RoleInclusionReloadCacheRepository{
		repo:          repo,
		userRoleRepo:  userRepo,
		cacheReloader: cacheReloader,
		logger:        elog.DefaultLogger.With(elog.FieldName("RoleInclusionReloadCache")),
	}
}

func (r *RoleInclusionReloadCacheRepository) Create(ctx context.Context, inclusion domain.RoleInclusion) (domain.RoleInclusion, error) {
	created, err := r.repo.Create(ctx, inclusion)
	if err != nil {
		return domain.RoleInclusion{}, err
	}
	err1 := r.cacheReloader.Reload(ctx, r.getAffectUsers(ctx, created.BizID, created.IncludingRole.ID))
	if err1 != nil {
		r.logger.Warn("创建角色包含成功后，重新加载所有受影响用户缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", created.BizID),
			elog.Any("includingRoleID", created.IncludingRole.ID),
			elog.Any("includedRoleID", created.IncludedRole.ID),
		)

	}
	return created, nil
}

func (r *RoleInclusionReloadCacheRepository) getAffectUsers(ctx context.Context, bizID, includedRoleId int64) []domain.User {
	roleIds, err := r.getAffectedRoleIDs(ctx, bizID, includedRoleId)
	if err != nil {
		return nil
	}
	userRoles, err := r.userRoleRepo.FindByBizIDAndRoleIDs(ctx, bizID, roleIds)
	if err != nil {
		return nil
	}
	return slice.Map(userRoles, func(idx int, src domain.UserRole) domain.User {
		return domain.User{
			BizID: src.BizID,
			ID:    src.UserID,
		}
	})
}

func (r *RoleInclusionReloadCacheRepository) getAffectedRoleIDs(ctx context.Context, bizID int64, includeRoleID int64) ([]int64, error) {
	allRoleIDs := make(map[int64]any)
	allRoleIDs[includeRoleID] = struct{}{}

	inlcudedIDs := []int64{includeRoleID}
	for {
		inclusions, err := r.repo.FindByBizIdAndIncludedIds(ctx, bizID, inlcudedIDs)
		if err != nil {
			return nil, err
		}
		if len(inclusions) == 0 {
			break
		}
		inlcudedIDs = slice.Map(inclusions, func(idx int, src domain.RoleInclusion) int64 {
			allRoleIDs[src.IncludingRole.ID] = struct{}{}
			return src.IncludingRole.ID
		})
	}
	return mapx.Keys(allRoleIDs), nil
}
func (r *RoleInclusionReloadCacheRepository) FindByBizIDAndID(ctx context.Context, bizID, id int64) (domain.RoleInclusion, error) {
	return r.repo.FindByBizIDAndID(ctx, bizID, id)
}

func (r *RoleInclusionReloadCacheRepository) FindByBizIDAndIncludingRoleIDs(ctx context.Context, bizID int64, includingRoleIDs []int64) ([]domain.RoleInclusion, error) {
	return r.repo.FindByBizIdAndIncludingIds(ctx, bizID, includingRoleIDs)
}

func (r *RoleInclusionReloadCacheRepository) FindByBizIDAndIncludedRoleIDs(ctx context.Context, bizID int64, includedRoleIDs []int64) ([]domain.RoleInclusion, error) {
	return r.repo.FindByBizIdAndIncludedIds(ctx, bizID, includedRoleIDs)
}
func (r *RoleInclusionReloadCacheRepository) DeleteByBizIDAndID(ctx context.Context, bizID, id int64) error {
	deleted, err := r.repo.FindByBizIDAndID(ctx, bizID, id)
	if err != nil {
		return err
	}
	err = r.repo.DeleteByBizIDAndID(ctx, bizID, id)
	if err != nil {
		return err
	}
	if err1 := r.cacheReloader.Reload(ctx, r.getAffectUsers(ctx, deleted.BizID, deleted.IncludingRole.ID)); err1 != nil {
		r.logger.Warn("删除角色包含关系成功后，重新加载所有受影响用户的缓存失败",
			elog.FieldErr(err1),
			elog.Any("bizID", deleted.BizID),
			elog.Any("includingRoleID", deleted.IncludingRole.ID),
			elog.Any("includedRoleID", deleted.IncludedRole.ID),
		)
	}
	return nil
}
func (r *RoleInclusionReloadCacheRepository) FindByBizID(ctx context.Context, bizID int64, offset, limit int) ([]domain.RoleInclusion, error) {
	return r.repo.FindByBizID(ctx, bizID, offset, limit)
}
