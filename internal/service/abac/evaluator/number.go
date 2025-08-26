package evaluator

import (
	"encoding/json"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/abac/converter"
	"strconv"
)

type NumberEvaluator struct {
	converter converter.Converter[int64]
}

func NewNumberEvaluator() *NumberEvaluator {
	return &NumberEvaluator{converter: converter.NewNumberConverter()}
}
func (n *NumberEvaluator) getData(wantVal, actualVal string) (convWantVal, convActualVal int64, err error) {
	convWantVal, err = n.converter.Decode(wantVal)
	if err != nil {
		return 0, 0, err
	}
	convActualVal, err = n.converter.Decode(actualVal)
	return convWantVal, convActualVal, err
}
func (n *NumberEvaluator) getSliceData(wantVal, actualVal string) (res []int64, ans int64, err error) {
	err = json.Unmarshal([]byte(wantVal), &res)
	if err != nil {
		return nil, 0, err
	}
	actualNumberVal, err := strconv.ParseInt(actualVal, 10, 64)
	return res, actualNumberVal, err
}
func (n *NumberEvaluator) Evaluator(wantVal, actualVal string, op domain.RuleOperator) (bool, error) {
	if isSlice(op) {
		list, convActualVal, err := n.getSliceData(wantVal, actualVal)
		if err != nil {
			return false, err
		}
		return sliceEvaluator[int64](list, convActualVal, op)
	}

	convWantVal, convActualVal, err := n.getData(wantVal, actualVal)
	if err != nil {
		return false, err
	}
	return baseEvaluator[int64](convWantVal, convActualVal, op)
}
