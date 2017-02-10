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
	"strconv"
	"strings"

	"github.com/nu7hatch/gouuid"

	"github.com/trustedanalytics-ng/tap-catalog/models"
)

const (
	keySeparator = "/"

	idFieldName       = "Id"
	nameFieldName     = "Name"
	bindingsFieldName = "Bindings"
	stateFieldName    = "State"
	classIdFieldName  = "ClassId"

	auditTrailKey = "AuditTrail"
)

var (
	runningStates = []models.InstanceState{models.InstanceStateRunning, models.InstanceStateStopReq}
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

func CheckIfIdFieldIsEmpty(entity interface{}) error {
	idPropertyValue := getStructID(reflect.ValueOf(entity))
	if idPropertyValue != "" {
		return errors.New("Id field has to be empty!")
	}
	return nil
}

func GetEntityKey(organization string, entity string) string {
	dataMapper := DataMapper{}
	org := dataMapper.ToKey("", organization)
	return dataMapper.ToKey(org, entity)
}

func GetFilteredInstances(expectedInstanceType models.InstanceType, expectedClassId string, org string, repositoryApi RepositoryApi) ([]models.Instance, error) {
	filteredInstances := []models.Instance{}

	result, err := repositoryApi.GetListOfData(GetEntityKey(org, Instances), models.Instance{})

	if err != nil {
		return filteredInstances, err
	}

	for _, el := range result {
		instance, ok := el.(models.Instance)
		if !ok {
			return filteredInstances, errors.New("Cannot convert element to models.Instance")
		}

		if instance.Type == expectedInstanceType &&
			(expectedClassId == "" || instance.ClassId == expectedClassId) {
			filteredInstances = append(filteredInstances, instance)
		}
	}
	return filteredInstances, nil
}

func IsInstanceTypeOf(instance models.Instance, instanceType models.InstanceType) bool {
	return instance.Type == instanceType
}

func IsRunningInstance(instance models.Instance) bool {
	return isInstanceInState(instance, runningStates)
}

func isInstanceInState(instance models.Instance, states []models.InstanceState) bool {
	return len(states) > 0 && isStateInSlice(instance.State, states)
}

func isStateInSlice(a models.InstanceState, list []models.InstanceState) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getStructID(structObject reflect.Value) string {
	structObject = unwrapPointer(structObject)
	idProperty := structObject.FieldByName(idFieldName)
	if idProperty == (reflect.Value{}) {
		return ""
	} else {
		return idProperty.Interface().(string)
	}
}

func getOrCreateStructID(structObject reflect.Value) string {
	structId := getStructID(structObject)
	if structId == "" {
		newId, _ := GenerateID()
		idProperty := structObject.FieldByName(idFieldName)
		idProperty.SetString(newId)
		return newId
	} else {
		return structId
	}
}

func GenerateID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func unwrapPointer(structObject reflect.Value) reflect.Value {
	if structObject.Kind() == reflect.Ptr {
		return unwrapPointer(reflect.Indirect(structObject))
	} else {
		return structObject
	}
}

func isCollection(kind reflect.Kind) bool {
	return kind == reflect.Array || kind == reflect.Slice
}

func isObject(property reflect.Value) bool {
	return isCollection(property.Kind()) || property.Kind() == reflect.Struct
}

func isSimpleType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

func buildEtcdKey(dirKey string, fieldName, id string, addIdToKey bool) string {
	if addIdToKey {
		return dirKey + keySeparator + id + keySeparator + fieldName
	} else {
		return dirKey + keySeparator + fieldName
	}
}

func unmarshalJSON(value []byte, fieldName string, structType reflect.Type) (interface{}, error) {
	entity := getNewInstance(fieldName, structType).Interface()
	// arrays/slices saved as marshaled strings should have quotes removed before unmarshalling
	if isCollection(structType.Kind()) {
		if unquoted, err := strconv.Unquote(string(value)); err == nil {
			value = []byte(unquoted)
		}
	}
	err := json.Unmarshal(value, &entity)
	return entity, err
}

func getNodeName(key string) string {
	nodeKeys := strings.Split(key, keySeparator)
	return nodeKeys[len(nodeKeys)-1]
}

func getNewInstance(fieldName string, structType reflect.Type) reflect.Value {
	if reflectType, ok := models.Registry[fieldName]; ok {
		v := reflect.New(reflectType)
		return v
	} else {
		v := reflect.New(structType)
		return v
	}
}

func isStateField(key string) bool {
	return strings.HasSuffix(key, keySeparator+stateFieldName)
}

func getIdFromKey(key, prefix, suffix string) string {
	result := strings.TrimPrefix(key, prefix)
	result = strings.TrimSuffix(result, suffix)
	return strings.Trim(result, keySeparator)
}

func isAuditTrailKey(key string) bool {
	keys := strings.Split(key, "/")
	if len(keys) >= 2 {
		if keys[len(keys)-2] == auditTrailKey {
			return true
		}
	}
	return false
}
