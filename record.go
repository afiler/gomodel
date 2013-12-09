package gomodel
// 
// import (
// 	"fmt"
// 	"reflect"
// 	"bytes"
// )
// 
// type Record struct {}
// 
// func (r Record) DataMap() (map[string]interface{}) {
// 	dataMap := make(map[string]interface{})
// 
// 	structVal := reflect.ValueOf(r)
// 	structType := structVal.Type()
// 	numField := structType.NumField()
// 
// 	for i := 0; i < numField; i++ {
// 		dataMap[structType.Field(i).Name] = structVal.Field(i).Interface()
// 	}
// 
// 	return dataMap
// }
// 
// func (r Record) String() string {
// 	var buffer bytes.Buffer
// 
// 	for k, v := range r.DataMap() {
// 		buffer.WriteString(fmt.Sprint("[", k, "] ", v, "\n"))
// 	}
// 	
// 	return buffer.String()
// }