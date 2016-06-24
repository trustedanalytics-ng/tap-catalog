package etcd

import (
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"encoding/json"
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
	kapi, err := getKVApiV2Connector()
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

	kapi, err := getKVApiV2Connector()

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

func (c *EtcdConnector) Set(key string, value interface{}) error {
	logger.Debug("Setting value of key:", key)
	valueByte, err := json.Marshal(value)
	if err != nil {
		logger.Error("Can't marshall etcd key value!:", err)
		return err
	}

	kapi, err := getKVApiV2Connector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Set(context.Background(), key, string(valueByte), nil)
	if err != nil {
		log.Println("Setting key value error", err)
		return err
	}
	return nil
}

func (c *EtcdConnector) Update(key string, value interface{}) error {
	logger.Debug("Updating value of key:", key)
	valueByte, err := json.Marshal(value)
	if err != nil {
		logger.Error("Can't marshall etcd key value!:", err)
		return err
	}

	kapi, err := getKVApiV2Connector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Update(context.Background(), key, string(valueByte))
	if err != nil {
		log.Println("Updating key value error", err)
		return err
	}
	return nil
}

func (c *EtcdConnector) Delete(key string) error {
	logger.Debug("Deleting value of key:", key)
	kapi, err := getKVApiV2Connector()
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return err
	}

	_, err = kapi.Delete(context.Background(), key, nil)
	if err != nil {
		logger.Error("Getinng key value error:", err)
		return err
	}
	return nil
}

func getKVApiV2Connector() (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Println("connection error:", err)
		return nil, err
	}
	return client.NewKeysAPI(c), nil
}
