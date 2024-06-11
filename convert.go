package invoke

import (
	"fmt"
	"reflect"
	"strconv"
)

// Int2Str converts an integer to a string.
func Int2Str(num int) string {
	return strconv.Itoa(num)
}

// Str2Int converts a string to an integer.
func Str2Int(str string) (int, error) {
	return strconv.Atoi(str)
}

// Str2Int64 converts a string to an int64.
func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// Float2Str converts a float64 to a string.
func Float2Str(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// Str2Float converts a string to a float64.
func Str2Float(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// Bool2Str converts a boolean to a string.
func Bool2Str(b bool) string {
	return strconv.FormatBool(b)
}

// Str2Bool converts a string to a boolean.
func Str2Bool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// Int2Float converts an integer to a float64.
func Int2Float(num int) float64 {
	return float64(num)
}

// Float2Int converts a float64 to an integer.
func Float2Int(num float64) int {
	return int(num)
}

// ToInt converts an interface{} to an integer.
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case string:
		return Str2Int(v)
	case float64:
		return Float2Int(v), nil
	default:
		return 0, fmt.Errorf("unable to convert %v to int", value)
	}
}

// ToString converts an interface{} to a string.
func ToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return Int2Str(v)
	case string:
		return v
	case float64:
		return Float2Str(v)
	case bool:
		return Bool2Str(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// ToFloat converts an interface{} to a float64.
func ToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return Int2Float(v), nil
	case string:
		return Str2Float(v)
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unable to convert %v to float64", value)
	}
}

// ConvertBool converts an interface{} to a boolean.
// Bool is a ValueConverter that converts input values to bools.
//
// The conversion rules are:
//   - booleans are returned unchanged
//   - for integer types,
//     1 is true
//     0 is false,
//     other integers are an error
//   - for strings and []byte, same rules as strconv.ParseBool
//   - all other types are an error
func ToBool(src interface{}) (interface{}, error) {
	switch s := src.(type) {
	case bool:
		return s, nil
	case string:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("ConvertBool: couldn't convert %q into type bool", s)
		}
		return b, nil
	case []byte:
		b, err := strconv.ParseBool(string(s))
		if err != nil {
			return nil, fmt.Errorf("ConvertBool: couldn't convert %q into type bool", s)
		}
		return b, nil
	}

	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := sv.Int()
		if iv == 1 || iv == 0 {
			return iv == 1, nil
		}
		return nil, fmt.Errorf("ConvertBool: couldn't convert %d into type bool", iv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uv := sv.Uint()
		if uv == 1 || uv == 0 {
			return uv == 1, nil
		}
		return nil, fmt.Errorf("ConvertBool: couldn't convert %d into type bool", uv)
	}

	return nil, fmt.Errorf("ConvertBool: couldn't convert %v (%T) into type bool", src, src)
}

// ConvertInt64 converts an interface{} to a int64
func ToInt64(v interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 > (1<<31)-1 || i64 < -(1<<31) {
			return nil, fmt.Errorf("ConvertInt64: value %d overflows int32", v)
		}
		return i64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := rv.Uint()
		if u64 > (1<<31)-1 {
			return nil, fmt.Errorf("ConvertInt64: value %d overflows int32", v)
		}
		return int64(u64), nil
	case reflect.String:
		i, err := strconv.Atoi(rv.String())
		if err != nil {
			return nil, fmt.Errorf("ConvertInt64: value %q can't be converted to int32", v)
		}
		return int64(i), nil
	}
	return nil, fmt.Errorf("ConvertInt64: unsupported value %v (type %T) converting to int32", v, v)
}

// Rune2Str This function converts a Unicode code point into its corresponding string representation,
// handling printable ASCII characters, control characters, and Unicode escape sequences.
func Rune2Str(r rune) string {
	if r >= 0x20 && r < 0x7f {
		return fmt.Sprintf("'%c'", r)
	}
	switch r {
	case 0x07:
		return "'\\a'"
	case 0x08:
		return "'\\b'"
	case 0x0C:
		return "'\\f'"
	case 0x0A:
		return "'\\n'"
	case 0x0D:
		return "'\\r'"
	case 0x09:
		return "'\\t'"
	case 0x0b:
		return "'\\v'"
	case 0x5c:
		return "'\\\\\\'"
	case 0x27:
		return "'\\''"
	case 0x22:
		return "'\\\"'"
	}
	if r < 0x10000 {
		return fmt.Sprintf("\\u%04x", r)
	}
	return fmt.Sprintf("\\U%08x", r)
}
