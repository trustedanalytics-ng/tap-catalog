package data

import (
	"github.com/coreos/etcd/client"
	"github.com/trustedanalytics/tap-catalog/models"
	"reflect"
)

type ApplicationParser struct {
}

func (t *ApplicationParser) ToApplication(rootKey string, dataNode client.Node) (models.Application, error) {
	result := models.Application{}

	reflectResultValues := reflect.ValueOf(&result).Elem()
	dataParser := DataParser{dataDirKey: rootKey}

	for _, node := range dataNode.Nodes {
		logger.Debug("Service Key: ", node.Key)
		logger.Debug("Service Value: ", node.Value)
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.Application{}, reflectResultValues)
		} else {
			//TODO add AuditTrail parser
		}
	}

	return result, nil
}
