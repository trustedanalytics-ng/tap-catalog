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

package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
)

const (
	instanceId string = "test-instance-id"
	serviceId  string = "test-service-id"
	planId     string = "test-plan-id"
)

func TestAddServiceInstance(t *testing.T) {
	Convey("Test Add Service Instance", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("Adding instance ok, response status is 201", func() {
			instance := getSampleInstance()
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				mocks.repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
				mocks.repositoryMock.EXPECT().IsExistByName(instance.Name, models.Instance{}, context.getInstanceKey()).Return(false, nil),
				mocks.repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(nil),
				mocks.repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil),
				mocks.repositoryMock.EXPECT().GetData(gomock.Any(), models.Instance{}).Return(instance, nil),
			)

			responseInstance, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldBeNil)
			So(status, ShouldEqual, http.StatusCreated)
			So(responseInstance, ShouldResemble, instance)
		})

		Convey("Adding service-broker instance ok if no PLAN_ID, response status is 201", func() {
			instance := getSampleInstance()
			instance.Type = models.InstanceTypeServiceBroker
			instance.Metadata = []models.Metadata{}
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
				mocks.repositoryMock.EXPECT().IsExistByName(instance.Name, models.Instance{}, context.getInstanceKey()).Return(false, nil),
				mocks.repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(nil),
				mocks.repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil),
				mocks.repositoryMock.EXPECT().GetData(gomock.Any(), models.Instance{}).Return(instance, nil),
			)

			responseInstance, status, err := catalogClient.AddServiceBrokerInstance(serviceId, instance)

			So(err, ShouldBeNil)
			So(status, ShouldEqual, http.StatusCreated)
			So(responseInstance, ShouldResemble, instance)
		})

		Convey("Service not exist, response status is 404", func() {
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, errors.New("not exist")),
			)

			_, status, err := catalogClient.AddServiceInstance(serviceId, models.Instance{})

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusNotFound)
			So(err.Error(), ShouldContainSubstring, "does not exists")
		})

		Convey("Id field not empty, response status is 400", func() {
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Id = instanceId

			_, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusBadRequest)
			So(err.Error(), ShouldContainSubstring, "Id field has to be empty!")
		})

		Convey("Plan not found, response status is 400", func() {
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Metadata = []models.Metadata{}

			_, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusBadRequest)
			So(err.Error(), ShouldContainSubstring, fmt.Sprintf("key %s not found!", models.OFFERING_PLAN_ID))
		})

		Convey("Instance name does not match lowercase DNS rule, response status is 400", func() {
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Name = "NOT_DNS"

			_, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusBadRequest)
			So(err.Error(), ShouldContainSubstring, "Field: Name has incorrect value: "+instance.Name)
		})

		Convey("Instance already exist, response status is 409", func() {
			instance := getSampleInstance()
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				mocks.repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
				mocks.repositoryMock.EXPECT().IsExistByName(instance.Name, models.Instance{}, context.getInstanceKey()).Return(true, nil),
			)

			_, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusConflict)
			So(err.Error(), ShouldContainSubstring, "already exists!")
		})

		Convey("Binding data does not match envs validation rule, response status is 400", func() {
			instance := getSampleInstance()
			instance.Bindings = []models.InstanceBindings{
				{Id: "bindingId", Data: map[string]string{"not-env-key": "value"}},
			}

			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				mocks.repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
			)

			_, status, err := catalogClient.AddServiceInstance(serviceId, instance)

			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusBadRequest)
			So(err.Error(), ShouldContainSubstring, "Field: data has incorrect value:")
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func getSampleInstance() models.Instance {
	return models.Instance{
		Name:  "test-name",
		Type:  models.InstanceTypeService,
		State: models.InstanceStateRunning,
		Metadata: []models.Metadata{
			{Id: models.OFFERING_PLAN_ID, Value: "plan_id"},
		},
		Bindings: []models.InstanceBindings{
			{Id: "bindingId", Data: map[string]string{"key": "value"}},
		},
	}
}
