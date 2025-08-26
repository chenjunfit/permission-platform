package abac

import (
	"context"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/abac"
)

type ABACAttributeValServer struct {
	baseServer
	svc abac.AttributeValueSvc
	permissionv1.UnsafeAttributeValueServiceServer
}

func (a *ABACAttributeValServer) SaveSubjectValue(ctx context.Context, request *permissionv1.AttributeValueServiceSaveSubjectValueRequest) (*permissionv1.AttributeValueServiceSaveSubjectValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	val := domain.AttributeValue{
		AttrDef: a.toDomainAttributeDefinition(request.Value.Definition),
		Value:   request.Value.Value,
	}
	id, err := a.svc.SaveSubjectValue(ctx, bizId, request.SubjectId, val)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceSaveSubjectValueResponse{Id: id}, nil
}

func (a *ABACAttributeValServer) DeleteSubjectValue(ctx context.Context, request *permissionv1.AttributeValueServiceDeleteSubjectValueRequest) (*permissionv1.AttributeValueServiceDeleteSubjectValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.DeleteSubjectValue(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceDeleteSubjectValueResponse{}, nil
}

func (a *ABACAttributeValServer) FindSubjectValueWithDefinition(ctx context.Context, request *permissionv1.AttributeValueServiceFindSubjectValueWithDefinitionRequest) (*permissionv1.AttributeValueServiceFindSubjectValueWithDefinitionResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	valObject, err := a.svc.FindSubjectValueWithDefinition(ctx, bizId, request.SubjectId)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceFindSubjectValueWithDefinitionResponse{
		Subject: &permissionv1.SubjectObject{
			Id:              valObject.ID,
			AttributeValues: a.convertToProtoSubjectAttributeValues(valObject.AttrValues),
		},
	}, nil
}

func (a *ABACAttributeValServer) SaveResourceValue(ctx context.Context, request *permissionv1.AttributeValueServiceSaveResourceValueRequest) (*permissionv1.AttributeValueServiceSaveResourceValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	val := domain.AttributeValue{
		AttrDef: a.toDomainAttributeDefinition(request.Value.Definition),
		Value:   request.Value.Value,
	}
	id, err := a.svc.SaveResourceValue(ctx, bizId, request.ResourceId, val)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceSaveResourceValueResponse{Id: id}, nil
}

func (a *ABACAttributeValServer) DeleteResourceValue(ctx context.Context, request *permissionv1.AttributeValueServiceDeleteResourceValueRequest) (*permissionv1.AttributeValueServiceDeleteResourceValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.DeleteResourceValue(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceDeleteResourceValueResponse{}, nil
}

func (a *ABACAttributeValServer) FindResourceValueWithDefinition(ctx context.Context, request *permissionv1.AttributeValueServiceFindResourceValueWithDefinitionRequest) (*permissionv1.AttributeValueServiceFindResourceValueWithDefinitionResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	abacObj, err := a.svc.FindResourceValueWithDefinition(ctx, bizId, request.ResourceId)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceFindResourceValueWithDefinitionResponse{
		Resource: &permissionv1.ResourceObject{
			Id:              abacObj.ID,
			AttributeValues: a.convertToResourceAttributeValues(abacObj.AttrValues),
		},
	}, nil

}

func (a *ABACAttributeValServer) SaveEnvironmentValue(ctx context.Context, request *permissionv1.AttributeValueServiceSaveEnvironmentValueRequest) (*permissionv1.AttributeValueServiceSaveEnvironmentValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	val := domain.AttributeValue{
		AttrDef: a.toDomainAttributeDefinition(request.Value.Definition),
		Value:   request.Value.Value,
	}
	id, err := a.svc.SaveEnvironmentValue(ctx, bizId, val)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceSaveEnvironmentValueResponse{Id: id}, nil
}

func (a *ABACAttributeValServer) DeleteEnvironmentValue(ctx context.Context, request *permissionv1.AttributeValueServiceDeleteEnvironmentValueRequest) (*permissionv1.AttributeValueServiceDeleteEnvironmentValueResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = a.svc.DeleteEnvironmentValue(ctx, bizId, request.Id)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceDeleteEnvironmentValueResponse{}, nil
}

func (a *ABACAttributeValServer) FindEnvironmentValueWithDefinition(ctx context.Context, request *permissionv1.AttributeValueServiceFindEnvironmentValueWithDefinitionRequest) (*permissionv1.AttributeValueServiceFindEnvironmentValueWithDefinitionResponse, error) {
	bizId, err := a.baseServer.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	abacObj, err := a.svc.FindEnvironmentValueWithDefinition(ctx, bizId)
	if err != nil {
		return nil, err
	}
	return &permissionv1.AttributeValueServiceFindEnvironmentValueWithDefinitionResponse{
		Environment: &permissionv1.EnvironmentObject{
			AttributeValues: a.convertToEnvironmentAttributeValues(abacObj.AttrValues),
		},
	}, nil
}

func NewABACAttributeValServer(svc abac.AttributeValueSvc) *ABACAttributeValServer {
	return &ABACAttributeValServer{svc: svc}
}
