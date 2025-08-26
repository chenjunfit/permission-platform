package evaluator

import (
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/errs"
	"github.com/permission-dev/internal/service/abac/converter"
)

type ArrayEvaluator struct {
	converter converter.Converter[[]string]
}

func NewArrayEvaluator() *ArrayEvaluator {
	return &ArrayEvaluator{
		converter: converter.NewArrayConverter(),
	}
}

func (a *ArrayEvaluator) Evaluator(wantVal, actualVal string, op domain.RuleOperator) (bool, error) {
	wantArray, err := a.converter.Decode(wantVal)
	if err != nil {
		return false, err
	}
	actualArray, err := a.converter.Decode(actualVal)
	if err != nil {
		return false, err
	}
	if len(actualArray) == 0 {
		return false, err
	}
	switch op {
	case domain.AnyMatch:
		return a.CheckAnyMatch(wantArray, actualArray), nil
	case domain.AllMatch:
		return a.CheckAllMatch(wantArray, actualArray), nil
	default:
		return false, errs.ErrUnkonwOperator
	}
}

func (a *ArrayEvaluator) CheckAnyMatch(wantArray, actualArray []string) bool {
	wantMap := make(map[string]struct{})
	for index := range wantArray {
		wantMap[wantArray[index]] = struct{}{}
	}
	for index := range actualArray {
		if _, ok := wantMap[actualArray[index]]; ok {
			return true
		}
	}
	return false
}
func (a *ArrayEvaluator) CheckAllMatch(wantArray, actualArray []string) bool {
	wantMap := make(map[string]struct{})
	for index := range wantArray {
		wantMap[wantArray[index]] = struct{}{}
	}
	for index := range actualArray {
		if _, ok := wantMap[actualArray[index]]; !ok {
			return false
		}
	}
	return true
}
