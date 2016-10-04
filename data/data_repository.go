package data

import (
	"github.com/trustedanalytics/tap-catalog/etcd"
	"reflect"
)

type RepositoryConnector struct {
	etcdClient etcd.EtcdConnector
	mapper     DataMapper
}

func (t *RepositoryConnector) StoreData(keyStore map[string]interface{}) error {
	var err error
	for k, v := range keyStore {
		err = t.etcdClient.Set(k, v, nil, 0)
		if err != nil {
			return err
		}
	}
	return err
}

func (t *RepositoryConnector) UpdateData(updates []PatchSingleUpdate) error {
	var err error
	for _, update := range updates {
		node, err := t.etcdClient.GetKeyNodes(update.Key)
		if err != nil {
			logger.Error("UpdateData in etcd error! Can't get key: ", update.Key, err)
			return err
		}

		err = t.etcdClient.Set(update.Key, update.Value, update.PreviousValue, node.ModifiedIndex)
		if err != nil {
			logger.Error("UpdateData in etcd error! key: ", update.Key, err)
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

	result := []interface{}{}

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
func (t *RepositoryConnector) CreateDirs(org string) error {

	dirs := []string{
		org,
		t.mapper.ToKey(org, Templates),
		t.mapper.ToKey(org, Instances),
		t.mapper.ToKey(org, Applications),
		t.mapper.ToKey(org, Services),
		t.mapper.ToKey(org, Images)}

	for _, dir := range dirs {
		if _, err := t.etcdClient.GetKeyNodes(dir); err != nil {
			err := t.etcdClient.AddOrUpdateDir(dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *RepositoryConnector) IsExistByName(expectedName string, model interface{}, key string) (bool, error) {

	result, err := t.GetListOfData(key, model)
	if err != nil {
		return false, err
	}

	for _, el := range result {

		nameField := reflect.ValueOf(el).FieldByName("Name").String()

		if nameField == expectedName {
			return true, nil
		}
	}

	return false, nil
}
