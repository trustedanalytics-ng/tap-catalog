package data

import (
	"errors"
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
		err := dataParser.proccesNode(node, reflectResultValues)
		if err != nil {
			logger.Error(err)
			return model, err
		}
	}
	return model, nil
}

func (d *DataParser) proccesNode(node *client.Node, output reflect.Value) error {
	d.dataNode = node
	if !node.Dir {
		logger.Debug("Instance Key: ", node.Key)
		logger.Debug("Instance Value: ", node.Value)

		err := d.parseToStruct(output)
		if err != nil {
			return err
		}
	} else {
		fieldName, fieldType, err := d.getFieldNameAndTypeIfKeyExistInEtcd(output)
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("Dir node - collection field type, fieldName:", fieldName)

		sliceElement := getNewInstance(fieldName, fieldType).Elem()
		slice := reflect.MakeSlice(reflect.SliceOf(sliceElement.Type()), len(node.Nodes), len(node.Nodes))

		for i, objectNode := range node.Nodes {
			objectId := getNodeName(objectNode.Key)
			childDataParser := DataParser{dataDirKey: d.dataDirKey + "/" + fieldName + "/" + objectId}
			sliceElement := slice.Index(i)
			for _, fieldNode := range objectNode.Nodes {
				if err := childDataParser.proccesNode(fieldNode, sliceElement); err != nil {
					return err
				}
			}
		}
		output.FieldByName(fieldName).Set(slice)
	}
	return nil
}

func (d *DataParser) parseToStruct(output reflect.Value) error {
	fieldName, _, err := d.getFieldNameAndTypeIfKeyExistInEtcd(output)
	if err != nil {
		return err
	}

	field := output.FieldByName(fieldName)
	return setValue(field, d.dataNode.Value, fieldName)
}

func (d *DataParser) getFieldNameAndTypeIfKeyExistInEtcd(structValue reflect.Value) (string, reflect.Type, error) {
	for i := 0; i < structValue.NumField(); i++ {
		key := d.mapToEtcdKey(structValue.Type().Field(i))
		if key == d.dataNode.Key {
			return structValue.Type().Field(i).Name, structValue.Type().Field(i).Type, nil
		}
	}
	return "", nil, errors.New("Cant't find any matching field in ETCD for key: " + d.dataNode.Key)
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
