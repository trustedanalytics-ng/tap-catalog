package data

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/logger"
)

var logger = logger_wrapper.InitLogger("template_wrapper")

type DataMapper struct {
	Username string
}

func (t *DataMapper) ToKeyValue(dirKey string, inputStruct interface{}, isRootElement bool) map[string]interface{} {
	result := map[string]interface{}{}
	structInputValues := reflect.ValueOf(inputStruct)

	if isCollection(structInputValues.Kind()) {
		for i := 0; i < structInputValues.Len(); i++ {
			ele := structInputValues.Index(i)
			objectAsMap := t.structToMap(dirKey, ele, true)
			result = MergeMap(result, objectAsMap)
		}
	} else {
		if structInputValues.Type() == reflect.TypeOf(models.AuditTrail{}) {
			result = MergeMap(result, t.updateAuditTrail(dirKey, false))
		} else {
			structAsMap := t.structToMap(dirKey, structInputValues, isRootElement)
			result = MergeMap(result, structAsMap)
		}
	}
	return result
}

func (t *DataMapper) updateAuditTrail(mainStructDirKey string, isUpdateAction bool) map[string]interface{} {
	result := map[string]interface{}{}
	auditTrail := models.AuditTrail{
		CreatedBy:     t.Username,
		CreatedOn:     time.Now().Unix(),
		LastUpdateBy:  t.Username,
		LastUpdatedOn: time.Now().Unix(),
	}
	valueOfAuditTrial := reflect.ValueOf(auditTrail)
	for i := 0; i < valueOfAuditTrial.NumField(); i++ {
		fieldName := valueOfAuditTrial.Type().Field(i).Name
		if shouldUpdateAuditTrailField(fieldName, isUpdateAction) {
			objectAsMap := t.SingleFieldToMap(buildEtcdKey(mainStructDirKey, fieldName, "", false), valueOfAuditTrial.Field(i), fieldName, "")
			result = MergeMap(result, objectAsMap)
		}
	}
	return result
}

func shouldUpdateAuditTrailField(fieldName string, isUpdateAction bool) bool {
	if isUpdateAction {
		if fieldName == "CreatedOn" || fieldName == "CreatedBy" {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

type PatchSingleUpdate struct {
	Key           string
	Value         interface{}
	PreviousValue interface{}
}

type PatchedKeyValues struct {
	Add    map[string]interface{}
	Update []PatchSingleUpdate
	Delete map[string]interface{}
}

func mapToPatchSingleUpdates(input map[string]interface{}, prevValue interface{}) []PatchSingleUpdate {
	result := []PatchSingleUpdate{}
	for k, v := range input {
		result = append(result, PatchSingleUpdate{
			Key:           k,
			Value:         v,
			PreviousValue: prevValue,
		})
	}
	return result
}

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
					result.Add = MergeMap(result.Add, t.structToMap(mainStructDirKey+"/"+patchFieldName, receivedElement, true))
					err := validatePatch(patchFieldName, patch, false)
					if err != nil {
						return result, err
					}
				} else {
					return result, errors.New("Add operation is allowed only for Collections!")
				}
			} else if patch.Operation == models.OperationUpdate {
				if isObject(originalField) {
					result.Update = append(result.Update, mapToPatchSingleUpdates(t.structToMap(mainStructDirKey+"/"+patchFieldName, receivedElement, isCollection(originalField.Kind())), nil)...)
				} else {
					var receivedPreviousValueInterface interface{}
					if len(patch.PrevValue) > 0 {
						previousValue, err := unmarshalJSON(patch.PrevValue, patchFieldName, originalField.Type())
						receivedPreviousValueInterface = reflect.ValueOf(previousValue).Elem().Interface()
						if err != nil {
							return result, err
						}
					}

					err := validatePatch(patchFieldName, patch, true)
					if err != nil {
						return result, err
					}
					result.Update = append(result.Update, mapToPatchSingleUpdates(t.SingleFieldToMap(mainStructDirKey+"/"+patchFieldName, receivedElement, patchFieldName, ""), receivedPreviousValueInterface)...)
				}
			} else if patch.Operation == models.OperationDelete {
				if isCollection(originalField.Kind()) {
					if structId := getStructId(receivedElement); structId != "" {
						result.Delete[mainStructDirKey+"/"+patchFieldName+"/"+structId] = nil
					} else {
						return result, errors.New("Delete operation required NOT EMPPTY ID field!")
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

	result.Update = append(result.Update, mapToPatchSingleUpdates(t.updateAuditTrail(mainStructDirKey+"/AuditTrail", true), nil)...)
	return result, nil
}

func validatePatch(patchFieldName string, patch models.Patch, isUpdateOp bool) error {
	if isUpdateOp {
		if patchFieldName == idFieldName || patchFieldName == nameFieldName {
			return errors.New("ID and Name fields can not be changed!")
		}
	}
	if patchFieldName == metadataFieldName {
		metadataEntity := models.Metadata{}
		err := json.Unmarshal([]byte(patch.Value), &metadataEntity)
		if err != nil {
			return err
		}
		err = CheckIfDNSLabelCompatible(metadataEntity.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DataMapper) ToKey(prefix string, key string) string {
	return prefix + "/" + key
}

func (t *DataMapper) structToMap(dirKey string, structObject reflect.Value, addIdToKey bool) map[string]interface{} {
	result := map[string]interface{}{}
	structObject = unwrapPointer(structObject)
	structId := getOrCreateStructId(structObject)

	for i := 0; i < structObject.NumField(); i++ {
		fieldName := structObject.Type().Field(i).Name
		objectAsMap := t.SingleFieldToMap(buildEtcdKey(dirKey, fieldName, structId, addIdToKey), structObject.Field(i), fieldName, structId)
		result = MergeMap(result, objectAsMap)
	}
	return result
}

func (t *DataMapper) SingleFieldToMap(key string, fieldValue reflect.Value, fieldName, structId string) map[string]interface{} {
	result := map[string]interface{}{}
	if isObject(fieldValue) {
		objectAsMap := t.ToKeyValue(key, fieldValue.Interface(), false)
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
