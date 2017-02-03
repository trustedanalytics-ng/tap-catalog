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
	"fmt"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	key1       = "key1"
	value1     = "value1"
	prevValue1 = "value2"
	prevIndex1 = 45
)

func TestNewEtcdKVStore(t *testing.T) {
	Convey("Test NewEtcdKVStore should return not nil error on bad port and address", t, func() {
		_, err := NewEtcdKVStore("bad_adress")
		So(err, ShouldBeNil)
	})
}

func TestGetKeyValue(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	Convey("Test GetKeyValue provided with proper key", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(createClientResponse(value1), nil)

		result, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("result should be proper", func() {
			So(result, ShouldEqual, value1)
		})
	})

	Convey("Test GetKeyValue in case etcd returns get error", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(nil, fmt.Errorf(""))

		_, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetKeyRawResponse(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	Convey("Test GetKeyRawResponse provided with proper key", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, gomock.Any()).Return(createClientResponse(value1), nil)

		result, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("result should be proper", func() {
			So(result, ShouldEqual, value1)
		})
	})

	Convey("Test GetKeyRawResponse in case etcd returns get error", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, gomock.Any()).Return(nil, fmt.Errorf(""))

		_, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetKeyIntoStruct(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	Convey("Test GetKeyIntoStruct provided with proper key", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(createClientResponse(value1), nil)

		var result string
		err := etcdKVStore.GetKeyIntoStruct(key1, &result)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("result should be proper", func() {
			So(result, ShouldEqual, value1)
		})
	})

	Convey("Test GetKeyIntoStruct in case etcd returns get error", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(createClientResponse(value1), fmt.Errorf(""))

		var result string
		err := etcdKVStore.GetKeyIntoStruct(key1, &result)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetKeyNodes(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	Convey("Test GetKeyNodes provided with proper key", t, func() {
		properGetOptions := client.GetOptions{Recursive: false, Sort: true}
		properResponse := createClientResponse(value1)
		keysAPI.EXPECT().Get(gomock.Any(), key1, &properGetOptions).Return(properResponse, nil)

		result, err := etcdKVStore.GetKeyNodes(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("result should be proper", func() {
			So(result, ShouldResemble, *properResponse.Node)
		})
	})

	Convey("Test GetKeyNodes in case etcd returns get error", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(nil, fmt.Errorf(""))

		_, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetKeyNodesRecursively(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	Convey("Test GetKeyNodesRecursively provided with proper key", t, func() {
		properGetOptions := client.GetOptions{Recursive: true, Sort: true}
		properResponse := createClientResponse(value1)
		keysAPI.EXPECT().Get(gomock.Any(), key1, &properGetOptions).Return(properResponse, nil)

		result, err := etcdKVStore.GetKeyNodesRecursively(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("result should be proper", func() {
			So(result, ShouldResemble, *properResponse.Node)
		})
	})

	Convey("Test GetKeyNodesRecursively in case etcd returns get error", t, func() {
		keysAPI.EXPECT().Get(gomock.Any(), key1, nil).Return(nil, fmt.Errorf(""))

		_, err := etcdKVStore.GetKeyValue(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreate(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	quotedValue1 := fmt.Sprintf("%q", value1)
	properSetOptions := client.SetOptions{PrevExist: client.PrevNoExist}
	Convey("Test Create provided with key and value", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.Create(key1, value1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Create in case etcd returns error", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.Create(key1, value1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestCreateDir(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	properSetOptions := client.SetOptions{PrevExist: client.PrevNoExist, Dir: true}
	Convey("Test CreateDir provided with key", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, "\"\"", &properSetOptions).Return(nil, nil)

		err := etcdKVStore.CreateDir(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test CreateDir in case etcd returns error", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, "\"\"", &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.CreateDir(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddOrUpdate(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	properSetOptions := client.SetOptions{PrevExist: client.PrevIgnore}
	quotedValue1 := fmt.Sprintf("%q", value1)
	Convey("Test AddOrUpdateDir provided with key", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.AddOrUpdate(key1, value1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test AddOrUpdateDir in case etcd returns error", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.AddOrUpdate(key1, value1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUpdate(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	quotedValue1 := fmt.Sprintf("%q", value1)
	quotedPrevValue1 := fmt.Sprintf("%q", prevValue1)
	Convey("Test Update provided with prevValue", t, func() {
		properSetOptions := client.SetOptions{PrevIndex: prevIndex1, PrevExist: client.PrevExist, PrevValue: quotedPrevValue1}
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.Update(key1, value1, prevValue1, prevIndex1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Update provided without prevValue", t, func() {
		properSetOptions := client.SetOptions{PrevIndex: prevIndex1, PrevExist: client.PrevExist}
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.Update(key1, value1, nil, prevIndex1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Update in case etcd returns error", t, func() {
		properSetOptions := client.SetOptions{PrevIndex: prevIndex1, PrevExist: client.PrevExist}
		keysAPI.EXPECT().Set(gomock.Any(), key1, quotedValue1, &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.Update(key1, value1, nil, prevIndex1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDelete(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	properSetOptions := client.DeleteOptions{Recursive: true, PrevIndex: prevIndex1}
	Convey("Test Delete provided with key", t, func() {
		keysAPI.EXPECT().Delete(gomock.Any(), key1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.Delete(key1, prevIndex1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Delete in case etcd returns error", t, func() {
		keysAPI.EXPECT().Delete(gomock.Any(), key1, &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.Delete(key1, prevIndex1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDeleteDir(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	properSetOptions := client.DeleteOptions{Recursive: true, Dir: true}
	Convey("Test DeleteDir provided with key", t, func() {
		keysAPI.EXPECT().Delete(gomock.Any(), key1, &properSetOptions).Return(nil, nil)

		err := etcdKVStore.DeleteDir(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test DeleteDir in case etcd returns error", t, func() {
		keysAPI.EXPECT().Delete(gomock.Any(), key1, &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.DeleteDir(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddOrUpdateDir(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	properSetOptions := client.SetOptions{PrevExist: client.PrevIgnore, Dir: true}
	Convey("Test AddOrUpdateDir provided with key", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, "", &properSetOptions).Return(nil, nil)

		err := etcdKVStore.AddOrUpdateDir(key1)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test AddOrUpdateDir in case etcd returns error", t, func() {
		keysAPI.EXPECT().Set(gomock.Any(), key1, "", &properSetOptions).Return(nil, fmt.Errorf(""))

		err := etcdKVStore.AddOrUpdateDir(key1)

		Convey("err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetLongPollWatcherForKey(t *testing.T) {
	etcdKVStore, keysAPI := prepareEtcdKVStoreAndKeysAPIMock(t)
	watcher := prepareWatcherMock(t)

	monitorSubNodes := true
	var afterIndex uint64 = 10
	properWatcherOptions := client.WatcherOptions{Recursive: monitorSubNodes, AfterIndex: afterIndex}
	Convey(fmt.Sprintf("Test AddOrUpdateDir provided with monitorSubNodes=%v and afterIndex=%v", monitorSubNodes, afterIndex), t, func() {
		keysAPI.EXPECT().Watcher(gomock.Any(), &properWatcherOptions).Return(watcher)

		result, err := etcdKVStore.GetLongPollWatcherForKey(key1, monitorSubNodes, afterIndex)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("result should be proper", func() {
			So(result, ShouldEqual, watcher)
		})
	})

	monitorSubNodes = false
	afterIndex = 0
	properWatcherOptions = client.WatcherOptions{Recursive: monitorSubNodes, AfterIndex: afterIndex}
	Convey(fmt.Sprintf("Test AddOrUpdateDir provided with monitorSubNodes=%v and afterIndex=%v", monitorSubNodes, afterIndex), t, func() {
		keysAPI.EXPECT().Watcher(gomock.Any(), &properWatcherOptions).Return(watcher)

		result, err := etcdKVStore.GetLongPollWatcherForKey(key1, monitorSubNodes, afterIndex)

		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("result should be proper", func() {
			So(result, ShouldEqual, watcher)
		})
	})
}

func createClientResponse(value string) *client.Response {
	return &client.Response{Node: &client.Node{Value: fmt.Sprintf("%q", value)}}
}

func prepareEtcdKVStoreAndKeysAPIMock(t *testing.T) (EtcdKVStore, *client.MockKeysAPI) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	keysAPIMock := client.NewMockKeysAPI(mockCtrl)
	etcdKVStore := &EtcdConnector{keysAPI: keysAPIMock}
	return etcdKVStore, keysAPIMock
}

func prepareWatcherMock(t *testing.T) *client.MockWatcher {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	return client.NewMockWatcher(mockCtrl)
}
