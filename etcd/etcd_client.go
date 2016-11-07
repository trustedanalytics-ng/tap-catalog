/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package etcd

import (
	"encoding/json"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

type EtcdKVStore interface {
	GetKeyValue(key string) (string, error)
	GetKeyIntoStruct(key string, result interface{}) error
	Set(key string, value interface{}) error
	Update(key string, value interface{}) error
	Delete(key string) error
}

type EtcdConnector struct{}

var logger, _ = commonLogger.InitLogger("etcd")

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
		logger.Errorf("Getinng key: %s, error:", key, err)
		return err
	}
	return json.Unmarshal([]byte(resp.Node.Value), result)
}

func (c *EtcdConnector) GetKeyNodes(key string) (client.Node, error) {
	return getKeyNodes(key, client.GetOptions{Recursive: false})
}

func (c *EtcdConnector) GetKeyNodesRecursively(key string) (client.Node, error) {
	return getKeyNodes(key, client.GetOptions{Recursive: true})
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
		logger.Errorf("Setinng key: %s, error:", key, err)
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

func getKeyNodes(key string, getOptions client.GetOptions) (client.Node, error) {
	logger.Debug("Getting nodes of key:", key)

	kapi, err := getKVApiV2DefaultConnector()

	resultNode := client.Node{}
	if err != nil {
		logger.Error("Can't connect with ETCD:", err)
		return resultNode, err
	}

	resp, err := kapi.Get(context.Background(), key, &getOptions)
	if err != nil {
		logger.Errorf("Getinng key: %s, error:", key, err)
		return resultNode, err
	}

	return *resp.Node, nil
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
