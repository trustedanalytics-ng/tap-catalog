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
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics-ng/tap-catalog/models"
)

const (
	sampleID1 = "1"
	sampleID2 = "2"

	sampleName1 = "sample-name-1"
	sampleName2 = "sample-name-2"
)

func TestAddService(t *testing.T) {
	Convey("Testing AddService", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("When Service is proper", func() {
			sampleService := getSampleServices()[0]
			sampleService.Id = ""
			sampleServiceInterface := interface{}(sampleService)

			mocks.repositoryMock.EXPECT().IsExistByName(sampleService.Name, models.Service{}, context.getServiceKey()).Return(false, nil)
			mocks.repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(nil)
			mocks.repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil)
			mocks.repositoryMock.EXPECT().GetData(gomock.Any(), models.Service{}).Return(sampleServiceInterface, nil)

			service, status, err := catalogClient.AddService(sampleService)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be Created", func() {
				So(status, ShouldEqual, http.StatusCreated)
			})
			Convey("returned services should be proper", func() {
				So(service, ShouldResemble, sampleService)
			})
		})

		Convey("When Service has no plan", func() {
			sampleService := models.Service{Id: "", Name: sampleName1}

			_, status, err := catalogClient.AddService(sampleService)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be BadRequest", func() {
				So(status, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When Service is provided with ID", func() {
			sampleService := getSampleServices()[0]

			_, status, err := catalogClient.AddService(sampleService)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be BadRequest", func() {
				So(status, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When Service is provided with wrong name", func() {
			sampleService := getSampleServices()[0]
			sampleService.Name = "name_with_underscore"

			_, status, err := catalogClient.AddService(sampleService)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be BadRequest", func() {
				So(status, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When Service is provided with name already used", func() {
			sampleService := getSampleServices()[0]
			sampleService.Id = ""

			mocks.repositoryMock.EXPECT().IsExistByName(sampleService.Name, models.Service{}, context.getServiceKey()).Return(true, nil)

			_, status, err := catalogClient.AddService(sampleService)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be Conflict", func() {
				So(status, ShouldEqual, http.StatusConflict)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestPatchServiceUpdate(t *testing.T) {
	Convey("Testing PatchService with Update operation", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		sampleService := models.Service{Id: sampleID1, Name: sampleName1, Plans: []models.ServicePlan{models.ServicePlan{}}, State: models.ServiceStateReady}
		sampleServiceInterface := interface{}(sampleService)
		fieldName := "state"
		newValueByte, _ := json.Marshal(models.ServiceStateOffline)
		patches := []models.Patch{models.Patch{Operation: models.OperationUpdate, Field: &fieldName, Value: (*json.RawMessage)(&newValueByte)}}
		patchedValues, _ := context.mapper.ToKeyValueByPatches(context.buildServiceKey(sampleService.Id), models.Service{}, patches)

		Convey("When field state is updated from READY to OFFLINE state", func() {
			sampleInstances := getSampleInstances()
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)
			sampleServices := getSampleServices()
			sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)

			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(sampleService.Id), models.Service{}).Return(sampleServiceInterface, nil)
			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)
			mocks.repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)
			mocks.repositoryMock.EXPECT().ApplyPatchedValues(patchedValues)
			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(sampleService.Id), models.Service{}).Return(sampleServiceInterface, nil)

			service, status, err := catalogClient.UpdateService(sampleService.Id, patches)

			Convey("response should be proper", func() {
				So(err, ShouldBeNil)
				So(status, ShouldEqual, http.StatusOK)
				So(service, ShouldResemble, sampleService)
			})
		})

		Convey("When field state is updated from READY to OFFLINE state and there is service instance", func() {
			sampleInstances := []models.Instance{models.Instance{ClassId: sampleService.Id}}
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(sampleService.Id), models.Service{}).Return(sampleServiceInterface, nil)
			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)

			_, status, err := catalogClient.UpdateService(sampleService.Id, patches)

			Convey("response should be proper", func() {
				So(err, ShouldNotBeNil)
				So(status, ShouldEqual, http.StatusForbidden)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestGetServices(t *testing.T) {
	Convey("Testing GetServices", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)
		sampleServices := getSampleServices()
		sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)

		Convey("When RepositoryAPI returns proper data", func() {
			mocks.repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			services, status, err := catalogClient.GetServices()

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
			Convey("returned services should be proper", func() {
				So(services, ShouldResemble, sampleServices)
			})
		})

		Convey("When RepositoryAPI returns improper data", func() {
			sampleServicesAsListOfInterfaces = append(sampleServicesAsListOfInterfaces, nil)

			mocks.repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			_, status, err := catalogClient.GetServices()

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be InternalServerError", func() {
				So(status, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestGetService(t *testing.T) {
	Convey("Testing GetService", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)
		sampleService := getSampleServices()[0]
		id := sampleID1

		Convey("When RepositoryAPI returns proper data", func() {
			sampleServiceInterface := interface{}(sampleService)

			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(sampleServiceInterface, nil)

			service, status, err := catalogClient.GetService(id)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be OK", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
			Convey("returned services should be proper", func() {
				So(service, ShouldResemble, sampleService)
			})
		})

		Convey("When service with given ID does not exist", func() {
			id = "not-existing-id"

			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(nil, errors.New("not found"))

			_, status, err := catalogClient.GetService(id)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When RepositoryAPI returns improper data", func() {
			sampleServiceInterface := interface{}(2)

			mocks.repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(sampleServiceInterface, nil)

			_, status, err := catalogClient.GetService(id)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be InternalServerError", func() {
				So(status, ShouldEqual, http.StatusInternalServerError)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestDeleteService(t *testing.T) {
	id := sampleID1

	Convey("Testing DeleteService", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("When offering can be deleted", func() {
			sampleServices := getSampleServices()
			sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)
			sampleInstances := getSampleInstances()
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)
			mocks.repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)
			mocks.repositoryMock.EXPECT().DeleteData(context.buildServiceKey(id)).Return(nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be NoContent", func() {
				So(status, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("When there exist an instance of offering", func() {
			sampleInstances := getSampleInstances()
			sampleInstances[0].ClassId = sampleID1
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be Forbidden", func() {
				So(status, ShouldEqual, http.StatusForbidden)
			})
		})

		Convey("When there is a dependency to offering being deleted", func() {
			sampleServices := getSampleServicesWithDependency()
			sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)
			sampleInstances := getSampleInstances()
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)
			mocks.repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be Forbidden", func() {
				So(status, ShouldEqual, http.StatusForbidden)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func getSampleServices() []models.Service {
	return []models.Service{
		{Id: sampleID1, Name: sampleName1, Plans: []models.ServicePlan{models.ServicePlan{}}},
		{Id: sampleID2, Name: sampleName2, Plans: []models.ServicePlan{models.ServicePlan{}}},
	}
}

func getSampleServicesWithDependency() []models.Service {
	dependencies := []models.ServiceDependency{{ServiceId: sampleID1}}
	plans := []models.ServicePlan{{Dependencies: dependencies}}
	return []models.Service{
		{Id: sampleID1, Name: sampleName1},
		{Id: sampleID2, Name: sampleName2, Plans: plans},
	}
}

func getSampleInstances() []models.Instance {
	return []models.Instance{
		{Id: sampleID1, Name: sampleName1, Type: models.InstanceTypeService, ClassId: serviceId},
		{Id: sampleID2, Name: sampleName2, Type: models.InstanceTypeService, ClassId: serviceId},
	}
}

func getSampleServicesAsListOfInterfaces(services []models.Service) []interface{} {
	result := make([]interface{}, len(services))

	for i, service := range services {
		result[i] = interface{}(service)
	}
	return result
}

func getSampleInstancesAsListOfInterfaces(instances []models.Instance) []interface{} {
	result := make([]interface{}, len(instances))

	for i, instance := range instances {
		result[i] = interface{}(instance)
	}
	return result
}
