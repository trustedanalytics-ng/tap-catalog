package data

import (
	"github.com/coreos/etcd/client"
	"github.com/trustedanalytics/tap-catalog/models"
	"reflect"
	"strings"
)

type InstanceParser struct {
}

func (t *InstanceParser) ToInstance(rootKey string, dataNode client.Node) (models.Instance, error) {
	result := models.Instance{}
	result.Bindings = []models.InstanceBindings{}
	result.Metadata = []models.InstanceMetadata{}
	reflectResultValues := reflect.ValueOf(&result).Elem()
	dataParser := DataParser{dataDirKey: rootKey}

	for _, node := range dataNode.Nodes {
		logger.Debug("Service Key: ", node.Key)
		logger.Debug("Service Value: ", node.Value)
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.Instance{}, reflectResultValues)
		} else {
			for _, childNode := range node.Nodes {
				if isBinding(*childNode) {
					binding := models.InstanceBindings{}
					toStruct(childNode.Key, *childNode, reflect.ValueOf(&binding).Elem(), models.InstanceBindings{})
					result.Bindings = append(result.Bindings, binding)
				} else {
					if isMetadata(*childNode) {
						metadata := models.InstanceMetadata{}
						toStruct(childNode.Key, *childNode, reflect.ValueOf(&metadata).Elem(), models.InstanceMetadata{})
						result.Metadata = append(result.Metadata, metadata)
					}

				}
			}
		}
	}

	return result, nil
}

func isBinding(node client.Node) bool {
	return strings.Contains(node.Key, Bindings)
}

func isMetadata(node client.Node) bool {
	return strings.Contains(node.Key, Metadata)
}
