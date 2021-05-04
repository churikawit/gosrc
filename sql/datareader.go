package sql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"golang.org/x/text/encoding/charmap"
)

// DataReader : Object
type DataReader struct {
	rows        *sql.Rows
	columnType  []*sql.ColumnType
	vals        []interface{}
	columnCount int
}

// CreateDataReader : get a wrapper for sql.Rows
func CreateDataReader(rows *sql.Rows) *DataReader {
	reader := new(DataReader)
	reader.rows = rows

	var err error
	reader.columnType, err = rows.ColumnTypes()
	if err != nil {
		log.Fatal(err)
	}

	reader.columnCount = len(reader.columnType)
	reader.vals = make([]interface{}, len(reader.columnType))
	for i, ct := range reader.columnType {
		// fmt.Printf("%- 25s: %s\n", (*ct).Name(), (*ColumnType)(ct).DatabaseTypeName2())
		switch (*ct).DatabaseTypeName() {
		case "VARCHAR":
			reader.vals[i] = new(sql.NullString)
		case "CHAR":
			reader.vals[i] = new(sql.NullString)
		case "TEXT":
			reader.vals[i] = new(sql.NullString)
		case "NVARCHAR":
			reader.vals[i] = new(sql.NullString)
		case "DECIMAL":
			reader.vals[i] = new(sql.NullFloat64)
		case "FLOAT":
			reader.vals[i] = new(sql.NullFloat64)
		case "BOOL":
			reader.vals[i] = new(sql.NullBool)
		case "INT":
			reader.vals[i] = new(sql.NullInt32)
		case "INT8":
			reader.vals[i] = new(sql.NullInt32)
		case "BIGINT":
			reader.vals[i] = new(sql.NullInt64)
		case "DATE":
			reader.vals[i] = new(sql.NullTime)
		default:
			reader.vals[i] = new(sql.NullString)
		}
	}
	return reader
}

func (dr *DataReader) Read() bool {

	output := dr.rows.Next()
	if !output {
		return output
	}
	err := dr.rows.Scan(dr.vals...)
	if err != nil {
		log.Fatal(err)
	}
	return output
}

// GetName : return FieldName
func (dr *DataReader) GetName(i int) string {
	if i >= dr.FieldCount() {
		return ""
	}
	return dr.columnType[i].Name()
}

func (dr *DataReader) GetNames() []string {
	c := dr.FieldCount()
	output := make([]string, c)
	for i := 0; i < c; i++ {
		output[i] = dr.columnType[i].Name()
	}
	return output
}

// GetDataTypeName : return {VARCHAR, DECIMAL, TEXT, BOOL, INT, BIGINT, DATE, etc...}
func (dr *DataReader) GetDataTypeName(i int) string {
	if i >= dr.FieldCount() {
		return ""
	}
	return dr.columnType[i].DatabaseTypeName()
}

// GetDataTypeName2 : return {VARCHAR(5), DECIMAL(10,2), TEXT, BOOL, INT, BIGINT, DATE, etc...}
func (dr *DataReader) GetDataTypeName2(i int) string {
	if i >= dr.FieldCount() {
		return ""
	}
	return (*ColumnType)(dr.columnType[i]).DatabaseTypeName2()
}

// GetFieldType : return Type of field
func (dr *DataReader) GetFieldType(i int) reflect.Type {
	return reflect.TypeOf(dr.vals[i])
}

// GetValue : return value
// value with type of {sql.NullString, sql.NullFloat64, sql.NullBool, sql.NullInt32, sql.NullInt64, sql.NullTime}
func (dr *DataReader) GetValue(i int) interface{} {
	return dr.vals[i]
}

// GetValue2 : return value
// value with type of {string, float64, bool, int32, int64, time.Time}
func (dr *DataReader) GetValue2(i int) interface{} {
	if i >= dr.FieldCount() {
		return nil
	}

	// s, ok := (vals[1]).(*string) // interface's type assertion
	switch s2 := (dr.vals[i]).(type) {
	case *sql.NullBool:
		if s2.Valid {
			return s2.Bool
		}
		return nil
	case *sql.NullTime:
		if s2.Valid {
			return s2.Time
		}
		return nil
	case *sql.NullString:
		if s2.Valid {
			return s2.String
			// return AsciiToUtf8(s2.String)
		}
		return nil

	case *sql.NullInt32:
		if s2.Valid {
			return s2.Int32
		}
		return nil

	case *sql.NullInt64:
		if s2.Valid {
			return s2.Int64
		}
		return nil
	case *sql.NullFloat64:
		if s2.Valid {
			return s2.Float64
		}
		return nil
	default:
		return nil
	}
}

func (dr *DataReader) GetValues() []interface{} {
	var output []interface{}
	output = make([]interface{}, dr.FieldCount())
	for i := 0; i < dr.FieldCount(); i++ {
		output[i] = dr.GetValue2(i)
	}
	return output
}

// FieldCount : return number of column
func (dr *DataReader) FieldCount() int {
	return dr.columnCount
}

// Close : close reader
func (dr *DataReader) Close() {
	dr.rows.Close()
}

// IsNull : return True if field is null value
func (dr *DataReader) IsNull(i int) bool {
	if i >= dr.FieldCount() {
		return true
	}

	switch s2 := (dr.vals[i]).(type) {
	case *sql.NullBool:
		if s2.Valid {
			return false
		}
		return true
	case *sql.NullTime:
		if s2.Valid {
			return false
		}
		return true
	case *sql.NullString:
		if s2.Valid {
			return false
		}
		return true

	case *sql.NullInt32:
		if s2.Valid {
			return false
		}
		return true

	case *sql.NullInt64:
		if s2.Valid {
			return false
		}
		return true
	case *sql.NullFloat64:
		if s2.Valid {
			return false
		}
		return true
	default:
		return true
	}
}

func AsciiToUtf8(input string) string {
	dewin874 := charmap.Windows874.NewDecoder()
	output, err := dewin874.String(input)
	if err != nil {
		return ""
	}
	return output
}

// -----------------------------------------------------------------------------------------------------------------------------
// ColumnType : override sql.ColumnType for new function()
type ColumnType sql.ColumnType

// DatabaseTypeName2 : return Datatype with length
func (v *ColumnType) DatabaseTypeName2() string {

	ct := (*(sql.ColumnType))(v)
	output := ""
	// output += fmt.Sprintf("%q:%s", (*ct).Name(), (*ct).DatabaseTypeName())
	output += fmt.Sprintf("%s", (*ct).DatabaseTypeName())
	if l, ok := (*ct).Length(); ok {
		output += fmt.Sprintf("(%d)", l)
	} else if precision, scale, ok := (*ct).DecimalSize(); ok {
		output += fmt.Sprintf("(%d,%d)", precision, scale)
	}
	return output
}
