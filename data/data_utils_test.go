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
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
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
