package abac

import (
	"context"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/service/abac"
)

type ABACAttributeDefinitionServer struct {
	baseServer
	svc abac.AttributeDefinitionSvc
	permissionv1.UnimplementedAttributeDefinitionServiceServer
}

func (a *ABACAttributeDefinitionServer) Save(ctx context.Context, request *permissionv1.AttributeDefinitionServiceSaveRequest) (*permissionv1.AttributeDefinitionServiceSaveResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	id, err := a.svc.Create(ctx, bizId, a.toDomainAttributeDefinition(request.Definition))
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeDefinitionServiceSaveResponse{Id: int64(id)}, nil
}

func (a *ABACAttributeDefinitionServer) First(ctx context.Context, request *permissionv1.AttributeDefinitionServiceFirstRequest) (*permissionv1.AttributeDefinitionServiceFirstResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	found, err := a.svc.FindByBizIdAndId(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeDefinitionServiceFirstResponse{Definition: a.convertToProtoAttributeDefinition(found)}, nil

}

func (a *ABACAttributeDefinitionServer) Delete(ctx context.Context, request *permissionv1.AttributeDefinitionServiceDeleteRequest) (*permissionv1.AttributeDefinitionServiceDeleteResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.Delete(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeDefinitionServiceDeleteResponse{}, nil
}

func (a *ABACAttributeDefinitionServer) Find(ctx context.Context, request *permissionv1.AttributeDefinitionServiceFindRequest) (*permissionv1.AttributeDefinitionServiceFindResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	bizDefinitions, err := a.svc.FindByBizID(ctx, bizId)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeDefinitionServiceFindResponse{BizDefinition: a.convertToProtoBizDefinition(bizDefinitions)}, nil
}
