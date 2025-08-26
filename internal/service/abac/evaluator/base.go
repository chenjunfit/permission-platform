package evaluator

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/errs"
)

type Numbered interface {
	~int | ~int8 | ~int16 |
		~int32 | ~int64 | ~uint |
		~uint8 | ~uint16 | ~uint32 |
		~uint64 | ~float32 | ~float64
}

func baseEvaluator[T Numbered](wantedVal, actualVal T, op domain.RuleOperator) (bool, error) {
	switch op {
	case domain.Equals:
		return wantedVal == actualVal, nil
	case domain.NotEquals:
		return wantedVal != actualVal, nil
	case domain.Greater:
		return actualVal > wantedVal, nil
	case domain.Less:
		return actualVal < wantedVal, nil
	case domain.GreaterOrEqual:
		return actualVal >= wantedVal, nil
	case domain.LessOrEqual:
		return actualVal <= wantedVal, nil
	default:
		return false, errs.ErrUnkonwOperator
	}

}

func isSlice(op domain.RuleOperator) bool {
	return op == domain.IN || op == domain.NotIn
}

func sliceEvaluator[T comparable](wantedVal []T, actualVal T, op domain.RuleOperator) (bool, error) {
	switch op {
	case domain.IN:
		for index := range wantedVal {
			if wantedVal[index] == actualVal {
				return true, nil
			}
		}
		return false, nil
	case domain.NotIn:
		for index := range wantedVal {
			if wantedVal[index] == actualVal {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, errs.ErrUnkonwOperator
	}
}
