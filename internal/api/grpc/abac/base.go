package abac

import (
	"context"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"github.com/permission-dev/internal/domain"
)

type baseServer struct{}

// 从gRPC上下文中获取业务ID
func (s *baseServer) getBizIDFromContext(ctx context.Context) (int64, error) {
	return auth.GetBizIDFromContext(ctx)
}
func (s *baseServer) toProtoAttributeDefinition(definition domain.AttributeDefinition) *permissionv1.AttributeDefinition {
	return &permissionv1.AttributeDefinition{
		Id:             definition.ID,
		Name:           definition.Name,
		Description:    definition.Description,
		DataType:       s.convertToProtoDataType(definition.DataType),
		EntityType:     s.convertToProtoEntityType(definition.EntityType),
		ValidationRule: definition.ValidationRule,
		Ctime:          definition.Ctime,
		Utime:          definition.Utime,
	}
}
func (s *baseServer) toDomainAttributeDefinition(definition *permissionv1.AttributeDefinition) domain.AttributeDefinition {
	if definition == nil {
		return domain.AttributeDefinition{}
	}
	return domain.AttributeDefinition{
		ID:             definition.Id,
		Name:           definition.Name,
		Description:    definition.Description,
		DataType:       s.toDomainDataType(definition.DataType),
		EntityType:     s.toDomainEntityType(definition.EntityType),
		ValidationRule: definition.ValidationRule,
		Ctime:          definition.Ctime,
		Utime:          definition.Utime,
	}
}
func (s *baseServer) toDomainDataType(dataType permissionv1.DataType) domain.DataType {
	switch dataType {
	case permissionv1.DataType_DATA_TYPE_STRING:
		return domain.DataTypeString
	case permissionv1.DataType_DATA_TYPE_NUMBER:
		return domain.DataTypeNumber
	case permissionv1.DataType_DATA_TYPE_BOOLEAN:
		return domain.DataTypeBoolean
	case permissionv1.DataType_DATA_TYPE_FLOAT:
		return domain.DataTypeFloat
	case permissionv1.DataType_DATA_TYPE_DATETIME:
		return domain.DataTypeDatetime
	default:
		return domain.DataType("")
	}

}
func (s *baseServer) toDomainEntityType(entityType permissionv1.EntityType) domain.EntityType {
	switch entityType {
	case permissionv1.EntityType_ENTITY_TYPE_SUBJECT:
		return domain.SubjectTypeEntity
	case permissionv1.EntityType_ENTITY_TYPE_RESOURCE:
		return domain.ResourceTypeEntity
	case permissionv1.EntityType_ENTITY_TYPE_ENVIRONMENT:
		return domain.EnvironmentTypeEntity
	default:
		return ""
	}

}
func (s *baseServer) convertToProtoAttributeDefinition(d domain.AttributeDefinition) *permissionv1.AttributeDefinition {
	return &permissionv1.AttributeDefinition{
		Id:             d.ID,
		Name:           d.Name,
		Description:    d.Description,
		DataType:       s.convertToProtoDataType(d.DataType),
		EntityType:     s.convertToProtoEntityType(d.EntityType),
		ValidationRule: d.ValidationRule,
		Ctime:          d.Ctime,
		Utime:          d.Utime,
	}
}
func (s *baseServer) convertToProtoDataType(d domain.DataType) permissionv1.DataType {
	switch d {
	case domain.DataTypeString:
		return permissionv1.DataType_DATA_TYPE_STRING
	case domain.DataTypeNumber:
		return permissionv1.DataType_DATA_TYPE_NUMBER
	case domain.DataTypeBoolean:
		return permissionv1.DataType_DATA_TYPE_BOOLEAN
	case domain.DataTypeFloat:
		return permissionv1.DataType_DATA_TYPE_FLOAT
	case domain.DataTypeDatetime:
		return permissionv1.DataType_DATA_TYPE_DATETIME
	default:
		return permissionv1.DataType_DATA_TYPE_UNKNOWN
	}
}
func (s *baseServer) convertToProtoEntityType(e domain.EntityType) permissionv1.EntityType {
	switch e {
	case domain.SubjectTypeEntity:
		return permissionv1.EntityType_ENTITY_TYPE_SUBJECT
	case domain.ResourceTypeEntity:
		return permissionv1.EntityType_ENTITY_TYPE_RESOURCE
	case domain.EnvironmentTypeEntity:
		return permissionv1.EntityType_ENTITY_TYPE_ENVIRONMENT
	default:
		return permissionv1.EntityType_ENTITY_TYPE_UNKNOWN
	}
}
func (s *baseServer) convertToProtoBizDefinition(d domain.BizAttrDefinition) *permissionv1.BizDefinition {
	return &permissionv1.BizDefinition{
		SubjectAttrs:     s.convertToProtoAttributeDefinitions(d.SubjectAttrDefs),
		ResourceAttrs:    s.convertToProtoAttributeDefinitions(d.ResourceAttrDefs),
		EnvironmentAttrs: s.convertToProtoAttributeDefinitions(d.EnvironmentAttrDefs),
	}
}

