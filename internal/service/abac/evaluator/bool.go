package evaluator

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/errs"
	"github.com/permission-dev/internal/service/abac/converter"
)

type BoolEvaluator struct {
	converter converter.Converter[bool]
}

func NewBoolEvaluator() *BoolEvaluator {
	return &BoolEvaluator{converter: converter.NewBoolConverter()}
}

func (b *BoolEvaluator) Evaluator(wantVal, actualVal string, op domain.RuleOperator) (bool, error) {
	boolWantVal, boolActualVal, err := b.getData(wantVal, actualVal)
	if err != nil {
		return false, err
	}
	switch op {
	case domain.Equals:
		return boolWantVal == boolActualVal, nil
	case domain.NotEquals:
		return !boolWantVal == boolActualVal, nil
	default:
		return false, errs.ErrUnkonwOperator
	}
}
func (b *BoolEvaluator) getData(wantVal, actualVal string) (convWantVal, convActualVal bool, err error) {
	convWantVal, err = b.converter.Decode(wantVal)
	if err != nil {
		return false, false, err
	}
	convActualVal, err = b.converter.Decode(actualVal)
	return convWantVal, convActualVal, err
}
