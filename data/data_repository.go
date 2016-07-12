package data

import (
	"github.com/trustedanalytics/tapng-catalog/etcd"
)

type RepositoryConnector struct {
	etcdClient etcd.EtcdConnector
	mapper     DataMapper
}

func (t *RepositoryConnector) StoreData(keyStore map[string]interface{}) error {
	var err error
	for k, v := range keyStore {
		err = t.etcdClient.Set(k, v, 0)
		if err != nil {
			return err
		}
	}
	return err
}

func (t *RepositoryConnector) UpdateData(keyStore map[string]interface{}) error {
	var err error
	for k, v := range keyStore {
		node, err := t.etcdClient.GetKeyNodes(k)
		if err != nil {
			logger.Error("UpdateData in etcd error! Can't get key: ", k, err)
			return err
		}

		err = t.etcdClient.Set(k, v, node.ModifiedIndex)
		if err != nil {
			logger.Error("UpdateData in etcd error! key: ", k, err)
			return err
		}
	}
	return err
}

func (t *RepositoryConnector) ApplyPatchedValues(patchedKeyValues PatchedKeyValues) error {
	err := t.StoreData(patchedKeyValues.Add)
	if err != nil {
		return err
	}

	err = t.UpdateData(patchedKeyValues.Update)
	if err != nil {
		return err
	}

	for k, _ := range patchedKeyValues.Delete {
		err = t.DeleteData(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *RepositoryConnector) DeleteData(key string) error {
	return t.etcdClient.DeleteDir(key)
}

func (t *RepositoryConnector) GetData(key string, model interface{}) (interface{}, error) {
	node, err := t.etcdClient.GetKeyNodes(key)
	if err != nil {
		return "", err
	}
	return t.mapper.ToModelInstance(key, node, model)
}

func (t *RepositoryConnector) GetListOfData(key string, model interface{}) ([]interface{}, error) {
	node, err := t.etcdClient.GetKeyNodes(key)

	var result []interface{}

	if err != nil {
		return result, err
	}

	for _, childNode := range node.Nodes {
		elem, err := t.mapper.ToModelInstance(childNode.Key, *childNode, model)
		if err != nil {
			return result, err
		}
		result = append(result, elem)
	}
	return result, nil
}
func (t *RepositoryConnector) CreateDirs() error {
	dirs := []string{Templates, Instances, Applications, Services, Images}

	_, err := t.etcdClient.GetKeyNodes(dirs[0])
	if err != nil {
		for _, dir := range dirs {
			err := t.etcdClient.AddOrUpdateDir(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
