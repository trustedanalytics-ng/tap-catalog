package data

import (
	"errors"
	"reflect"
	"strings"

	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/logger"
)

var logger = logger_wrapper.InitLogger("template_wrapper")

type DataMapper struct {
}

//todo update AuditTrial
func (t *DataMapper) ToKeyValue(dirKey string, inputStruct interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	structInputValues := reflect.ValueOf(inputStruct)

	if isCollection(structInputValues.Kind()) {
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

//todo update AuditTrial
func (t *DataMapper) ToKeyValueByPatches(mainStructDirKey string, inputStruct interface{}, patches []models.Patch) (PatchedKeyValues, error) {
	result := PatchedKeyValues{
		Delete: make(map[string]interface{}),
	}

	for _, patch := range patches {
		patchFieldName := strings.Title(patch.Field)

		if patch.Field == "" {
			return result, errors.New("field value is empty!")
		} else if originalField := reflect.ValueOf(inputStruct).FieldByName(patchFieldName); originalField.IsValid() {
			newValue, err := unmarshalJSON(patch.Value, patchFieldName, originalField.Type())
			if err != nil {
				return result, err
			}

			receivedElement := reflect.ValueOf(newValue).Elem()
			if patch.Operation == models.OperationAdd {
				if isCollection(originalField.Kind()) {
					objectAsMap, _ := t.structToMap(mainStructDirKey+"/"+patchFieldName, receivedElement)
					result.Add = MergeMap(result.Add, objectAsMap)
				} else {
					return result, errors.New("Add operation is allowed only for Collections!")
				}
				//todo we should make possibility to add new object, not only collection!
			} else if patch.Operation == models.OperationUpdate {
				//todo check if object already exist!
				singlePatchMap := map[string]interface{}{}
				objectAsMap := map[string]interface{}{}
				structId := ""
				if isObject(originalField) {
					//todo here we should check if object with specific ID exist
					objectAsMap, structId = t.structToMap(mainStructDirKey+"/"+patchFieldName, receivedElement)
				} else {
					if patchFieldName == idFieldName {
						return result, errors.New("ID field can not be changed!")
					}
					objectAsMap = t.SingleFieldToMap(mainStructDirKey+"/"+patchFieldName, receivedElement, patchFieldName, structId)
				}
				singlePatchMap = MergeMap(singlePatchMap, objectAsMap)
				result.Update = MergeMap(result.Update, singlePatchMap)
			} else if patch.Operation == models.OperationDelete {
				if isCollection(originalField.Kind()) {
					if structId := getStructId(receivedElement); structId != "" {
						result.Delete[mainStructDirKey+"/"+patchFieldName+"/"+structId] = nil
					} else {
						return result, errors.New("Delete operation required not empty Id field!")
					}
				} else {
					return result, errors.New("Delete operation is allowed only for Collections!")
				}
			} else {
				return result, errors.New("Patch operation type unknown: " + string(patch.Operation))
			}
		} else {
			return result, errors.New("Original field not found: " + patchFieldName)
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
	structId := getOrCreateStructId(structObject)

	for i := 0; i < structObject.NumField(); i++ {
		fieldName := structObject.Type().Field(i).Name
		objectAsMap := t.SingleFieldToMap(buildEtcdKey(dirKey, fieldName, structId), structObject.Field(i), fieldName, structId)
		result = MergeMap(result, objectAsMap)
	}

	return result, structId
}

func (t *DataMapper) SingleFieldToMap(key string, fieldValue reflect.Value, fieldName, structId string) map[string]interface{} {
	result := map[string]interface{}{}
	if isObject(fieldValue) {
		objectAsMap := t.ToKeyValue(key, fieldValue.Interface())
		result = MergeMap(result, objectAsMap)
	} else {
		if fieldName == idFieldName {
			result[key] = structId
		} else {
			result[key] = fieldValue.Interface()
		}
	}
	return result
}
