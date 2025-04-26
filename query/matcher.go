// query/matcher.go
package query

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Matcher struct{}

func NewMatcher() *Matcher {
	return &Matcher{}
}

func (m *Matcher) Matches(data map[string]interface{}, filter map[string]interface{}) bool {
	for key, condition := range filter {
		if IsLogicalOperator(Operator(key)) {
			if !m.evaluateLogicalOperator(Operator(key), data, condition) {
				return false
			}
			continue
		}

		value := getNestedValue(data, key)
		if !m.evaluateCondition(value, condition) {
			return false
		}
	}
	return true
}

func (m *Matcher) evaluateLogicalOperator(operator Operator, data map[string]interface{}, condition interface{}) bool {
	conditions, ok := condition.([]interface{})
	if !ok {
		return false
	}

	switch operator {
	case OpAnd:
		for _, c := range conditions {
			if cond, ok := c.(map[string]interface{}); ok {
				if !m.Matches(data, cond) {
					return false
				}
			}
		}
		return true
	case OpOr:
		for _, c := range conditions {
			if cond, ok := c.(map[string]interface{}); ok {
				if m.Matches(data, cond) {
					return true
				}
			}
		}
		return false
	}
	return false
}

func (m *Matcher) evaluateCondition(value, condition interface{}) bool {
	switch cond := condition.(type) {
	case map[string]interface{}:
		return m.evaluateOperators(value, cond)
	default:
		return reflect.DeepEqual(value, condition)
	}
}

func (m *Matcher) evaluateOperators(value interface{}, operators map[string]interface{}) bool {
	for op, condition := range operators {
		switch Operator(op) {
		case OpEquals:
			if !m.matchEQ(value, condition) {
				return false
			}
		case OpNotEquals:
			if m.matchEQ(value, condition) {
				return false
			}
		case OpGreater:
			if !compareValues(value, condition, ">") {
				return false
			}
		case OpGreaterEqual:
			if !compareValues(value, condition, ">=") {
				return false
			}
		case OpLess:
			if !compareValues(value, condition, "<") {
				return false
			}
		case OpLessEqual:
			if !compareValues(value, condition, "<=") {
				return false
			}
		case OpIn:
			if !containsValue(condition, value) {
				return false
			}
		case OpNotIn:
			if containsValue(condition, value) {
				return false
			}
		case OpExists:
			exists := condition.(bool)
			if exists != (value != nil) {
				return false
			}
		case OpRegex:
			if !matchRegex(value, condition) {
				return false
			}
		}
	}
	return true
}

func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}

		next, ok := current[part].(map[string]interface{})
		if !ok {
			return nil
		}
		current = next
	}
	return nil
}

func compareValues(a, b interface{}, op string) bool {
	aVal, errA := toNumber(a)
	bVal, errB := toNumber(b)

	if errA != nil && errB != nil {
		switch op {
		case ">":
			return aVal > bVal
		case ">=":
			return aVal >= bVal
		case "<":
			return aVal < bVal
		case "<=":
			return aVal <= bVal
		default:
			return false
		}
	}

	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	switch op {
	case ">":
		return aStr > bStr
	case ">=":
		return aStr >= bStr
	case "<":
		return aStr < bStr
	case "<=":
		return aStr <= bStr
	}

	return false
}

func containsValue(list interface{}, value interface{}) bool {
	listValue := reflect.ValueOf(list)
	if listValue.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

func matchRegex(value interface{}, pattern interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	patternStr, ok := pattern.(string)
	if !ok {
		return false
	}
	matched, err := regexp.MatchString(patternStr, str)
	return err == nil && matched
}

func toNumber(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to number", v)
}

func (m *Matcher) matchEQ(value, filterValue interface{}) bool {

	if value == nil || filterValue == nil {
		return value == filterValue
	}

	if v1, err1 := toNumber(value); err1 == nil {
		if v2, err2 := toNumber(filterValue); err2 == nil {
			return v1 == v2
		}
	}

	if str1, ok1 := value.(string); ok1 {
		if str2, ok2 := filterValue.(string); ok2 {
			return str1 == str2
		}
	}

	return reflect.DeepEqual(value, filterValue)
}
