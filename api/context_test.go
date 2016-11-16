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
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInitDB(t *testing.T) {
	_, context, repositoryMock := prepareMocksAndRouter(t)
	const org = "SAMPLE_ORG"

	Convey(fmt.Sprintf("Given some Context instance and %s organization", org), t, func() {
		Convey(fmt.Sprintf("Context.initBD should call repository.CreateDirs with %s parameter", org), func() {
			repositoryMock.EXPECT().CreateDirs(org).Return(nil)

			err := context.initDB(org)

			Convey("initDB response should be proper", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestReserveID(t *testing.T) {
	_, context, repositoryMock := prepareMocksAndRouter(t)
	const samplePath = "application"

	Convey("Testing reserveID", t, func() {
		Convey("when first UUID generation trial is successful", func() {
			repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(nil)

			retUUID, err := context.reserveID(samplePath)

			Convey("ReserveID response is proper", func() {
				Convey("Err is proper", func() {
					So(err, ShouldBeNil)
				})
				Convey("Uuid is not empty", func() {
					So(len(retUUID), ShouldBeGreaterThan, 0)
				})
			})
		})

		Convey("When second UUID generation trial is successful", func() {
			gomock.InOrder(
				repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(errors.New("")),
				repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(nil),
			)

			retUUID, err := context.reserveID(samplePath)

			Convey("ReserveID response is proper", func() {
				Convey("Err is proper", func() {
					So(err, ShouldBeNil)
				})
				Convey("Uuid is not empty", func() {
					So(len(retUUID), ShouldBeGreaterThan, 0)
				})
			})
		})

		Convey("When maxUUIDGenerationTrials=%d generation trials are not successful", func() {
			repositoryMock.EXPECT().CreateDir(gomock.Any()).Return(errors.New("")).AnyTimes()

			_, err := context.reserveID(samplePath)

			Convey("Response err is not nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
