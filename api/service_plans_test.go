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
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/golang/mock/gomock"
	"github.com/trustedanalytics-ng/tap-catalog/models"
)

func TestDeletePlan(t *testing.T) {
	Convey("Testing DeleteService", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)
		samplePlan := getSampleServicePlan()
		sampleInstances := getSampleInstances()
		sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

		Convey("Should delete plan if is not in use by any instance", func() {
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.getServicedPlanIDKey(serviceId, planId), models.ServicePlan{}).Return(samplePlan, nil),
				mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil),
				mocks.repositoryMock.EXPECT().DeleteData(context.getServicedPlanIDKey(serviceId, planId)).Return(nil),
			)

			status, err := catalogClient.DeleteServicePlan(serviceId, planId)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, http.StatusNoContent)
		})

		Convey("Should return error on delete plan which is in use by instance", func() {
			planIdMetadata := models.Metadata{
				Id:    models.OFFERING_PLAN_ID,
				Value: planId,
			}
			sampleInstances[0].Metadata = append(sampleInstances[0].Metadata, planIdMetadata)
			sampleInstancesAsListOfInterfaces = getSampleInstancesAsListOfInterfaces(sampleInstances)

			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetData(context.getServicedPlanIDKey(serviceId, planId), models.ServicePlan{}).Return(samplePlan, nil),
				mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil),
			)

			status, err := catalogClient.DeleteServicePlan(serviceId, planId)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, http.StatusBadRequest)
			So(err.Error(), ShouldContainSubstring, "can not be deleted - is in use by instance")
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func getSampleServicePlan() models.ServicePlan {
	return models.ServicePlan{
		Id: sampleID1, Name: sampleName1,
	}
}
