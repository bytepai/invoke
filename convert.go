package invoke

import (
	"fmt"
	"strconv"
)

// IntToString converts an integer to a string.
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// StringToInt converts a string to an integer.
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringToInt64 converts a string to an int64.
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// FloatToString converts a float64 to a string.
func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// StringToFloat converts a string to a float64.
func StringToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// BoolToString converts a boolean to a string.
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// StringToBool converts a string to a boolean.
func StringToBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// IntToFloat converts an integer to a float64.
func IntToFloat(num int) float64 {
	return float64(num)
}

// FloatToInt converts a float64 to an integer.
func FloatToInt(num float64) int {
	return int(num)
}

// InterfaceToInt converts an interface{} to an integer.
func InterfaceToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		return StringToInt(v)
	case float64:
		return FloatToInt(v), nil
	default:
		return 0, fmt.Errorf("unable to convert %v to int", value)
	}
}

// InterfaceToString converts an interface{} to a string.
func InterfaceToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return IntToString(v)
	case string:
		return v
	case float64:
		return FloatToString(v)
	case bool:
		return BoolToString(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// InterfaceToFloat converts an interface{} to a float64.
func InterfaceToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return IntToFloat(v), nil
	case string:
		return StringToFloat(v)
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unable to convert %v to float64", value)
	}
}

// InterfaceToBool converts an interface{} to a boolean.
func InterfaceToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return StringToBool(v)
	default:
		return false, fmt.Errorf("unable to convert %v to bool", value)
	}
}