func (s *baseServer) convertToProtoAttributeDefinitions(defs []domain.AttributeDefinition) []*permissionv1.AttributeDefinition {
	res := make([]*permissionv1.AttributeDefinition, 0, len(defs))
	for _, def := range defs {
		res = append(res, s.convertToProtoAttributeDefinition(def))
	}
	return res
}
func (s *baseServer) convertToProtoSubjectAttributeValues(values []domain.AttributeValue) []*permissionv1.SubjectAttributeValue {
	result := make([]*permissionv1.SubjectAttributeValue, 0, len(values))
	for _, v := range values {
		result = append(result, &permissionv1.SubjectAttributeValue{
			Id:         v.ID,
			Definition: s.convertToProtoAttributeDefinition(v.AttrDef),
			Value:      v.Value,
			Ctime:      v.Ctime,
			Utime:      v.Utime,
		})
	}
	return result
}
func (s *baseServer) convertToResourceAttributeValues(values []domain.AttributeValue) []*permissionv1.ResourceAttributeValue {
	result := make([]*permissionv1.ResourceAttributeValue, 0, len(values))
	for _, v := range values {
		result = append(result, &permissionv1.ResourceAttributeValue{
			Id:         v.ID,
			Definition: s.convertToProtoAttributeDefinition(v.AttrDef),
			Value:      v.Value,
			Ctime:      v.Ctime,
			Utime:      v.Utime,
		})
	}
	return result
}
func (s *baseServer) convertToEnvironmentAttributeValues(values []domain.AttributeValue) []*permissionv1.EnvironmentAttributeValue {
	result := make([]*permissionv1.EnvironmentAttributeValue, 0, len(values))
	for _, v := range values {
		result = append(result, &permissionv1.EnvironmentAttributeValue{
			Id:         v.ID,
			Definition: s.convertToProtoAttributeDefinition(v.AttrDef),
			Value:      v.Value,
			Ctime:      v.Ctime,
			Utime:      v.Utime,
		})
	}
	return result
}
func (s *baseServer) convertToDomainPolicy(p *permissionv1.Policy) domain.Policy {
	if p == nil {
		return domain.Policy{}
	}
	return domain.Policy{
		BizID:       p.BizId,
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Status:      s.convertToDomainPolicyStatus(p.Status),
		Permissions: []domain.UserPermission{
			{
				Effect: domain.Effect(p.Effect),
			},
		},

		Rules: s.convertToDomainPolicyRules(p.Rules),
		Ctime: p.Ctime,
		Utime: p.Utime,
	}
}
func (s *baseServer) convertToDomainPolicyStatus(status permissionv1.PolicyStatus) domain.PolicyStatusType {
	switch status {
	case permissionv1.PolicyStatus_POLICY_STATUS_ACTIVE:
		return domain.PolicyStatusActive
	case permissionv1.PolicyStatus_POLICY_STATUS_INACTIVE:
		return domain.PolicyStatusInActive
	default:
		return ""
	}
}
func (s *baseServer) convertToDomainPolicyRules(rules []*permissionv1.PolicyRule) []domain.PolicyRule {
	if rules == nil {
		return nil
	}
	result := make([]domain.PolicyRule, 0, len(rules))
	for _, rule := range rules {
		result = append(result, s.convertToDomainPolicyRule(rule))
	}
	return result
}
func (s *baseServer) convertToDomainPolicyRule(r *permissionv1.PolicyRule) domain.PolicyRule {
	if r == nil {
		return domain.PolicyRule{}
	}
	left := s.convertToDomainPolicyRule(r.LeftRule)
	right := s.convertToDomainPolicyRule(r.RightRule)
	return domain.PolicyRule{
		ID:        r.Id,
		AttrDef:   s.convertToDomainAttributeDefinition(r.AttributeDefinition),
		Value:     r.Value,
		Operator:  s.convertToDomainOperator(r.Operator),
		LeftRule:  &left,
		RightRule: &right,
	}
}
func (s *baseServer) convertToDomainOperator(o permissionv1.RuleOperator) domain.RuleOperator {
	switch o {
	case permissionv1.RuleOperator_RULE_OPERATOR_EQUALS:
		return domain.Equals
	case permissionv1.RuleOperator_RULE_OPERATOR_NOT_EQUALS:
		return domain.NotEquals
	case permissionv1.RuleOperator_RULE_OPERATOR_GREATER:
		return domain.Greater
	case permissionv1.RuleOperator_RULE_OPERATOR_LESS:
		return domain.Less
	case permissionv1.RuleOperator_RULE_OPERATOR_GREATER_OR_EQUAL:
		return domain.GreaterOrEqual
	case permissionv1.RuleOperator_RULE_OPERATOR_LESS_OR_EQUAL:
		return domain.LessOrEqual
	case permissionv1.RuleOperator_RULE_OPERATOR_AND:
		return domain.AND
	case permissionv1.RuleOperator_RULE_OPERATOR_OR:
		return domain.OR
	case permissionv1.RuleOperator_RULE_OPERATOR_IN:
		return domain.IN
	case permissionv1.RuleOperator_RULE_OPERATOR_NOT_IN:
		return domain.NotIn
	case permissionv1.RuleOperator_RULE_OPERATOR_NOT:
		return domain.NOT
	default:
		return domain.RuleOperator("")
	}
}
func (s *baseServer) convertToDomainAttributeDefinition(d *permissionv1.AttributeDefinition) domain.AttributeDefinition {
	if d == nil {
		return domain.AttributeDefinition{}
	}
	return domain.AttributeDefinition{
		ID:             d.Id,
		Name:           d.Name,
		Description:    d.Description,
		DataType:       s.convertToDomainDataType(d.DataType),
		EntityType:     s.convertToDomainEntityType(d.EntityType),
		ValidationRule: d.ValidationRule,
		Ctime:          d.Ctime,
		Utime:          d.Utime,
	}
}
func (s *baseServer) convertToDomainDataType(d permissionv1.DataType) domain.DataType {
	switch d {
	case permissionv1.DataType_DATA_TYPE_STRING:
		return domain.DataTypeString
	case permissionv1.DataType_DATA_TYPE_NUMBER:
		return domain.DataTypeNumber
	case permissionv1.DataType_DATA_TYPE_BOOLEAN:
		return domain.DataTypeBoolean
	case permissionv1.DataType_DATA_TYPE_FLOAT:
		return domain.DataTypeFloat
	case permissionv1.DataType_DATA_TYPE_DATETIME:
		return domain.DataTypeDatetime
	default:
		return domain.DataType("")
	}
}
func (s *baseServer) convertToDomainEntityType(e permissionv1.EntityType) domain.EntityType {
	switch e {
	case permissionv1.EntityType_ENTITY_TYPE_SUBJECT:
		return domain.SubjectTypeEntity
	case permissionv1.EntityType_ENTITY_TYPE_RESOURCE:
		return domain.ResourceTypeEntity
	case permissionv1.EntityType_ENTITY_TYPE_ENVIRONMENT:
		return domain.EnvironmentTypeEntity
	default:
		return ""
	}
}

