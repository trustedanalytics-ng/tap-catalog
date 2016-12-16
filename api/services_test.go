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
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/trustedanalytics/tap-catalog/models"
)

const (
	sampleID1 = "1"
	sampleID2 = "2"

	sampleName1 = "sample_name_1"
	sampleName2 = "sample_name_2"
)

func TestGetServices(t *testing.T) {
	Convey("Testing GetServices", t, func() {
		mockCtrl, context, repositoryMock, catalogClient := prepareMocksAndClient(t)
		sampleServices := getSampleServices()
		sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)

		Convey("When RepositoryAPI returns proper data", func() {
			repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			services, status, err := catalogClient.GetServices()

			Convey("err should be proper", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
			Convey("returned services should be proper", func() {
				So(services, ShouldResemble, sampleServices)
			})
		})

		Convey("When RepositoryAPI returns improper data", func() {
			sampleServicesAsListOfInterfaces = append(sampleServicesAsListOfInterfaces, nil)

			repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			_, status, err := catalogClient.GetServices()

			Convey("err should be proper", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
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
		mockCtrl, context, repositoryMock, catalogClient := prepareMocksAndClient(t)
		sampleService := getSampleServices()[0]
		id := sampleID1

		Convey("When RepositoryAPI returns proper data", func() {
			sampleServiceInterface := interface{}(sampleService)

			repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(sampleServiceInterface, nil)

			service, status, err := catalogClient.GetService(id)

			Convey("err should be proper", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
			Convey("returned services should be proper", func() {
				So(service, ShouldResemble, sampleService)
			})
		})

		Convey("When service with given ID does not exist", func() {
			id = "not-existing-id"

			repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(nil, fmt.Errorf("not found"))

			_, status, err := catalogClient.GetService(id)

			Convey("err should be proper", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When RepositoryAPI returns improper data", func() {
			sampleServiceInterface := interface{}(2)

			repositoryMock.EXPECT().GetData(context.buildServiceKey(id), models.Service{}).Return(sampleServiceInterface, nil)

			_, status, err := catalogClient.GetService(id)

			Convey("err should be proper", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
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
		mockCtrl, context, repositoryMock, catalogClient := prepareMocksAndClient(t)

		Convey("When offering can be deleted", func() {
			sampleServices := getSampleServices()
			sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)
			sampleInstances := getSampleInstances()
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)
			repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)
			repositoryMock.EXPECT().DeleteData(context.buildServiceKey(id)).Return(nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should be proper", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("When there exist an instance of offering", func() {
			sampleInstances := getSampleInstances()
			sampleInstances[0].ClassId = sampleID1
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should be proper", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusForbidden)
			})
		})

		Convey("When there is a dependency to offering being deleted", func() {
			sampleServices := getSampleServicesWithDependency()
			sampleServicesAsListOfInterfaces := getSampleServicesAsListOfInterfaces(sampleServices)
			sampleInstances := getSampleInstances()
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)
			repositoryMock.EXPECT().GetListOfData(context.getServiceKey(), models.Service{}).Return(sampleServicesAsListOfInterfaces, nil)

			status, err := catalogClient.DeleteService(id)

			Convey("err should be proper", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("status should be proper", func() {
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
		{Id: sampleID1, Name: sampleName1},
		{Id: sampleID2, Name: sampleName2},
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
		{Id: sampleID1, Name: sampleName1},
		{Id: sampleID2, Name: sampleName2},
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
