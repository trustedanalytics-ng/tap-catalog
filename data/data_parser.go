package data

import (
	"github.com/coreos/etcd/client"
	"reflect"
	"strconv"
)

type DataParser struct {
	dataNode   *client.Node
	dataDirKey string
}

func (d *DataParser) parseToStruct(outputTemplate interface{}, output reflect.Value) error {
	outputStructValues := reflect.ValueOf(outputTemplate)

	if !isCollection(outputStructValues) {
		for i := 0; i < outputStructValues.NumField(); i++ {
			if d.isMatchingEtcdKey(outputStructValues.Type().Field(i)) {
				field := output.FieldByName(outputStructValues.Type().Field(i).Name)
				err := setValue(field, d.dataNode.Value)
				if err != nil {
					return err
				}
			}
		}
	} else {
		//TODO add collection handling
	}

	return nil
}

func (d *DataParser) isMatchingEtcdKey(field reflect.StructField) bool {
	key := d.mapToEtcdKey(field)
	return key == d.dataNode.Key
}

func (d *DataParser) mapToEtcdKey(field reflect.StructField) string {
	return d.dataDirKey + "/" + field.Name
}

func setValue(field reflect.Value, value string) error {
	if field.CanSet() {
		if field.Kind() == reflect.String {
			field.SetString(value)
		}
		if field.Kind() == reflect.Bool {
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		}
		if field.Kind() == reflect.Int {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		}
	}
	return nil
}

func (t *DataMapper) FromKeyValue(dataType string, rootKey string, dataNode client.Node) (interface{}, error) {

	switch dataType {
	case Templates:
		template_parser := TemplateParser{}
		return template_parser.ToTemplate(rootKey, dataNode)
	case Services:
		service_parser := ServiceParser{}
		return service_parser.ToService(rootKey, dataNode)
	}

	//TODO add errror
	return nil, nil
}
