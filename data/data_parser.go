package data

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/coreos/etcd/client"
)

type DataParser struct {
	dataNode   *client.Node
	dataDirKey string
}

func (t *DataMapper) ToModelInstance(rootKey string, dataNode client.Node, model interface{}) (interface{}, error) {
	reflectResultValues := reflect.ValueOf(model).Elem()
	dataParser := DataParser{dataDirKey: rootKey}

	for _, node := range dataNode.Nodes {
		err := dataParser.processNode(node, reflectResultValues)
		if err != nil {
			logger.Error(err)
			return model, err
		}
	}
	return model, nil
}

func (d *DataParser) processNode(node *client.Node, output reflect.Value) error {
	d.dataNode = node
	if !node.Dir {
		logger.Debug("Instance Key: ", node.Key)
		logger.Debug("Instance Value: ", node.Value)

		err := d.parseToStruct(output)
		if err != nil {
			return err
		}
	} else {
		structField, err := d.getStructFieldIfKeyExistInEtcd(output)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("Dir node case - collection or struct field type. FieldName:", structField.Name)

		sliceElement := getNewInstance(structField.Name, structField.Type).Elem()
		slice := reflect.MakeSlice(reflect.SliceOf(sliceElement.Type()), len(node.Nodes), len(node.Nodes))

		for i, objectNode := range node.Nodes {
			objectId := getNodeName(objectNode.Key)
			childDataParser := DataParser{dataDirKey: d.dataDirKey + "/" + structField.Name + "/" + objectId}
			sliceElement := slice.Index(i)
			for _, fieldNode := range objectNode.Nodes {
				if err := childDataParser.processNode(fieldNode, sliceElement); err != nil {
					return err
				}
			}
		}

		if isCollection(structField.Type.Kind()) {
			output.FieldByName(structField.Name).Set(slice)
		} else {
			if slice.Len() == 1 {
				output.FieldByName(structField.Name).Set(slice.Index(0))
			} else {
				err = errors.New(fmt.Sprintf("PARSING ERROR! Struct field: %s with type %v was parsed to %d elements instead od 1",
					structField.Name, structField.Type, slice.Len()))
				logger.Error(err)
				return err
			}
		}
	}
	return nil
}

func (d *DataParser) parseToStruct(output reflect.Value) error {
	structField, err := d.getStructFieldIfKeyExistInEtcd(output)
	if err != nil {
		return err
	}

	field := output.FieldByName(structField.Name)
	return setValue(field, d.dataNode.Value, structField.Name)
}

func (d *DataParser) getStructFieldIfKeyExistInEtcd(structValue reflect.Value) (reflect.StructField, error) {
	for i := 0; i < structValue.NumField(); i++ {
		if d.mapToEtcdKey(structValue.Type().Field(i)) == d.dataNode.Key {
			return structValue.Type().Field(i), nil
		}
	}
	return reflect.StructField{}, errors.New("Cant't find any matching field in ETCD for key: " + d.dataNode.Key)
}

func (d *DataParser) mapToEtcdKey(field reflect.StructField) string {
	return d.dataDirKey + "/" + field.Name
}

func setValue(field reflect.Value, value, fieldName string) error {
	if field.CanSet() {
		value, err := unmarshalJSON([]byte(value), fieldName, field.Type())
		if err != nil {
			return err
		} else {
			v := reflect.ValueOf(value)
			v = unwrapPointer(v)
			field.Set(v)
		}
	} else {
		logger.Error("Field can not be set! Field type:", field.Type())
	}
	return nil
}
