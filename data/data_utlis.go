package data

import (
	"github.com/coreos/etcd/client"
	"reflect"
)

func getStructId(structObject reflect.Value) string {
	idProperty := structObject.FieldByName("Id")
	if idProperty == (reflect.Value{}) {
		return ""
	}
	return idProperty.Interface().(string)
}

func isCollection(property reflect.Value) bool {
	return property.Kind() == reflect.Array || property.Kind() == reflect.Slice
}

func (t *DataMapper) isObject(property reflect.Value) bool {
	return property.Kind() == reflect.Slice || isCollection(property)
}

func buildEtcdKey(dirKey string, field reflect.StructField, id interface{}) string {
	return dirKey + "/" + id.(string) + "/" + field.Name
}

func mergeMap(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}

func toStruct(rootKey string, node client.Node, output reflect.Value, outputTemplate interface{}) {
	dataParser := DataParser{dataDirKey: rootKey}
	for _, node := range node.Nodes {
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(outputTemplate, output)
		}
	}
}
