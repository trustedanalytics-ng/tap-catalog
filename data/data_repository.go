package data

import (
	"github.com/trustedanalytics/tap-catalog/etcd"
)

type RepositoryConnector struct {
	etcdClient etcd.EtcdConnector
	wrapper    DataMapper
}

func (t *RepositoryConnector) StoreData(keyStore map[string]interface{}) error {
	var err error
	for k, v := range keyStore {
		err = t.etcdClient.Set(k, v)
		//TODO push at once all values from map
		//TODO add error handling

	}

	if err != nil {
		return err
	}

	return nil
}

func (t *RepositoryConnector) GetData(dataType string, id string) (interface{}, error) {
	key := t.wrapper.ToKey(dataType, id)
	node, err := t.etcdClient.GetKeyNodes(key)

	if err != nil {
		return "", err
	}

	return t.wrapper.FromKeyValue(dataType, key, node)
}
