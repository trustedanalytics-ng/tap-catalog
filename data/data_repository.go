package data

import (
	"github.com/trustedanalytics/tap-catalog/etcd"
)

type RepositoryConnector struct {
	etcdClient etcd.EtcdConnector
	mapper     DataMapper
}

func (t *RepositoryConnector) StoreData(keyStore map[string]interface{}) error {
	////TODO for updates  add another methods with SetOptions.PrevValue options. In order to
	var err error
	for k, v := range keyStore {
		err = t.etcdClient.Set(k, v)
		//TODO push at once all values from map
		//TODO add transactions
		//TODO add error handling

	}

	if err != nil {
		return err
	}

	return nil
}

func (t *RepositoryConnector) DeleteData(key string) error {
	return t.etcdClient.DeleteDir(key)
}

func (t *RepositoryConnector) GetData(dataType string, key string) (interface{}, error) {
	node, err := t.etcdClient.GetKeyNodes(key)

	if err != nil {
		return "", err
	}

	return t.mapper.FromKeyValue(dataType, key, node)
}

func (t *RepositoryConnector) GetListOfData(dataType string, key string) ([]interface{}, error) {
	node, err := t.etcdClient.GetKeyNodes(key)

	var result []interface{}

	if err != nil {
		return result, err
	}

	for _, childNode := range node.Nodes {
		elem, err := t.mapper.FromKeyValue(dataType, childNode.Key, *childNode)
		if err != nil {
			return result, err
		}
		result = append(result, elem)
	}

	return result, nil
}
