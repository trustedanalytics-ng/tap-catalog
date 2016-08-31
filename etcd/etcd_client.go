package etcd

import (
	"encoding/json"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/trustedanalytics/tap-go-common/logger"
)

type EtcdKVStore interface {
	GetKeyValue(key string) (string, error)
	GetKeyIntoStruct(key string, result interface{}) error
	Set(key string, value interface{}) error
	Update(key string, value interface{}) error
	Delete(key string) error
}

type EtcdConnector struct{}

var logger = logger_wrapper.InitLogger("etcd")

func (c *EtcdConnector) GetKeyValue(key string) (string, error) {
	logger.Debug("Getting value of key:", key)
	result := ""
	err := c.GetKeyIntoStruct(key, &result)
	return result, err
}

func (c *EtcdConnector) GetKeyIntoStruct(key string, result interface{}) error {
	logger.Debug("Getting value of key:", key)
	kapi, err := getKVApiV2DefaultConnector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	resp, err := kapi.Get(context.Background(), key, nil)

	if err != nil {
		logger.Error("Getinng key value error:", err)
		return err
	}
	err = json.Unmarshal([]byte(resp.Node.Value), result)
	return nil
}

func (c *EtcdConnector) GetKeyNodes(key string) (client.Node, error) {
	logger.Debug("Getting nodes of key:", key)

	kapi, err := getKVApiV2DefaultConnector()

	resultNode := client.Node{}
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return resultNode, err
	}

	resp, err := kapi.Get(context.Background(), key, &client.GetOptions{Recursive: true})

	if err != nil {
		logger.Error("Getinng key value error:", err)
		return resultNode, err
	}

	return *resp.Node, nil
}

func (c *EtcdConnector) Set(key string, value, prevValue interface{}, prevIndex uint64) error {
	logger.Debug("Setting value of key:", key)

	valueByte, err := json.Marshal(value)
	if err != nil {
		logger.Error("Can't marshall etcd key value!:", err)
		return err
	}

	kapi, err := getKVApiV2DefaultConnector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	options := &client.SetOptions{PrevIndex: prevIndex}

	if prevValue != nil {
		prevValueByte, err := json.Marshal(prevValue)
		if err != nil {
			logger.Error("Can't marshall prevValue!:", err)
			return err
		}
		prevValueString := string(prevValueByte)
		if isNotEmptyValue(prevValueString) {
			options.PrevValue = prevValueString
		}
	}

	_, err = kapi.Set(context.Background(), key, string(valueByte), options)
	if err != nil {
		logger.Error("Setting key value error", err)
		return err
	}
	return nil
}

func isNotEmptyValue(value string) bool {
	if value != "" && value != `""` {
		return true
	}
	return false
}

func (c *EtcdConnector) Update(key string, value interface{}) error {
	logger.Debug("Updating value of key:", key)
	valueByte, err := json.Marshal(value)
	if err != nil {
		logger.Error("Can't marshall etcd key value!:", err)
		return err
	}

	kapi, err := getKVApiV2DefaultConnector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Update(context.Background(), key, string(valueByte))
	if err != nil {
		logger.Error("Updating key value error", err)
		return err
	}
	return nil
}

func (c *EtcdConnector) Delete(key string, prevIndex uint64) error {
	options := client.DeleteOptions{Recursive: true, PrevIndex: prevIndex}
	return c.delete(key, &options)
}

func (c *EtcdConnector) DeleteDir(key string) error {
	options := client.DeleteOptions{Recursive: true, Dir: true}
	return c.delete(key, &options)
}

func (c *EtcdConnector) AddOrUpdateDir(key string) error {
	kapi, err := getKVApiV2DefaultConnector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevIgnore})
	if err != nil {
		logger.Error("Setting key value error", err)
		return err
	}
	return nil
}

func (c *EtcdConnector) delete(key string, options *client.DeleteOptions) error {
	logger.Debug("Deleting value of key:", key)
	kapi, err := getKVApiV2DefaultConnector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Delete(context.Background(), key, options)
	if err != nil {
		logger.Error("Getinng key value error:", err)
		return err
	}
	return nil
}

func getKVApiV2DefaultConnector() (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		logger.Error("connection error:", err)
		return nil, err
	}
	return client.NewKeysAPI(c), nil
}
