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
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-go-common/util"
)

type EtcdKVStore interface {
	Connect() error
	GetKeyValue(key string) (string, error)
	GetKeyIntoStruct(key string, result interface{}) error
	GetKeyRawResponse(key string) (*client.Response, error)
	GetKeyNodes(key string) (client.Node, error)
	GetKeyNodesRecursively(key string) (client.Node, error)
	Create(key string, value interface{}) error
	CreateDir(key string) error
	AddOrUpdate(key string, value interface{}) error
	AddOrUpdateDir(key string) error
	Update(key string, value, prevValue interface{}, prevIndex uint64) error
	Delete(key string, prevIndex uint64) error
	DeleteDir(key string) error
	GetLongPollWatcherForKey(key string, monitorSubNodes bool, afterIndex uint64) (client.Watcher, error)
}

type EtcdConnector struct {
	address string
	port    int
	keysAPI client.KeysAPI
}

var logger, _ = commonLogger.InitLogger("etcd")

func NewEtcdKVStore(address string, port int) (EtcdKVStore, error) {
	res := &EtcdConnector{address: address, port: port}
	err := res.Connect()
	return res, err
}

func (c *EtcdConnector) Connect() error {
	headerTimeoutFromEnv, _ := util.GetInt64EnvValueOrDefault(EtcdConnectionHeaderTimeout, EtcdConnectionHeaderTimeoutDefault)
	headerTimeout := time.Duration(headerTimeoutFromEnv) * time.Millisecond

	cfg := client.Config{
		Endpoints:               []string{fmt.Sprintf("https://%s:%d", c.address, c.port)},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: headerTimeout,
	}
	newClient, err := client.New(cfg)
	if err != nil {
		err := fmt.Errorf("connection error: %v", err)
		return err
	}
	c.keysAPI = client.NewKeysAPI(newClient)
	return nil
}

func (c *EtcdConnector) GetKeyValue(key string) (string, error) {
	logger.Debug("Getting value of key:", key)
	result := ""
	err := c.GetKeyIntoStruct(key, &result)
	return result, err
}

func (c *EtcdConnector) GetKeyIntoStruct(key string, result interface{}) error {
	logger.Debug("Getting value of key:", key)

	resp, err := c.keysAPI.Get(context.Background(), key, nil)
	if err != nil {
		return fmt.Errorf("getting key %q error: %v", key, err)
	}
	return json.Unmarshal([]byte(resp.Node.Value), result)
}

func (c *EtcdConnector) GetKeyRawResponse(key string) (*client.Response, error) {
	options := client.GetOptions{Recursive: false, Sort: true}
	return c.keysAPI.Get(context.Background(), key, &options)
}

func (c *EtcdConnector) GetKeyNodes(key string) (client.Node, error) {
	return c.getKeyNodes(key, client.GetOptions{Recursive: false, Sort: true})
}

func (c *EtcdConnector) GetKeyNodesRecursively(key string) (client.Node, error) {
	return c.getKeyNodes(key, client.GetOptions{Recursive: true, Sort: true})
}

func (c *EtcdConnector) Create(key string, value interface{}) error {
	logger.Debug("Creating value of key: ", key)

	options := &client.SetOptions{PrevExist: client.PrevNoExist}

	return c.set(key, value, options)
}

func (c *EtcdConnector) CreateDir(key string) error {
	logger.Debug("Creating value of key: ", key)

	options := &client.SetOptions{PrevExist: client.PrevNoExist, Dir: true}

	return c.set(key, "", options)
}

func (c *EtcdConnector) AddOrUpdate(key string, value interface{}) error {
	logger.Debug("Setting value of key: ", key)

	options := &client.SetOptions{PrevExist: client.PrevIgnore}

	return c.set(key, value, options)
}

func (c *EtcdConnector) Update(key string, value, prevValue interface{}, prevIndex uint64) error {
	logger.Debug("Updating value of key: ", key)

	options := &client.SetOptions{PrevIndex: prevIndex, PrevExist: client.PrevExist}

	if err := addPrevValueToOptions(prevValue, options); err != nil {
		return err
	}

	return c.set(key, value, options)
}

func addPrevValueToOptions(prevValue interface{}, options *client.SetOptions) error {
	if prevValue != nil {
		prevValueByte, err := json.Marshal(prevValue)
		if err != nil {
			err = fmt.Errorf("cannot marshal prevValue: %v", err)
			return err
		}
		prevValueString := string(prevValueByte)
		if isNotEmptyValue(prevValueString) {
			options.PrevValue = prevValueString
		}
	}
	return nil
}

func (c *EtcdConnector) set(key string, value interface{}, options *client.SetOptions) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		err = fmt.Errorf("cannot marshal etcd key value: %v", err)
		return err
	}

	_, err = c.keysAPI.Set(context.Background(), key, string(valueByte), options)
	if err != nil {
		err = fmt.Errorf("setting key %s error: %v", key, err)
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

func (c *EtcdConnector) Delete(key string, prevIndex uint64) error {
	options := client.DeleteOptions{Recursive: true, PrevIndex: prevIndex}
	return c.delete(key, &options)
}

func (c *EtcdConnector) DeleteDir(key string) error {
	options := client.DeleteOptions{Recursive: true, Dir: true}
	return c.delete(key, &options)
}

func (c *EtcdConnector) AddOrUpdateDir(key string) error {
	logger.Debugf("Adding or updating directory of key %s", key)

	if _, err := c.keysAPI.Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevIgnore}); err != nil {
		return fmt.Errorf("setting key value error: %v", err)
	}
	return nil
}

func (c *EtcdConnector) delete(key string, options *client.DeleteOptions) error {
	logger.Debug("Deleting value of key:", key)

	_, err := c.keysAPI.Delete(context.Background(), key, options)
	if err != nil {
		return fmt.Errorf("getting key value error: %v", err)
	}
	return nil
}

func (c *EtcdConnector) getKeyNodes(key string, getOptions client.GetOptions) (client.Node, error) {
	logger.Debug("Getting nodes of key:", key)

	resultNode := client.Node{}
	resp, err := c.keysAPI.Get(context.Background(), key, &getOptions)
	if err != nil {
		return resultNode, fmt.Errorf("getting key %q error: %v", key, err)
	}

	return *resp.Node, nil
}

func (c *EtcdConnector) GetLongPollWatcherForKey(key string, monitorSubNodes bool, afterIndex uint64) (client.Watcher, error) {
	logger.Debug("Long pulling for key:", key)

	opts := client.WatcherOptions{
		Recursive:  monitorSubNodes,
		AfterIndex: afterIndex, //0 is from currentTime, 1 from the beginning
	}
	return c.keysAPI.Watcher(key, &opts), nil
}
