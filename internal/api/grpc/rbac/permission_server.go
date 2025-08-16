package rbac

import (
	"context"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/rbac"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PermissionServer struct {
	baseServer
	permissionv1.UnimplementedPermissionServiceServer
	permissionSvc rbac.PermissionService
}

func (p *PermissionServer) CheckPermission(ctx context.Context, in *permissionv1.CheckPermissionRequest) (*permissionv1.CheckPermissionResponse, error) {
	//参数校验
	if in.Uid <= 0 || in.Permission.ResourceKey == "" || in.Permission.ResourceType == "" || len(in.Permission.Actions) == 0 {
		return &permissionv1.CheckPermissionResponse{Allowed: false}, status.Error(codes.InvalidArgument, "参数无效")
	}
	bizId, err := p.getBizIDFromContext(ctx)
	if err != nil {
		return &permissionv1.CheckPermissionResponse{Allowed: false}, status.Error(codes.Unauthenticated, err.Error())
	}
	allow, err := p.permissionSvc.Check(ctx, bizId, in.Uid, domain.Resource{
		BizID: bizId,
		Type:  in.Permission.ResourceType,
		Key:   in.Permission.ResourceKey,
	}, in.Permission.Actions)
	if err != nil {
		return &permissionv1.CheckPermissionResponse{Allowed: false}, status.Error(codes.Internal, err.Error())

	}
	return &permissionv1.CheckPermissionResponse{Allowed: allow}, nil
}

func NewPermissionServer(permissionSvc rbac.PermissionService) *PermissionServer {
	return &PermissionServer{permissionSvc: permissionSvc}
}
