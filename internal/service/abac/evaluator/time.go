package evaluator

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/abac/converter"
	"strconv"
	"strings"
	"time"
)

const (
	timeType  string = "time"
	dayType          = "day"
	monthType        = "month"
	weekType         = "week"

	hoursInDay    = 24
	minutesInHour = 60
	daysInWeek    = 7
	minWeekday    = 0
	maxWeekday    = 6
	minMonthDay   = 1
	maxMonthDay   = 31
	number2       = 2
)

type timeRule struct {
	Type     string
	Value    string
	Operator string
}

func parseTimeRule(rule string) (*timeRule, error) {
	//@day(9:30)
	rule = strings.TrimPrefix(rule, "@")
	parts := strings.SplitN(rule, "(", number2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rule format %s", rule)
	}
	ruleType := parts[0]
	if !slice.Contains[string]([]string{timeType, dayType, monthType, timeType}, ruleType) {
		return nil, fmt.Errorf("invalid %s rule type  %s", rule, ruleType)
	}
	value := strings.TrimSuffix(parts[1], ")")
	return &timeRule{
		Type:  ruleType,
		Value: value,
	}, nil
}

type TimeEvaluator struct {
	converter converter.Converter[time.Time]
}

func NewTimeEvaluator() *TimeEvaluator {
	return &TimeEvaluator{converter: converter.NewTimeConverter()}
}

func (t *TimeEvaluator) Evaluator(wantVal, actualVal string, op domain.RuleOperator) (bool, error) {
	rule, err := parseTimeRule(wantVal)
	if err != nil {
		return false, err
	}
	actualTime, err := t.converter.Decode(actualVal)
	if err != nil {
		return false, err
	}
	switch rule.Type {
	case timeType:
		return checkExactTime(actualTime, rule, op)
	case dayType:
		return checkDailyTime(actualTime, rule, op)
	case weekType:
		return checkWeeklyTime(actualTime, rule, op)
	case monthType:
		return checkMonthlyTime(actualTime, rule, op)
	default:
		return false, fmt.Errorf("unkown rule type %s", rule.Type)
	}
}
func compareTimes(actualTime, targetTime time.Time, op domain.RuleOperator) (bool, error) {
	switch op {
	case domain.Greater:
		return actualTime.After(targetTime), nil
	case domain.Less:
		return actualTime.Before(targetTime), nil
	case domain.GreaterOrEqual:
		return !actualTime.Before(targetTime), nil
	case domain.LessOrEqual:
		return !actualTime.After(targetTime), nil
	default:
		return false, fmt.Errorf("unkonw operator %s", op)
	}
}
func checkExactTime(actualTime time.Time, wantTimeRule *timeRule, op domain.RuleOperator) (bool, error) {
	targetTime, err := strconv.ParseInt(wantTimeRule.Value, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid time value: %s", wantTimeRule.Value)
	}
	target := time.UnixMilli(targetTime)
	return compareTimes(actualTime, target, op)
}
func checkDailyTime(actualTime time.Time, wantTimeRule *timeRule, op domain.RuleOperator) (bool, error) {
	hour, minute, err := parseTimeOfDay(wantTimeRule.Value)
	if err != nil {
		return false, err
	}
	today := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), hour, minute, 0, 0, actualTime.Location())
	return compareTimes(actualTime, today, op)
}
func parseTimeOfDay(timeStr string) (hour, minute int, err error) {
	parts := strings.SplitN(timeStr, ":", number2)
	if len(parts) != number2 {
		return 0, 0, fmt.Errorf("invalid time format %s", timeStr)
	}
	hour, err = strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour >= hoursInDay {
		return 0, 0, fmt.Errorf("invalid hour format %s", parts[0])
	}
	minute, err = strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute >= minutesInHour {
		return 0, 0, fmt.Errorf("invalid minute: %s", parts[1])
	}
	return hour, minute, nil
}

func checkWeeklyTime(actualTime time.Time, wantTimeRule *timeRule, op domain.RuleOperator) (bool, error) {
	//@week(1,9:30)
	parts := strings.SplitN(wantTimeRule.Value, ",", number2)
	if len(parts) != number2 {
		return false, fmt.Errorf("invalid week time format", wantTimeRule.Value)
	}
	week, err := strconv.Atoi(parts[0])
	if err != nil || week < minWeekday || week > maxWeekday {
		return false, fmt.Errorf("invalid week format", wantTimeRule.Value)

	}
	hour, minute, err := parseTimeOfDay(parts[1])
	if err != nil {
		return false, err

	}
	target := time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), hour, minute, 0, 0, actualTime.Location())
	target = target.AddDate(0, 0, week-int(actualTime.Weekday()))
	return compareTimes(actualTime, target, op)
}
func checkMonthlyTime(actualTime time.Time, rule *timeRule, op domain.RuleOperator) (bool, error) {
	parts := strings.Split(rule.Value, ",")
	if len(parts) != number2 {
		return false, fmt.Errorf("invalid monthly time format: %s", rule.Value)
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil || day < minMonthDay || day > maxMonthDay {
		return false, fmt.Errorf("invalid day of month: %s", parts[0])
	}

	hour, minute, err := parseTimeOfDay(parts[1])
	if err != nil {
		return false, err
	}

	// 只处理当前月的时间
	targetDate := time.Date(actualTime.Year(), actualTime.Month(), day, hour, minute, 0, 0, actualTime.Location())
	return compareTimes(actualTime, targetDate, op)
}
