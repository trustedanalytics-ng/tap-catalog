package data

import (
	"errors"
	"reflect"

	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/logger"
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
			objectAsMap, _ := t.structToMap(dirKey, ele)
			result = MergeMap(result, objectAsMap)
		}
	} else {
		structAsMap, _ := t.structToMap(dirKey, structInputValues)
		result = MergeMap(result, structAsMap)
	}
	return result
}

type PatchedKeyValues struct {
	Add    map[string]interface{}
	Update map[string]interface{}
	Delete map[string]interface{}
}

func (t *DataMapper) ToKeyValueByPatches(mainStructDirKey string, inputStruct interface{}, patches []models.Patch) (PatchedKeyValues, error) {
	result := PatchedKeyValues{
		Delete: make(map[string]interface{}),
	}

	for _, patch := range patches {
		singlePatchmap := map[string]interface{}{}
		originalField := reflect.ValueOf(inputStruct).FieldByName(patch.Field)
		if originalField.IsValid() {
			newValue, err := unmarshalJSON(patch.Value, patch.Field, originalField.Type())
			if err != nil {
				return result, err
			}

			var objectAsMap map[string]interface{}
			structId := ""
			receivedElement := reflect.ValueOf(newValue).Elem()

			if isObject(originalField) {
				objectAsMap, structId = t.structToMap(mainStructDirKey+"/"+patch.Field, receivedElement)
			} else {
				objectAsMap = t.SingleFieldToMap(mainStructDirKey+"/"+patch.Field, receivedElement)
			}
			singlePatchmap = MergeMap(singlePatchmap, objectAsMap)

			if patch.Operation == models.OperationAdd {
				result.Add = MergeMap(result.Add, singlePatchmap)
			} else if patch.Operation == models.OperationUpdate {
				result.Update = MergeMap(result.Update, singlePatchmap)
			} else if patch.Operation == models.OperationDelete {
				result.Delete[mainStructDirKey+"/"+patch.Field+"/"+structId] = nil
			} else {
				return result, errors.New("Patch operation type unknown: " + string(patch.Operation))
			}
		} else {
			return result, errors.New("Original field not found: " + patch.Field)
		}
	}
	return result, nil
}

func (t *DataMapper) ToKey(prefix string, key string) string {
	return prefix + "/" + key
}

func (t *DataMapper) structToMap(dirKey string, structObject reflect.Value) (map[string]interface{}, string) {
	result := map[string]interface{}{}
	structObject = unwrapPointer(structObject)
	structId := getStructId(structObject)

	for i := 0; i < structObject.NumField(); i++ {
		objectAsMap := t.SingleFieldToMap(buildEtcdKey(dirKey, structObject.Type().Field(i).Name, structId), structObject.Field(i))
		result = MergeMap(result, objectAsMap)
	}

	return result, structId
}

func (t *DataMapper) SingleFieldToMap(key string, fieldValue reflect.Value) map[string]interface{} {
	result := map[string]interface{}{}
	if isObject(fieldValue) {
		objectAsMap := t.ToKeyValue(key, fieldValue.Interface())
		result = MergeMap(result, objectAsMap)
	} else {
		result[key] = fieldValue.Interface()
	}
	return result
}
