package gomodel

import (
	"fmt"
	"reflect"
	"strings"
	"bytes"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	Prototype interface{}
	Config
}

func (m *Model) Find(i interface{}) interface{} {
	m.recordOptions()
	
	switch i := i.(type) {
	default:
		panic(fmt.Sprintf("Unexpected argument type to Find: %T", i))
	case string:
		return m.FindByString(i)
	case []interface{}:
		return m.FindBy(i)
	case Query:
		return m.FindByQuery(i)
	}
}

func (m *Model) FindOne(i interface{}) Record {
	return m.FindByQuery(i.(Query))
}

func (m *Model) FindBy(q Query) Record {
	return RecordBase{}
}

func (m *Model) FindByString(q Query) Record {
	return RecordBase{}
}

func (m *Model) FindByQuery(q Query) Record {
	query, valList := m.BuildQuery(q)
	fmt.Println(query)
	
	db, err := sql.Open(m.DriverName, m.DataSourceName)
	if err != nil { panic(err) }
	defer db.Close()

	rows, err := db.Query(query, valList...)
	if err != nil { panic(err) }

	return m.loadRows(rows)
}

func (m *Model) loadRows(rows *sql.Rows) Record {
	columns, err := rows.Columns()
	if err != nil { panic(err) }
	rows.Next()
	return m.loadRow(rows, columns)
}

func (m *Model) loadRow(row *sql.Rows, columns []string) Record {
	valuePointers := make([]interface{}, len(columns))
	
	t := reflect.TypeOf(m.Prototype)
	obj := reflect.New(t).Elem()
	v := reflect.ValueOf(m.Prototype)
	fmt.Println("t: ", reflect.TypeOf(t))
	
	for i, colName := range columns {
		field := v.FieldByName(colName)
		valuePointers[i] = reflect.New(field.Type()).Interface()
	}

	row.Scan(valuePointers...)
	
	for i, colName := range columns {
		val := reflect.ValueOf(valuePointers[i]).Elem()
		col := obj.FieldByName(colName)
		col.Set(val)
	}

	return obj.Interface().(Record)
}

func (m *Model) BuildQuery(q Query) (string, []interface{}) {
	columnNames := strings.Join(m.ColumnNames(), ", ")
	where, whereVals := Where(q)
	
	return "select " + columnNames + " from " + sqlIdent(m.TableName()) + " where " + where, whereVals
}

func (m *Model) TableName() string {
 	return m.recordOptions().Get("table")
}

func (m *Model) recordOptions() *reflect.StructTag {
	structType := reflect.TypeOf(m.Prototype)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		if field.Anonymous && field.Name == "RecordBase" {
			return &field.Tag
		}
	}
	
	return nil
}

func (m *Model) ColumnNames() []string {
	structType := reflect.TypeOf(m.Prototype)
	numColumns := structType.NumField()

	columnNames := make([]string, numColumns)
	fieldCount := 0
		
	for i := 0; i < numColumns; i++ {

		field := structType.Field(i)
		if !field.Anonymous {
			nameFromTag := field.Tag.Get("column")
			
			if nameFromTag != "" {
				columnNames[fieldCount] = sqlIdent(nameFromTag) + " as " + sqlIdent(field.Name)
			} else {
				columnNames[fieldCount] = sqlIdent(field.Name)
			}
			
			fieldCount++
		}
	}

	return columnNames[:fieldCount]
}

func sqlIdent(identifier string) string {
	// XXX: No built-in quoting in the go sql interface, so...
	return "`"+identifier+"`"
}

func (m *Model) DataMap() (map[string]interface{}) {
	dataMap := make(map[string]interface{})

	//structType := reflect.TypeOf(m.Data)
	structVal := reflect.ValueOf(m.Prototype)
	structType := structVal.Type()
	numField := structType.NumField()

	for i := 0; i < numField; i++ {
		dataMap[structType.Field(i).Name] = structVal.Field(i).Interface()
	}

	return dataMap
}

func Where(q Query) (string, []interface{}) {
	structType := reflect.TypeOf(q)
	structVal := reflect.ValueOf(q)
	numField := structType.NumField()
	
	whereList := make([]string, numField)
	valList := make([]interface{}, numField)

	for i := 0; i < numField; i++ {
		whereList[i] = "`" + structType.Field(i).Name + "` = ?"
		valList[i] = structVal.Field(i).Interface()
		
		fmt.Println(whereList[i], " -> ", valList[i])
	}

	return strings.Join(whereList, " and "), valList
}

func (m *Model) String() string {
	var buffer bytes.Buffer

	for k, v := range m.DataMap() {
		buffer.WriteString(fmt.Sprint("[", k, "] ", v, "\n"))
	}
	
	return buffer.String()
}


func dump(obj interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	numColumns := t.NumField()
	
	fmt.Println("Type: ", t)

	for i := 0; i < numColumns; i++ {
		fmt.Println(t.Field(i).Name, "[", t.Field(i).Type, "]", v.Field(i).Interface())
		_=v
	}
}