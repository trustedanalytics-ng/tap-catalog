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
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics-ng/tap-catalog/models"
)

const (
	testOrgName = "test-org"
)

func prepareMocks(t *testing.T) *MockRepositoryApi {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repositoryApiMock := NewMockRepositoryApi(mockCtrl)
	return repositoryApiMock
}

func TestGetEntityKey(t *testing.T) {
	Convey("Test GetEntityKey", t, func() {
		Convey("applicationKey should have a proper value", func() {
			applicationKey := GetEntityKey(testOrgName, Applications)
			So(applicationKey, ShouldEqual, fmt.Sprintf("/%s/Applications", testOrgName))
		})
		Convey("serviceKey should have a proper value", func() {
			serviceKey := GetEntityKey(testOrgName, Services)
			So(serviceKey, ShouldEqual, fmt.Sprintf("/%s/Services", testOrgName))
		})
		Convey("instanceKey should have a proper value", func() {
			instanceKey := GetEntityKey(testOrgName, Instances)
			So(instanceKey, ShouldEqual, fmt.Sprintf("/%s/Instances", testOrgName))
		})
	})
}

func TestGetFilteredInstances(t *testing.T) {
	testApplicationInstance := models.Instance{
		Id:   "1",
		Name: "app-1",
		Type: models.InstanceTypeApplication,
	}
	testServiceInstance1 := models.Instance{
		Id:   "2",
		Name: "service-1",
		Type: models.InstanceTypeService,
	}
	testServiceInstance2 := models.Instance{
		Id:   "3",
		Name: "service-2",
		Type: models.InstanceTypeService,
	}

	testAllInstances := make([]interface{}, 3)
	testAllInstances[0] = testApplicationInstance
	testAllInstances[1] = testServiceInstance1
	testAllInstances[2] = testServiceInstance2

	var testAllServicesInstances = []models.Instance{testServiceInstance1, testServiceInstance2}
	var testAllApplicationsInstances = []models.Instance{testApplicationInstance}

	repositoryApiMock := prepareMocks(t)

	Convey("Test GetFilteredInstances", t, func() {
		Convey("GetFilteredInstances should return only instances of services", func() {
			repositoryApiMock.EXPECT().GetListOfData(GetEntityKey(testOrgName, Instances), models.Instance{}).Return(testAllInstances, nil)
			servicesInstances, err := GetFilteredInstances(models.InstanceTypeService, "", testOrgName, repositoryApiMock)
			Convey("So returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("and length of returned array should by 2", func() {
				So(len(servicesInstances), ShouldEqual, 2)
			})
			Convey("and returned array should ressemble to testAllServicesInstances", func() {
				So(servicesInstances, ShouldResemble, testAllServicesInstances)
			})

		})

		Convey("GetFilteredInstances should return only instances of applications", func() {
			repositoryApiMock.EXPECT().GetListOfData(GetEntityKey(testOrgName, Instances), models.Instance{}).Return(testAllInstances, nil)
			servicesInstances, err := GetFilteredInstances(models.InstanceTypeApplication, "", testOrgName, repositoryApiMock)
			Convey("So returned error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("and length of returned array should by 1", func() {
				So(len(servicesInstances), ShouldEqual, 1)
			})
			Convey("and returned array should resemble testAllApplicationsInstances", func() {
				So(servicesInstances, ShouldResemble, testAllApplicationsInstances)
			})
		})
	})
}

func TestGetOrCreateStructID(t *testing.T) {
	Convey("testing getOrCreateStructID", t, func() {
		Convey("When structure ID field is empty", func() {
			entity := &models.Application{Id: ""}
			input := reflect.ValueOf(entity)
			input = unwrapPointer(input)

			id := getOrCreateStructID(input)

			Convey("returned ID should not be empty", func() {
				So(len(id), ShouldBeGreaterThan, 0)
			})
		})

		Convey("When structure ID field is not empty", func() {
			sampleID := "123456"
			entity := &models.Application{Id: sampleID}
			input := reflect.ValueOf(entity)
			input = unwrapPointer(input)

			id := getOrCreateStructID(input)

			Convey("returned ID should be equal to that from the structure", func() {
				So(id, ShouldEqual, sampleID)
			})

		})
	})
}

func TestIsStateField(t *testing.T) {
	Convey("testing isStateField", t, func() {
		goodExample := isStateField("/Instance/State")
		wrongExample1 := isStateField("/Instance/Statex")
		wrongExample2 := isStateField("/Instance/XState")

		So(goodExample, ShouldBeTrue)
		So(wrongExample1, ShouldBeFalse)
		So(wrongExample2, ShouldBeFalse)
	})
}

func TestGetIdFromKey(t *testing.T) {
	Convey("testing getIdFromKey", t, func() {
		id := "test-id"
		prefix := "/prefix/"
		suffix := "/suffix/"
		key := prefix + id + suffix

		goodExample := getIdFromKey(key, prefix, suffix)
		So(goodExample, ShouldEqual, id)

		goodExample = getIdFromKey(id, "", "")
		So(goodExample, ShouldEqual, id)

		key = prefix + keySeparator + id + keySeparator + suffix
		goodExample = getIdFromKey(key, prefix, suffix)
		So(goodExample, ShouldEqual, id)
	})
}

