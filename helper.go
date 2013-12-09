package gomodel

import (
	"fmt"
	"reflect"
	"bytes"
)

type Helper struct {}

func (h Helper) MapObj(obj interface{}) (map[string]interface{}) {
	dataMap := make(map[string]interface{})

	structVal := reflect.ValueOf(obj).Elem()
	structType := structVal.Type()
	numField := structType.NumField()

	for i := 0; i < numField; i++ {
		dataMap[structType.Field(i).Name] = structVal.Field(i).Interface()
	}

	return dataMap
}

func (h Helper) StringObj(obj interface{}) string {
	var buffer bytes.Buffer

	for k, v := range h.MapObj(obj) {
		buffer.WriteString(fmt.Sprint("[", k, "] ", v, "\n"))
	}
	
	return buffer.String()
}
