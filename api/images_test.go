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
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

const (
	urlPrefix       = "/api/v1"
	imageIDWildcard = ":imageId"
	urlPostImage    = urlPrefix + "/images"
)

func TestAddImage(t *testing.T) {
	router, context, repositoryMock := prepareMocksAndRouter(t)
	router.Post(urlPostImage, context.AddImage)

	Convey("Testing AddImage", t, func() {
		Convey("When providing AddImage with proper Image", func() {
			image := getSampleImage()
			gomock.InOrder(
				repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil),
				repositoryMock.EXPECT().GetData(gomock.Any(), models.Image{}).Return(image, nil),
			)

			byteBody, err := json.Marshal(image)
			if err != nil {
				t.Fatalf("cannot marshal %v", image)
			}

			response := util.SendRequest("POST", urlPostImage, byteBody, router)

			Convey("response should be proper", func() {
				util.AssertResponse(response, "", http.StatusCreated)

				Convey("status code should be proper", func() {
					So(response.Code, ShouldEqual, http.StatusCreated)
				})

				responseImage := models.Image{}
				err := util.ReadJsonFromByte(response.Body.Bytes(), &responseImage)
				Convey("unmarshal error is nil", func() {
					So(err, ShouldBeNil)

					Convey("returned Image should be proper", func() {
						So(responseImage, ShouldResemble, image)
					})
				})
			})
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
