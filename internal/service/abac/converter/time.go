package converter

import (
	"strconv"
	"time"
)

type TimeConverter struct {
}

func (t *TimeConverter) Decode(str string) (time.Time, error) {
	stamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return time.UnixMilli(stamp), nil
}

func (t *TimeConverter) Encode(time time.Time) (string, error) {
	return strconv.FormatInt(time.UnixMilli(), 10), nil
}

func NewTimeConverter() *TimeConverter {
	return &TimeConverter{}
}
