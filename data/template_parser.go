package data

import (
	"github.com/coreos/etcd/client"
	"github.com/trustedanalytics/tap-catalog/api/models"
	"reflect"
)

type TemplateParser struct {
}

func (t *TemplateParser) ToTemplate(rootKey string, dataNode client.Node) (models.Template, error) {
	result := models.Template{}
	reflectResultValues := reflect.ValueOf(&result).Elem()
	dataParser := DataParser{dataDirKey: rootKey}

	for _, node := range dataNode.Nodes {
		dataParser.dataNode = node
		if !node.Dir {
			dataParser.parseToStruct(models.Template{}, reflectResultValues)
			logger.Debug("Template Id - ", result.Id)
			logger.Debug("Templaet State - ", result.State)
		} else {
			//TODO add Audittrail
			//return result, errors.New("Unexpected key found in template Node")
		}
	}

	return result, nil
}
