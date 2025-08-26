package converter

import "strconv"

type NumberConverter struct {
}

func (n *NumberConverter) Decode(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func (n *NumberConverter) Encode(t int64) (string, error) {
	return strconv.Itoa(int(t)), nil
}

func NewNumberConverter() *NumberConverter {
	return &NumberConverter{}
}
