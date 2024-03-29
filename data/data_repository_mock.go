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
// Automatically generated by MockGen. DO NOT EDIT!
// Source: data/data_repository.go

package data

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/trustedanalytics-ng/tap-catalog/models"
)

// Mock of RepositoryApi interface
type MockRepositoryApi struct {
	ctrl     *gomock.Controller
	recorder *_MockRepositoryApiRecorder
}

// Recorder for MockRepositoryApi (not exported)
type _MockRepositoryApiRecorder struct {
	mock *MockRepositoryApi
}

func NewMockRepositoryApi(ctrl *gomock.Controller) *MockRepositoryApi {
	mock := &MockRepositoryApi{ctrl: ctrl}
	mock.recorder = &_MockRepositoryApiRecorder{mock}
	return mock
}

func (_m *MockRepositoryApi) EXPECT() *_MockRepositoryApiRecorder {
	return _m.recorder
}

func (_m *MockRepositoryApi) CreateData(keyStore map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "CreateData", keyStore)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) CreateData(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateData", arg0)
}

func (_m *MockRepositoryApi) SetData(keyStore map[string]interface{}) error {
	ret := _m.ctrl.Call(_m, "SetData", keyStore)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) SetData(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetData", arg0)
}

func (_m *MockRepositoryApi) UpdateData(updates []PatchSingleUpdate) error {
	ret := _m.ctrl.Call(_m, "UpdateData", updates)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) UpdateData(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateData", arg0)
}

func (_m *MockRepositoryApi) ApplyPatchedValues(patchedKeyValues PatchedKeyValues) error {
	ret := _m.ctrl.Call(_m, "ApplyPatchedValues", patchedKeyValues)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) ApplyPatchedValues(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ApplyPatchedValues", arg0)
}

func (_m *MockRepositoryApi) DeleteData(key string) error {
	ret := _m.ctrl.Call(_m, "DeleteData", key)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) DeleteData(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteData", arg0)
}

func (_m *MockRepositoryApi) CreateDir(key string) error {
	ret := _m.ctrl.Call(_m, "CreateDir", key)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) CreateDir(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDir", arg0)
}

func (_m *MockRepositoryApi) GetLatestIndex(key string) (uint64, error) {
	ret := _m.ctrl.Call(_m, "GetLatestIndex", key)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) GetLatestIndex(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLatestIndex", arg0)
}

func (_m *MockRepositoryApi) GetData(key string, model interface{}) (interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetData", key, model)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) GetData(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetData", arg0, arg1)
}

func (_m *MockRepositoryApi) GetListOfData(key string, model interface{}) ([]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetListOfData", key, model)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) GetListOfData(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetListOfData", arg0, arg1)
}

func (_m *MockRepositoryApi) GetListOfDataFlat(key string, model interface{}) ([]interface{}, error) {
	ret := _m.ctrl.Call(_m, "GetListOfDataFlat", key, model)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) GetListOfDataFlat(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetListOfDataFlat", arg0, arg1)
}

func (_m *MockRepositoryApi) GetDataCounter(key string, model interface{}) (int, error) {
	ret := _m.ctrl.Call(_m, "GetDataCounter", key, model)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) GetDataCounter(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDataCounter", arg0, arg1)
}

func (_m *MockRepositoryApi) CreateDirs(org string) error {
	ret := _m.ctrl.Call(_m, "CreateDirs", org)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryApiRecorder) CreateDirs(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateDirs", arg0)
}

func (_m *MockRepositoryApi) IsExistByName(expectedName string, model interface{}, key string) (bool, error) {
	ret := _m.ctrl.Call(_m, "IsExistByName", expectedName, model, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) IsExistByName(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsExistByName", arg0, arg1, arg2)
}

func (_m *MockRepositoryApi) MonitorObjectsStates(key string, afterIndex uint64) (models.StateChange, error) {
	ret := _m.ctrl.Call(_m, "MonitorObjectsStates", key, afterIndex)
	ret0, _ := ret[0].(models.StateChange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryApiRecorder) MonitorObjectsStates(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MonitorObjectsStates", arg0, arg1)
}