func (s *baseServer) convertToProtoPolicy(p domain.Policy) *permissionv1.Policy {
	var effect domain.Effect
	if len(p.Permissions) > 0 {
		effect = p.Permissions[0].Effect
	}
	return &permissionv1.Policy{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Status:      s.convertToProtoPolicyStatus(p.Status),
		Effect:      s.convertToProtoEffect(effect),
		Rules:       s.convertToProtoPolicyRules(p.Rules),
		Ctime:       p.Ctime,
		Utime:       p.Utime,
	}
}
func (s *baseServer) convertToProtoPolicyStatus(status domain.PolicyStatusType) permissionv1.PolicyStatus {
	switch status {
	case domain.PolicyStatusActive:
		return permissionv1.PolicyStatus_POLICY_STATUS_ACTIVE
	case domain.PolicyStatusInActive:
		return permissionv1.PolicyStatus_POLICY_STATUS_INACTIVE
	default:
		return permissionv1.PolicyStatus_POLICY_STATUS_UNKNOWN
	}
}

func (s *baseServer) convertToProtoEffect(e domain.Effect) permissionv1.Effect {
	switch e {
	case domain.EffectAllow:
		return permissionv1.Effect_EFFECT_ALLOW
	case domain.EffectDeny:
		return permissionv1.Effect_EFFECT_DENY
	default:
		return permissionv1.Effect_EFFECT_UNKNOWN
	}
}
func (s *baseServer) convertToProtoPolicyRules(rules []domain.PolicyRule) []*permissionv1.PolicyRule {
	if rules == nil {
		return nil
	}
	result := make([]*permissionv1.PolicyRule, 0, len(rules))
	for _, rule := range rules {
		result = append(result, s.convertToProtoPolicyRule(rule))
	}
	return result
}
func (s *baseServer) convertToProtoPolicyRule(r domain.PolicyRule) *permissionv1.PolicyRule {
	return &permissionv1.PolicyRule{
		Id:                  r.ID,
		AttributeDefinition: s.convertToProtoAttributeDefinition(r.AttrDef),
		Value:               r.Value,
		Operator:            s.convertToProtoOperator(r.Operator),
		LeftRule:            s.convertToProtoPolicyRule(r.SafeLeft()),
		RightRule:           s.convertToProtoPolicyRule(r.SafeRight()),
	}
}
func (s *baseServer) convertToProtoOperator(o domain.RuleOperator) permissionv1.RuleOperator {
	switch o {
	case domain.Equals:
		return permissionv1.RuleOperator_RULE_OPERATOR_EQUALS
	case domain.NotEquals:
		return permissionv1.RuleOperator_RULE_OPERATOR_NOT_EQUALS
	case domain.Greater:
		return permissionv1.RuleOperator_RULE_OPERATOR_GREATER
	case domain.Less:
		return permissionv1.RuleOperator_RULE_OPERATOR_LESS
	case domain.GreaterOrEqual:
		return permissionv1.RuleOperator_RULE_OPERATOR_GREATER_OR_EQUAL
	case domain.LessOrEqual:
		return permissionv1.RuleOperator_RULE_OPERATOR_LESS_OR_EQUAL
	case domain.AND:
		return permissionv1.RuleOperator_RULE_OPERATOR_AND
	case domain.OR:
		return permissionv1.RuleOperator_RULE_OPERATOR_OR
	case domain.IN:
		return permissionv1.RuleOperator_RULE_OPERATOR_IN
	case domain.NotIn:
		return permissionv1.RuleOperator_RULE_OPERATOR_NOT_IN
	case domain.NOT:
		return permissionv1.RuleOperator_RULE_OPERATOR_NOT
	default:
		return permissionv1.RuleOperator_RULE_OPERATOR_UNKNOWN
	}
}
func (s *baseServer) convertToDomainEffect(e permissionv1.Effect) domain.Effect {
	switch e {
	case permissionv1.Effect_EFFECT_ALLOW:
		return domain.EffectAllow
	case permissionv1.Effect_EFFECT_DENY:
		return domain.EffectDeny
	default:
		return "" // Use zero value instead of undefined constant
	}
}
func (s *baseServer) convertToProtoPolicies(policies []domain.Policy) []*permissionv1.Policy {
	if policies == nil {
		return nil
	}
	result := make([]*permissionv1.Policy, 0, len(policies))
	for _, policy := range policies {
		result = append(result, s.convertToProtoPolicy(policy))
	}
	return result
}
