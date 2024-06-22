package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// myDB wraps a sql.DB connection pool.
type MyDB struct {
	*sql.DB
}

// NewDB initializes a new database connection.
func NewDB(dataSourceName string) (*MyDB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &MyDB{db}, nil
}

// NamedExec executes a named query with the provided arguments.
func (db *MyDB) NamedExec(query string, arg map[string]interface{}) (sql.Result, error) {
	for k, v := range arg {
		placeholder := fmt.Sprintf(":%s", k)
		value := "NULL"
		if v != nil {
			value = fmt.Sprintf("'%v'", v)
		}
		query = strings.ReplaceAll(query, placeholder, value)
	}
	fmt.Println("Executing Query:", query)
	return db.Exec(query)
}

// In expands slice arguments for SQL IN queries.
func (db *MyDB) In(query string, args ...interface{}) (string, []interface{}, error) {
	var inArgs []interface{}
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Slice {
			placeholders := make([]string, val.Len())
			for i := 0; i < val.Len(); i++ {
				placeholders[i] = "?"
				inArgs = append(inArgs, val.Index(i).Interface())
			}
			query = strings.Replace(query, "?", strings.Join(placeholders, ","), 1)
		} else {
			query = strings.Replace(query, "?", "?", 1)
			inArgs = append(inArgs, arg)
		}
	}
	finalQuery := db.formatQuery(query, inArgs...)
	fmt.Println("Executing Query:", finalQuery)
	return query, inArgs, nil
}

// Begin starts a new database transaction.
func (db *MyDB) Begin() (*sql.Tx, error) {
	return db.DB.Begin()
}

// BulkInsert inserts multiple rows into the specified table.
func (db *MyDB) BulkInsert(table string, data []map[string]interface{}) (sql.Result, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to insert")
	}

	columns := make([]string, 0, len(data[0]))
	values := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)*len(data[0]))

	for k := range data[0] {
		columns = append(columns, k)
	}

	for _, row := range data {
		valuePlaceholders := make([]string, len(row))
		for i, col := range columns {
			valuePlaceholders[i] = "?"
			args = append(args, row[col])
		}
		values = append(values, "("+strings.Join(valuePlaceholders, ",")+")")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, strings.Join(columns, ","), strings.Join(values, ","))
	finalQuery := db.formatQuery(query, args...)
	fmt.Println("Executing Query:", finalQuery)
	return db.Exec(query, args...)
}

// BulkUpdate updates multiple rows in the specified table based on the key column.
func (db *MyDB) BulkUpdate(table string, data []map[string]interface{}, key string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to update")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	columns := make([]string, 0, len(data[0]))

	for k := range data[0] {
		if k != key {
			columns = append(columns, k)
		}
	}

	for _, row := range data {
		var setClauses []string
		var args []interface{}
		for _, col := range columns {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
			args = append(args, row[col])
		}
		args = append(args, row[key])
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table, strings.Join(setClauses, ", "), key)

		finalQuery := db.formatQuery(query, args...)
		fmt.Println("Executing Query:", finalQuery)

		_, err := tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Select executes a query and scans all rows into the destination slice.
func (db *MyDB) Select(dest interface{}, query string, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type().Elem()

	for rows.Next() {
		values := make([]sql.RawBytes, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		elem := reflect.New(destType).Elem()
		for i, v := range values {
			if err := mapToStruct(v, columns[i], elem.Addr().Interface()); err != nil {
				return err
			}
		}

		destValue.Set(reflect.Append(destValue, elem))
	}

	return rows.Err()
}

// Get executes a query and scans the first row into the destination struct.
func (db *MyDB) Get(dest interface{}, query string, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]sql.RawBytes, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return err
	}

	destValue := reflect.ValueOf(dest).Elem()
	for i, v := range values {
		if err := mapToStruct(v, columns[i], destValue.Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}

// TypeParser defines the function signature for custom type parsers.
type TypeParser func([]byte) (interface{}, error)

var typeParsers = make(map[reflect.Type]TypeParser)

// RegisterTypeParser registers a custom parser for a specific type.
func RegisterTypeParser(t reflect.Type, parser TypeParser) {
	typeParsers[t] = parser
}

// mapToStruct maps a single value to the corresponding struct field.
func mapToStruct(value []byte, column string, dest interface{}) error {
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	if destType.Kind() != reflect.Struct {
		return setValue(destValue, string(value))
		//return fmt.Errorf("dest must be a pointer to a struct")
	}
	fieldName := ""
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		dbTag := field.Tag.Get("db")
		if strings.EqualFold(dbTag, column) {
			fieldName = field.Name
			break
		}
	}

	if fieldName == "" {
		return fmt.Errorf("column %s not found in struct", column)
	}

	field := destValue.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	fieldType := field.Type()

	valueStr := string(value)

	switch field.Kind() {
	case reflect.Ptr:
		elemType := fieldType.Elem()
		newValue := reflect.New(elemType).Elem()
		err := setValue(newValue, valueStr)
		if err != nil {
			return err
		}
		field.Set(newValue.Addr())
	default:
		err := setValue(field, valueStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func setValue(field reflect.Value, valueStr string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(valueStr, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(valueStr, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(valueStr, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Slice:
		if field.CanConvert(reflect.TypeOf(json.RawMessage{})) {
			newVal := reflect.New(field.Type()).Elem()
			newVal.Set(reflect.ValueOf(valueStr).Convert(field.Type()))
			field.Set(newVal)
		} else {
			return fmt.Errorf("unsupported slice type: %v", field.Type())
		}
	case reflect.Struct:
		if field.CanConvert(reflect.TypeOf(time.Time{})) {
			if val, err := time.Parse(time.RFC3339, valueStr); err == nil {
				newVal := reflect.New(field.Type()).Elem()
				newVal.Set(reflect.ValueOf(val).Convert(field.Type()))
				field.Set(newVal)
			}
		} else {
			return fmt.Errorf("unsupported struct type: %v", field.Type())
		}

	default:
		return fmt.Errorf("unsupported type: %v", field.Kind())
	}
	return nil
}

// Helper function to format query with args
func (db *MyDB) formatQuery(query string, args ...interface{}) string {
	for _, arg := range args {
		query = strings.Replace(query, "?", fmt.Sprintf("'%v'", arg), 1)
	}
	return query
}
