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

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
)

func TestLatestIndex(t *testing.T) {
	Convey("Testing LatestIndex", t, func() {
		mockCtrl, _, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("When making request", func() {
			latestIndex := uint64(5)
			index := models.Index{Latest: latestIndex}
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().GetLatestIndex(gomock.Any()).Return(latestIndex, nil),
			)

			responseIndex, status, err := catalogClient.GetLatestIndex()

			Convey("response should be proper", func() {
				So(err, ShouldBeNil)

				Convey("status code should be proper", func() {
					So(status, ShouldEqual, http.StatusOK)
				})

				Convey("returned Image should be proper", func() {
					So(responseIndex, ShouldResemble, index)
				})
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}
