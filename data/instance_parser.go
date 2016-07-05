package data

import (
	"reflect"
	"strings"

	"github.com/coreos/etcd/client"

	"github.com/trustedanalytics/tap-catalog/models"
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
		logger.Debug("Instance Key: ", node.Key)
		logger.Debug("Instance Value: ", node.Value)
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.Instance{}, reflectResultValues)
		} else {
			for _, childNode := range node.Nodes {
				logger.Debug("Instance childNode Key: ", childNode.Key)
				logger.Debug("Instance childNode Value: ", childNode.Value)
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
