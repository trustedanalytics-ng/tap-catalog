package data

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/trustedanalytics/tap-catalog/models"
)

func MergeMap(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}

func getStructId(structObject reflect.Value) string {
	idProperty := structObject.FieldByName("Id")
	if idProperty == (reflect.Value{}) {
		return ""
	}
	return idProperty.Interface().(string)
}

func unwrapPointer(structObject reflect.Value) reflect.Value {
	if structObject.Kind() == reflect.Ptr {
		return unwrapPointer(reflect.Indirect(structObject))
	} else {
		return structObject
	}
}

func isCollection(property reflect.Value) bool {
	return property.Kind() == reflect.Array || property.Kind() == reflect.Slice
}

func isObject(property reflect.Value) bool {
	return property.Kind() == reflect.Slice || isCollection(property)
}

func buildEtcdKey(dirKey string, fieldName, id string) string {
	return dirKey + "/" + id + "/" + fieldName
}

func unmarshalJSON(value []byte, entityType string, structType reflect.Type) (interface{}, error) {
	if t, ok := models.Registry[entityType]; ok {
		v := reflect.New(t).Interface()
		err := json.Unmarshal(value, &v)
		return v, err
	} else {
		v := reflect.New(structType).Interface()
		err := json.Unmarshal(value, &v)
		return v, err
	}
}

func getNodeName(key string) string {
	nodeKeys := strings.Split(key, "/")
	return nodeKeys[len(nodeKeys)-1]
}

func getNewInstance(fieldName string, structType reflect.Type) reflect.Value {
	if t, ok := models.Registry[fieldName]; ok {
		v := reflect.New(t)
		return v
	} else {
		v := reflect.New(structType)
		return v
	}
}
