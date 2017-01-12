/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package data

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/trustedanalytics/tap-catalog/models"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

var logger, _ = commonLogger.InitLogger("DataMapper")

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
			result = MergeMap(result, t.updateAuditTrail(dirKey, false, structInputValues.Interface().(models.AuditTrail)))
		} else {
			structAsMap := t.structToMap(dirKey, structInputValues, isRootElement)
			result = MergeMap(result, structAsMap)
		}
	}
	return result
}

func (t *DataMapper) updateAuditTrail(mainStructDirKey string, isUpdateAction bool, auditTrail models.AuditTrail) map[string]interface{} {
	result := map[string]interface{}{}
	auditTrail.CreatedOn = time.Now().Unix()
	auditTrail.LastUpdatedOn = time.Now().Unix()

	valueOfAuditTrial := reflect.ValueOf(auditTrail)
	for i := 0; i < valueOfAuditTrial.NumField(); i++ {
		fieldName := valueOfAuditTrial.Type().Field(i).Name
		fieldValue := valueOfAuditTrial.Field(i)

		if shouldUpdateAuditTrailField(fieldName, fieldValue, isUpdateAction) {
			objectAsMap := t.SingleFieldToMap(buildEtcdKey(mainStructDirKey, fieldName, "", false), fieldValue, fieldName, "")
			result = MergeMap(result, objectAsMap)
		}
	}
	return result
}

func shouldUpdateAuditTrailField(fieldName string, fieldValue reflect.Value, isUpdateAction bool) bool {
	if isUpdateAction {
		if fieldName == "CreatedOn" || fieldName == "CreatedBy" {
			return false
		} else if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
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

	username := ""
	for _, patch := range patches {
		username = patch.Username

		if err := models.ValidatePatchStructure(patch); err != nil {
			return result, err
		}

		patchFieldName := strings.Title(*patch.Field)
		if originalField := reflect.ValueOf(inputStruct).FieldByName(patchFieldName); originalField.IsValid() {
			newValue, err := unmarshalJSON(*patch.Value, patchFieldName, originalField.Type())
			if err != nil {
				return result, err
			}

			receivedElement := reflect.ValueOf(newValue).Elem()
			if patch.Operation == models.OperationAdd {
				if isCollection(originalField.Kind()) {
					result.Add = MergeMap(result.Add, t.structToMap(mainStructDirKey+keySeparator+patchFieldName, receivedElement, true))
					err := validatePatch(patchFieldName, patch, false)
					if err != nil {
						return result, err
					}
				} else {
					return result, errors.New("Add operation is allowed only for Collections!")
				}
			} else if patch.Operation == models.OperationUpdate {
				if isObject(originalField) {
					result.Update = append(result.Update, mapToPatchSingleUpdates(t.structToMap(mainStructDirKey+keySeparator+patchFieldName, receivedElement, isCollection(originalField.Kind())), nil)...)
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
					result.Update = append(result.Update, mapToPatchSingleUpdates(t.SingleFieldToMap(mainStructDirKey+keySeparator+patchFieldName, receivedElement, patchFieldName, ""), receivedPreviousValueInterface)...)
				}
			} else if patch.Operation == models.OperationDelete {
				if isCollection(originalField.Kind()) {
					if structId := getStructID(receivedElement); structId != "" {
						result.Delete[mainStructDirKey+keySeparator+patchFieldName+keySeparator+structId] = nil
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

	result.Update = append(
		result.Update,
		mapToPatchSingleUpdates(t.updateAuditTrail(mainStructDirKey+keySeparator+"AuditTrail", true, models.AuditTrail{LastUpdateBy: username}), nil)...,
	)
	return result, nil
}

func validatePatch(patchFieldName string, patch models.Patch, isUpdateOp bool) error {
	if isUpdateOp {
		if patchFieldName == idFieldName || patchFieldName == nameFieldName {
			return errors.New("ID and Name fields can not be changed!")
		} else if patchFieldName == classIdFieldName {
			return errors.New("ClassID fields can not be changed!")
		}
	}
	if patchFieldName == bindingsFieldName {
		instanceBinding := models.InstanceBindings{}
		if err := json.Unmarshal([]byte(*patch.Value), &instanceBinding); err != nil {
			return err
		}
		for k, _ := range instanceBinding.Data {
			if err := models.CheckIfMatchingRegexp(k, models.RegexpProperSystemEnvName); err != nil {
				return errors.New("Field: Data has incorrect value: " + k)
			}
		}
	}
	return nil
}

func (t *DataMapper) ToKey(prefix string, key string) string {
	return prefix + keySeparator + key
}

func (t *DataMapper) structToMap(dirKey string, structObject reflect.Value, addIdToKey bool) map[string]interface{} {
	result := map[string]interface{}{}
	structObject = unwrapPointer(structObject)
	structId := getOrCreateStructID(structObject)

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
