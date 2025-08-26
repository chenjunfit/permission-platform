package evaluator

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/errs"
)

type PolicyRuleEvaluator interface {
	Evaluator(actualVal, wantVal string, op domain.RuleOperator) (bool, error)
}

type Selector interface {
	Select(dataType domain.DataType) (PolicyRuleEvaluator, error)
}

type selector struct {
	checkMap map[domain.DataType]PolicyRuleEvaluator
}

func NewSelector() *selector {
	return &selector{
		checkMap: map[domain.DataType]PolicyRuleEvaluator{
			domain.DataTypeString:   NewStringEvaluator(),
			domain.DataTypeBoolean:  NewBoolEvaluator(),
			domain.DataTypeArray:    NewArrayEvaluator(),
			domain.DataTypeDatetime: NewTimeEvaluator(),
			domain.DataTypeNumber:   NewNumberEvaluator(),
			domain.DataTypeFloat:    NewNumberEvaluator(),
		},
	}
}
func (s *selector) Select(dataType domain.DataType) (PolicyRuleEvaluator, error) {
	evaluator, ok := s.checkMap[dataType]
	if !ok {
		return nil, errs.ErrUnkonwDataType
	}
	return evaluator, nil
}
