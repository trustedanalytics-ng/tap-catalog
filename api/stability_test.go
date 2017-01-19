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

func TestCheckStateStability(t *testing.T) {
	Convey("Testing CheckStateStability", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)
		sampleInstances := getInstancesInStableState()

		Convey("When all instances are in stable state", func() {
			sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

			mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)

			result, status, err := catalogClient.CheckStateStability()

			Convey("err should be proper", func() {
				So(err, ShouldBeNil)
			})
			Convey("status should be proper", func() {
				So(status, ShouldEqual, http.StatusOK)
			})
			Convey("result should be proper", func() {
				So(result.Stable, ShouldEqual, true)
			})
		})

		statesNotStable := []models.InstanceState{
			models.InstanceStateRequested,
			models.InstanceStateDeploying,
			models.InstanceStateStartReq,
			models.InstanceStateStarting,
			models.InstanceStateReconfiguration,
			models.InstanceStateStopReq,
			models.InstanceStateStopping,
			models.InstanceStateDestroyReq,
			models.InstanceStateDestroying,
			models.InstanceStateUnavailable,
		}

		for _, stateNotReady := range statesNotStable {
			Convey(fmt.Sprintf("When there is an instance with state %q", stateNotReady), func() {
				sampleInstances = append(sampleInstances, models.Instance{State: stateNotReady})
				sampleInstancesAsListOfInterfaces := getSampleInstancesAsListOfInterfaces(sampleInstances)

				mocks.repositoryMock.EXPECT().GetListOfData(context.getInstanceKey(), models.Instance{}).Return(sampleInstancesAsListOfInterfaces, nil)

				result, status, err := catalogClient.CheckStateStability()

				Convey("err should be proper", func() {
					So(err, ShouldBeNil)
				})
				Convey("status should be proper", func() {
					So(status, ShouldEqual, http.StatusOK)
				})
				Convey("result should be proper", func() {
					So(result.Stable, ShouldEqual, false)
				})
			})
		}

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func getInstancesInStableState() []models.Instance {
	return []models.Instance{
		{State: models.InstanceStateStopped},
		{State: models.InstanceStateRunning},
		{State: models.InstanceStateFailure},
	}
}
