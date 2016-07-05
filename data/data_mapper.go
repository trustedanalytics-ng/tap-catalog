package data

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/logger"
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
			result = MergeMap(result, objectAsMap)
		}
	} else {
		structAsMap := t.structToMap(dirKey, structInputValues)
		result = MergeMap(result, structAsMap)
	}
	return result
}

func (t *DataMapper) ToKeyValueByPatches(mainStructDirKey string, inputStruct interface{}, patches []models.Patch) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, patch := range patches {
		originalField := reflect.ValueOf(inputStruct).FieldByName(patch.Field)
		if originalField.IsValid() {
			newValue, err := UnmarshalJSON(patch.Value, patch.Field, originalField.Type())
			if err != nil {
				return result, err
			}

			receivedElement := reflect.ValueOf(newValue).Elem()
			//todo check operation type!
			if t.isObject(originalField) {
				objectAsMap := t.structToMap(mainStructDirKey+"/"+patch.Field, receivedElement)
				result = MergeMap(result, objectAsMap)
			} else {
				objectAsMap := t.SingleFieldToMap(mainStructDirKey+"/"+patch.Field, receivedElement)
				result = MergeMap(result, objectAsMap)
			}
		} else {
			return result, errors.New("Original field not found: " + patch.Field)
		}
	}
	return result, nil
}

func UnmarshalJSON(value json.RawMessage, entityType string, structType reflect.Type) (interface{}, error) {
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

func (t *DataMapper) ToKey(prefix string, key string) string {
	return prefix + "/" + key
}

func (t *DataMapper) structToMap(dirKey string, structObject reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}
	structObject = unwrapPointer(structObject)
	structId := getStructId(structObject)

	for i := 0; i < structObject.NumField(); i++ {
		objectAsMap := t.SingleFieldToMap(buildEtcdKey(dirKey, structObject.Type().Field(i).Name, structId), structObject.Field(i))
		result = MergeMap(result, objectAsMap)
	}

	return result
}

func (t *DataMapper) SingleFieldToMap(key string, fieldValue reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}
	if t.isObject(fieldValue) {
		objectAsMap := t.ToKeyValue(key, fieldValue.Interface())
		result = MergeMap(result, objectAsMap)
	} else {
		result[key] = fieldValue.Interface()
	}
	return result
}
