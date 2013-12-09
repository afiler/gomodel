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
	Data interface{}
	Prototype interface{}
	TableName string
	Config
}

func (m Model) Find(q Query) interface{} {
	query, valList := m.BuildQuery(q)
	fmt.Println(query)
	
	db, err := sql.Open(m.DriverName, m.DataSourceName)
	if err != nil { panic(err) }
	defer db.Close()

	rows, err := db.Query(query, valList...)
	return m.loadRows(rows)
	// fmt.Println(m)
}

func (m Model) loadRows(rows *sql.Rows) interface{} {
	columns, err := rows.Columns()
	if err != nil { panic(err) }
	rows.Next()
	return m.loadRow(rows, columns)
}

func (m Model) loadRow(row *sql.Rows, columns []string) interface{} {
	valuePointers := make([]interface{}, len(columns))
	
	obj := reflect.New(reflect.TypeOf(m.Prototype)).Elem()
	
	for i, colName := range columns {
		colValue := reflect.ValueOf(m.Prototype).FieldByName(colName)
		valuePointers[i] = reflect.New(colValue.Type()).Interface()
	}

	row.Scan(valuePointers...)
	
	for i, colName := range columns {
		val := reflect.ValueOf(valuePointers[i]).Elem()
		col := obj.FieldByName(colName)
		col.Set(val)
	}
	
	return obj.Interface()
}

func (m Model) BuildQuery(q Query) (string, []interface{}) {
	columnNames := strings.Join(m.ColumnNames(), ", ")
	where, whereVals := Where(q)
	
	return "select " + columnNames + " from " + m.TableName + " where " + where, whereVals
}

func (m Model) ColumnNames() []string {
	structType := reflect.TypeOf(m.Data).Elem()
	numColumns := structType.NumField()

	columnNames := make([]string, numColumns)
		
	for i := 0; i < numColumns; i++ {
		columnNames[i] = structType.Field(i).Name
	}

	return columnNames
}

func (m Model) DataMap() (map[string]interface{}) {
	dataMap := make(map[string]interface{})

	//structType := reflect.TypeOf(m.Data)
	structVal := reflect.ValueOf(m.Data).Elem()
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
	}

	return strings.Join(whereList, " and "), valList
}

func (m Model) String() string {
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