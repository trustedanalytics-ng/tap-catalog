package data

import (
	"github.com/trustedanalytics/tap-go-common/logger"
	"reflect"
)

var logger = logger_wrapper.InitLogger("template_wrapper")

type DataMapper struct {
}

func (t *DataMapper) ToKeyValue(dirKey string, inputStruct interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	structInputValues := reflect.ValueOf(inputStruct)

	if isCollection(structInputValues) {
		for i := 0; i < structInputValues.Len(); i++ {
			ele := structInputValues.Index(i)
			objectAsMap := t.structToMap(dirKey, ele)
			result = mergeMap(result, objectAsMap)
		}
	} else {
		structAsMap := t.structToMap(dirKey, structInputValues)
		result = mergeMap(result, structAsMap)
	}

	return result
}

func (t *DataMapper) ToKey(prefix string, key string) string {
	return prefix + "/" + key
}

func (t *DataMapper) structToMap(dirKey string, structObject reflect.Value) map[string]interface{} {
	values := make([]interface{}, structObject.NumField())
	result := map[string]interface{}{}

	for i := 0; i < structObject.NumField(); i++ {
		values[i] = structObject.Field(i).Interface()

		key := buildEtcdKey(dirKey, structObject.Type().Field(i), getStructId(structObject))

		if t.isObject(structObject.Field(i)) {
			objectAsMap := t.ToKeyValue(key, values[i])
			result = mergeMap(result, objectAsMap)
		} else {
			result[key] = values[i]
		}
	}

	return result
}
