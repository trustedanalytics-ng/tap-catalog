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

package data

import (
	"errors"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/etcd"
)

const (
	key1           string = "key1"
	data1          string = "data1"
	key2           string = "key2"
	data2          string = "data2"
	key3           string = "key3"
	data3          string = "data3"
	auditTrailPath        = "/" + auditTrailKey + "/data"
	auditTrailData string = "auditTrailData"
	prevData3      string = "prevData3"

	modifiedIndex uint64 = 17
)

func TestCreateData(t *testing.T) {
	repository, etcdClientMock := prepareDataRepositoryWithMocks(t)

	Convey("testing CreateData", t, func() {
		Convey("When all data is created successfuly", func() {
			etcdClientMock.EXPECT().Create(key1, data1).Return(nil)
			etcdClientMock.EXPECT().Create(key2, data2).Return(nil)

			input := map[string]interface{}{
				key1: data1,
				key2: data2,
			}
			err := repository.CreateData(input)
			Convey("response error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("State field should be saved last", func() {
			gomock.InOrder(
				etcdClientMock.EXPECT().Create(key1, data1).Return(nil),
				etcdClientMock.EXPECT().Create(keySeparator+stateFieldName, data3).Return(nil),
			)
			input := map[string]interface{}{
				keySeparator + stateFieldName: data3,
				key1: data1,
			}
			err := repository.CreateData(input)
			Convey("response error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When data is not created successfuly", func() {
			etcdClientMock.EXPECT().Create(key1, data1).Return(errors.New(""))

			input := map[string]interface{}{
				key1: data1,
			}
			err := repository.CreateData(input)
			Convey("response error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCreateDir(t *testing.T) {
	repository, etcdClientMock := prepareDataRepositoryWithMocks(t)
	key := "sampleKey"

	Convey("CreateDir should call etcdClient CreateDir method", t, func() {
		etcdClientMock.EXPECT().CreateDir(key).Return(nil)

		repository.CreateDir(key)
	})
}

func TestApplyPatchedValues(t *testing.T) {
	repository, etcdClientMock := prepareDataRepositoryWithMocks(t)

	Convey("ApplyPatchedValues should execute 3 types of operations: create, update and delete", t, func() {
		input := PatchedKeyValues{}
		input.Add = map[string]interface{}{
			key1: data1,
			key2: data2,
		}
		input.Update = []PatchSingleUpdate{
			PatchSingleUpdate{
				Key:           key3,
				Value:         data3,
				PreviousValue: prevData3,
			},
			PatchSingleUpdate{
				Key:           auditTrailPath,
				Value:         auditTrailData,
				PreviousValue: prevData3,
			},
		}
		input.Delete = map[string]interface{}{
			key3: nil,
		}

		etcdClientMock.EXPECT().GetKeyNodesRecursively(key1).Return(client.Node{Key: ""}, errors.New("key not found"))
		etcdClientMock.EXPECT().Create(key1, data1).Return(nil)

		etcdClientMock.EXPECT().GetKeyNodesRecursively(key2).Return(client.Node{Key: key1, ModifiedIndex: modifiedIndex}, nil)
		etcdClientMock.EXPECT().Update(key2, data2, nil, modifiedIndex).Return(nil)

		etcdClientMock.EXPECT().GetKeyNodesRecursively(key3).Return(client.Node{ModifiedIndex: modifiedIndex}, nil)

		etcdClientMock.EXPECT().GetKeyNodesRecursively(key3).Return(client.Node{ModifiedIndex: modifiedIndex}, nil)
		etcdClientMock.EXPECT().Update(key3, data3, prevData3, modifiedIndex).Return(nil)
		etcdClientMock.EXPECT().GetKeyNodesRecursively(auditTrailPath).Return(client.Node{ModifiedIndex: modifiedIndex}, nil)
		etcdClientMock.EXPECT().Set(auditTrailPath, auditTrailData).Return(nil)

		etcdClientMock.EXPECT().DeleteDir(key3).Return(nil)

		err := repository.ApplyPatchedValues(input)

		Convey("response err should be proper", func() {
			So(err, ShouldBeNil)
		})
	})
}

func prepareDataRepositoryWithMocks(t *testing.T) (RepositoryApi, *etcd.MockEtcdKVStore) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	etcdClientMock := etcd.NewMockEtcdKVStore(mockCtrl)
	return NewRepositoryAPI(etcdClientMock, DataMapper{}), etcdClientMock
}
