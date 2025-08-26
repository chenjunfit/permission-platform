package evaluator

import (
	"encoding/json"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/errs"
)

type StringEvaluator struct{}

func NewStringEvaluator() *StringEvaluator {
	return &StringEvaluator{}
}

func (s *StringEvaluator) Evaluator(wantVal, actualVal string, op domain.RuleOperator) (bool, error) {
	if isSlice(op) {
		list, err := s.GetSliceData(wantVal)
		if err != nil {
			return false, err
		}
		return sliceEvaluator[string](list, actualVal, op)
	}
	switch op {
	case domain.Equals:
		return wantVal == wantVal, nil
	case domain.NotEquals:
		return wantVal != wantVal, nil
	default:
		return false, errs.ErrUnkonwOperator
	}
}

func (s *StringEvaluator) GetSliceData(wantVal string) (res []string, err error) {
	err = json.Unmarshal([]byte(wantVal), &res)
	if err != nil {
		return nil, err
	}
	return res, err
}
