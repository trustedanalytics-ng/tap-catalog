package data

import (
	"github.com/coreos/etcd/client"
	"github.com/trustedanalytics/tap-catalog/api/models"
	"reflect"
)

type ServiceParser struct {
}

func (t *ServiceParser) ToService(rootKey string, dataNode client.Node) (models.Service, error) {
	result := models.Service{}
	result.Plans = []models.ServicePlan{}
	reflectResultValues := reflect.ValueOf(&result).Elem()
	dataParser := DataParser{dataDirKey: rootKey}

	for _, node := range dataNode.Nodes {
		logger.Debug("Service Key: ", node.Key)
		logger.Debug("Service Value: ", node.Value)
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.Service{}, reflectResultValues)
		} else {
			for _, planNode := range node.Nodes {
				plan := models.ServicePlan{}
				t.ToPlan(planNode.Key, *planNode, reflect.ValueOf(&plan).Elem())
				result.Plans = append(result.Plans, plan)
			}
		}
	}

	return result, nil
}

func (t *ServiceParser) ToPlan(rootKey string, node client.Node, output reflect.Value) {
	dataParser := DataParser{dataDirKey: rootKey}
	for _, node := range node.Nodes {
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.ServicePlan{}, output)
		} else {
			//TODO add audittrail and cost parsing
		}
	}
}