func TestIsInstanceTypeOf(t *testing.T) {
	testApplicationInstance1 := models.Instance{
		Id:   "1",
		Name: "app-1",
		Type: models.InstanceTypeApplication,
	}
	testServiceInstance1 := models.Instance{
		Id:   "2",
		Name: "service-1",
		Type: models.InstanceTypeService,
	}
	testServiceBrokerInstance1 := models.Instance{
		Id:   "3",
		Name: "service-broker-1",
		Type: models.InstanceTypeServiceBroker,
	}
	Convey("testing IsIsInstanceTypeOf", t, func() {
		Convey("For application instance and expected type application should return true", func() {
			isApplication := IsInstanceTypeOf(testApplicationInstance1, models.InstanceTypeApplication)
			So(isApplication, ShouldBeTrue)
		})
		Convey("For application instance and expected type service should return false", func() {
			isApplication := IsInstanceTypeOf(testApplicationInstance1, models.InstanceTypeService)
			So(isApplication, ShouldBeFalse)
		})
		Convey("For application instance and expected type service broker should return false", func() {
			isApplication := IsInstanceTypeOf(testApplicationInstance1, models.InstanceTypeServiceBroker)
			So(isApplication, ShouldBeFalse)
		})

		Convey("For service instance and expected type application should return false", func() {
			isApplication := IsInstanceTypeOf(testServiceInstance1, models.InstanceTypeApplication)
			So(isApplication, ShouldBeFalse)
		})
		Convey("For service instance and expected type service should return true", func() {
			isApplication := IsInstanceTypeOf(testServiceInstance1, models.InstanceTypeService)
			So(isApplication, ShouldBeTrue)
		})
		Convey("For service instance and expected type service broker should return false", func() {
			isApplication := IsInstanceTypeOf(testServiceInstance1, models.InstanceTypeServiceBroker)
			So(isApplication, ShouldBeFalse)
		})

		Convey("For service broker instance and expected type application should return false", func() {
			isApplication := IsInstanceTypeOf(testServiceBrokerInstance1, models.InstanceTypeApplication)
			So(isApplication, ShouldBeFalse)
		})
		Convey("For service broker instance and expected type service should return false", func() {
			isApplication := IsInstanceTypeOf(testServiceBrokerInstance1, models.InstanceTypeService)
			So(isApplication, ShouldBeFalse)
		})
		Convey("For service broker instance and expected type service broker should return true", func() {
			isApplication := IsInstanceTypeOf(testServiceBrokerInstance1, models.InstanceTypeServiceBroker)
			So(isApplication, ShouldBeTrue)
		})
	})
}

func TestIsRunnungInstance(t *testing.T) {
	testApplicationInstance1 := models.Instance{
		Id:    "1",
		Name:  "app-1",
		Type:  models.InstanceTypeApplication,
		State: models.InstanceStateRunning,
	}
	testApplicationInstance2 := models.Instance{
		Id:    "2",
		Name:  "app-2",
		Type:  models.InstanceTypeApplication,
		State: models.InstanceStateStopReq,
	}
	testApplicationInstance3 := models.Instance{
		Id:    "2",
		Name:  "app-2",
		Type:  models.InstanceTypeApplication,
		State: models.InstanceStateFailure,
	}
	testServiceInstance1 := models.Instance{
		Id:    "2",
		Name:  "service-1",
		Type:  models.InstanceTypeService,
		State: models.InstanceStateRunning,
	}
	testServiceBrokerInstance1 := models.Instance{
		Id:    "2",
		Name:  "service-1",
		Type:  models.InstanceTypeServiceBroker,
		State: models.InstanceStateRunning,
	}
	Convey("testing IsApplicationInstance", t, func() {
		Convey("For application instance in state RUNNING should return true", func() {
			isRunningInstance := IsRunningInstance(testApplicationInstance1)
			So(isRunningInstance, ShouldBeTrue)
		})
		Convey("For application instance in state STOP_REQ should return true", func() {
			isRunningInstance := IsRunningInstance(testApplicationInstance2)
			So(isRunningInstance, ShouldBeTrue)
		})
		Convey("For application instance in state FAILURE should return false", func() {
			isRunningInstance := IsRunningInstance(testApplicationInstance3)
			So(isRunningInstance, ShouldBeFalse)
		})
		Convey("For service instance in state RUNNING should return false", func() {
			isRunningInstance := IsRunningInstance(testServiceInstance1)
			So(isRunningInstance, ShouldBeTrue)
		})
		Convey("For service broker in state RUNNING should return false", func() {
			isRunningInstance := IsRunningInstance(testServiceBrokerInstance1)
			So(isRunningInstance, ShouldBeTrue)
		})
	})
}

func TestIsAuditTrailKey(t *testing.T) {
	testCases := []struct {
		path   string
		result bool
	}{
		{path: "", result: false},
		{path: "/", result: false},
		{path: "//", result: false},
		{path: "/asfd/dsg", result: false},
		{path: "/AuditTrail", result: false},
		{path: "/AuditTrail/LastUpdate", result: true},
		{path: "/org/AuditTrail/LastUpdate", result: true},
	}

	Convey("For set of test cases, IsAuditTrail should return proper response", t, func() {
		for _, tc := range testCases {
			Convey(fmt.Sprintf("For path %q result should be proper", tc.path), func() {
				result := isAuditTrailKey(tc.path)
				So(result, ShouldEqual, tc.result)
			})
		}
	})
}
