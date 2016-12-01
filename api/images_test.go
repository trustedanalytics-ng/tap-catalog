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

func TestAddImage(t *testing.T) {
	router, _, repositoryMock := prepareMocksAndRouter(t)

	Convey("Testing AddImage", t, func() {
		Convey("When providing AddImage with proper Image", func() {
			image := getSampleImage()
			gomock.InOrder(
				repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil),
				repositoryMock.EXPECT().GetData(gomock.Any(), models.Image{}).Return(image, nil),
			)

			catalogClient := getCatalogClient(router, t)
			responseImage, status, err := catalogClient.AddImage(image)

			Convey("response should be proper", func() {
				So(err, ShouldBeNil)

				Convey("status code should be proper", func() {
					So(status, ShouldEqual, http.StatusCreated)
				})

				Convey("returned Image should be proper", func() {
					So(responseImage, ShouldResemble, image)
				})
			})
		})
	})
}

func TestMonitorImagesState(t *testing.T) {
	router, context, repositoryMock := prepareMocksAndRouter(t)

	stateChange := models.StateChange{
		Id: "test",
	}

	Convey("Testing MonitorImagesState", t, func() {
		Convey("Request correct, response status is 200", func() {
			afterIndex := models.WatchFromNow
			gomock.InOrder(
				repositoryMock.EXPECT().MonitorObjectsStates(context.buildImagesKey(""), afterIndex).Return(stateChange, nil),
			)

			catalogClient := getCatalogClient(router, t)
			response, status, err := catalogClient.WatchImages(afterIndex)

			So(status, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(response, ShouldResemble, stateChange)
		})
	})

	Convey("Testing MonitorSpecificImageState", t, func() {
		Convey("Request correct, response status is 200", func() {
			afterIndex := models.WatchFromNow
			imageId := "test-image"
			gomock.InOrder(
				repositoryMock.EXPECT().MonitorObjectsStates(context.buildImagesKey(imageId), afterIndex).Return(stateChange, nil),
			)

			catalogClient := getCatalogClient(router, t)
			response, status, err := catalogClient.WatchImage(imageId, afterIndex)

			So(status, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(response, ShouldResemble, stateChange)
		})
	})
}

func getSampleImage() models.Image {
	return models.Image{
		Type:     models.ImageTypeJava,
		BlobType: models.BlobTypeJar,
		State:    models.ImageStateBuilding,
	}
}
