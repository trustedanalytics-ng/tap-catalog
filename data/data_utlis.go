package data

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/nu7hatch/gouuid"

	"github.com/trustedanalytics/tapng-catalog/models"
)

const idFieldName = "Id"
const nameFieldName = "Name"
const metadataFieldName = "Metadata"

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
	idPropertyValue := getStructId(reflect.ValueOf(entity))
	if idPropertyValue != "" {
		return errors.New("Id field has to be empty!")
	} else {
		return nil
	}
}

func CheckIfDNSLabelCompatible(content string) error {

	const dnsLabelRegexp = "^[A-Za-z_][A-Za-z0-9_]*$"

	ok, _ := regexp.MatchString(dnsLabelRegexp, content)

	if !ok {
		return errors.New(content + " doesn't match DNS label rule: " + dnsLabelRegexp)
	}

	return nil
}

func getStructId(structObject reflect.Value) string {
	structObject = unwrapPointer(structObject)
	idProperty := structObject.FieldByName(idFieldName)
	if idProperty == (reflect.Value{}) {
		return ""
	} else {
		return idProperty.Interface().(string)
	}
}

func getOrCreateStructId(structObject reflect.Value) string {
	structId := getStructId(structObject)
	if structId == "" {
		newId, _ := uuid.NewV4()
		idProperty := structObject.FieldByName(idFieldName)
		idProperty.SetString(newId.String())
		return newId.String()
	} else {
		return structId
	}
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

func buildEtcdKey(dirKey string, fieldName, id string, addIdToKey bool) string {
	if addIdToKey {
		return dirKey + "/" + id + "/" + fieldName
	} else {
		return dirKey + "/" + fieldName
	}
}

func unmarshalJSON(value []byte, fieldName string, structType reflect.Type) (interface{}, error) {
	entity := getNewInstance(fieldName, structType).Interface()
	err := json.Unmarshal(value, &entity)
	return entity, err
}

func getNodeName(key string) string {
	nodeKeys := strings.Split(key, "/")
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
