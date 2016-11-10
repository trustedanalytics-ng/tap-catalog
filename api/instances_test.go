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
	"strings"
	"testing"

	"github.com/gocraft/web"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

const (
	instanceId string = "test-instance-id"
	serviceId  string = "test-service-id"

	urlPrefix              = "/api/v1"
	serviceIDWildcard      = ":serviceId"
	urlPostServiceInstance = urlPrefix + "/services/" + serviceIDWildcard + "/instances"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router, c Context, repositoryMock *data.MockRepositoryApi) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repositoryMock = data.NewMockRepositoryApi(mockCtrl)
	c = Context{
		Repository: repositoryMock,
	}
	router = web.New(c)
	return
}

func TestAddServiceInstance(t *testing.T) {
	router, context, repositoryMock := prepareMocksAndRouter(t)
	router.Post(urlPostServiceInstance, context.AddServiceInstance)

	Convey("Test Add Service Instance", t, func() {
		Convey("Adding instance ok, response status is 201", func() {
			instance := getSampleInstance()
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
				repositoryMock.EXPECT().IsExistByName(instance.Name, models.Instance{}, context.getInstanceKey()).Return(false, nil),
				repositoryMock.EXPECT().StoreData(gomock.Any()).Return(nil),
				repositoryMock.EXPECT().GetData(gomock.Any(), models.Instance{}).Return(instance, nil),
			)

			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "", http.StatusCreated)

			responseInstance := models.Instance{}
			err := util.ReadJsonFromByte(response.Body.Bytes(), &responseInstance)
			So(err, ShouldBeNil)
			So(responseInstance, ShouldResemble, instance)
		})

		Convey("Service not exist, response status is 404", func() {
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, errors.New("not exist")),
			)

			byteBody := util.PrepareAndValidateRequest(getSampleInstance(), t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "does not exists", http.StatusNotFound)
		})

		Convey("Id field not empty, response status is 400", func() {
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Id = instanceId
			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "Id field has to be empty!", http.StatusBadRequest)
		})

		Convey("Plan not found, response status is 400", func() {
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Metadata = []models.Metadata{}
			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, fmt.Sprintf("key %s not found!", models.OFFERING_PLAN_ID), http.StatusBadRequest)
		})

		Convey("Instance name does not match lowercase DNS rule, response status is 400", func() {
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
			)

			instance := getSampleInstance()
			instance.Name = "NOT_DNS"
			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "Field: Name has incorrect value: "+instance.Name, http.StatusBadRequest)
		})

		Convey("Instance already exist, response status is 409", func() {
			instance := getSampleInstance()
			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
				repositoryMock.EXPECT().IsExistByName(instance.Name, models.Instance{}, context.getInstanceKey()).Return(true, nil),
			)

			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "already exists!", http.StatusConflict)
		})

		Convey("Binding data does not match envs validation rule, response status is 400", func() {
			instance := getSampleInstance()
			instance.Bindings = []models.InstanceBindings{
				{Id: "bindingId", Data: map[string]string{"not-env-key": "value"}},
			}

			gomock.InOrder(
				repositoryMock.EXPECT().GetData(context.buildServiceKey(serviceId), models.Service{}).Return(models.Service{}, nil),
				repositoryMock.EXPECT().GetData(context.buildInstanceKey(instance.Bindings[0].Id), models.Instance{}).Return(models.Instance{}, nil),
			)

			byteBody := util.PrepareAndValidateRequest(instance, t)
			requestPath := strings.Replace(urlPostServiceInstance, serviceIDWildcard, serviceId, 1)
			response := util.SendRequest("POST", requestPath, byteBody, router)
			util.AssertResponse(response, "Field: data has incorrect value:", http.StatusBadRequest)
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
