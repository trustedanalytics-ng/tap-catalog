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
	"github.com/nu7hatch/gouuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
)

func TestAddImage(t *testing.T) {
	Convey("Testing AddImage", t, func() {
		mockCtrl, _, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("When providing AddImage with proper Image", func() {
			image := getSampleImage()
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().CreateData(gomock.Any()).Return(nil),
				mocks.repositoryMock.EXPECT().GetData(gomock.Any(), models.Image{}).Return(image, nil),
			)

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

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestMonitorImagesState(t *testing.T) {
	stateChange := models.StateChange{
		Id: "test",
	}

	Convey("Testing MonitorImagesState", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("Request correct, response status is 200", func() {
			afterIndex := models.WatchFromNow
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().MonitorObjectsStates(context.buildImagesKey(""), afterIndex).Return(stateChange, nil),
			)

			response, status, err := catalogClient.WatchImages(afterIndex)

			So(status, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(response, ShouldResemble, stateChange)
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})

	Convey("Testing MonitorSpecificImageState", t, func() {
		mockCtrl, context, mocks, catalogClient := prepareMocksAndClient(t)

		Convey("Request correct, response status is 200", func() {
			afterIndex := models.WatchFromNow
			imageId := "test-image"
			gomock.InOrder(
				mocks.repositoryMock.EXPECT().MonitorObjectsStates(context.buildImagesKey(imageId), afterIndex).Return(stateChange, nil),
			)

			response, status, err := catalogClient.WatchImage(imageId, afterIndex)

			So(status, ShouldEqual, http.StatusOK)
			So(err, ShouldBeNil)
			So(response, ShouldResemble, stateChange)
		})

		Reset(func() {
			mockCtrl.Finish()
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

type AppArray []models.Application

func (a AppArray) ConvertToInterfaces() []interface{} {
	var result []interface{}

	for _, app := range a {
		result = append(result, app)
	}
	return result
}

type SrvArray []models.Service

func (a SrvArray) ConvertToInterfaces() []interface{} {
	var result []interface{}

	for _, el := range a {
		result = append(result, el)
	}
	return result
}

const fakeDockerRegistryAddress = "10.10.5.10:3211"

type GetImageCheckRefsTestCaseDefinition struct {
	TestType         GetImageCheckRefsTestCaseType
	ImageIDToTest    string
	Image            models.Image
	Applications     AppArray
	Offerings        SrvArray
	ExpectedResponse models.ImageRefsResponse
	IsError          bool
}

func GenerateID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id.String()
}

type GetImageCheckRefsTestCaseType string

const (
	OnlyApplication         GetImageCheckRefsTestCaseType = "only application"
	OnlyService             GetImageCheckRefsTestCaseType = "only service"
	Both                    GetImageCheckRefsTestCaseType = "application and service"
	None                    GetImageCheckRefsTestCaseType = "none"
	BadImageAddress         GetImageCheckRefsTestCaseType = "bad image address"
	GetListOfDataAppFailure GetImageCheckRefsTestCaseType = "GetListOfData failure for application"
	GetListOfDataSrvFailure GetImageCheckRefsTestCaseType = "GetListOfData failure for service"
	GetListOfDataAppBadData GetImageCheckRefsTestCaseType = "GetListOfData returned bad data for application"
	GetListOfDataSrvBadData GetImageCheckRefsTestCaseType = "GetListOfData returned bad data for service"
)

func makeImageAddress(ImageID string) string {
	return fmt.Sprintf("%s/%s%s", fakeDockerRegistryAddress, models.USER_DEFINED_APPLICATION_IMAGE_PREFIX, ImageID)
}

func NewGetImageCheckRefsTestCase(testType GetImageCheckRefsTestCaseType) GetImageCheckRefsTestCaseDefinition {
	fakeImageID := GenerateID()
	expectedAppID := GenerateID()
	expectedSrvID := GenerateID()
	imageAddress := makeImageAddress(fakeImageID)

	testCase := GetImageCheckRefsTestCaseDefinition{
		TestType:      testType,
		ImageIDToTest: fakeImageID,
		Image: models.Image{
			Id: fakeImageID,
		},

		Applications: []models.Application{
			{
				Id:      GenerateID(),
				ImageId: GenerateID(),
			},
			{
				Id:      GenerateID(),
				ImageId: GenerateID(),
			},
		},
		Offerings: []models.Service{
			{
				Id: GenerateID(),
				Metadata: []models.Metadata{
					{
						Id:    models.APPLICATION_IMAGE_ADDRESS,
						Value: makeImageAddress(GenerateID()),
					},
				},
			},
			{
				Id: GenerateID(),
			},
		},
	}

	if testType == OnlyApplication || testType == Both {
		appToBeFound := models.Application{
			Id:      expectedAppID,
			ImageId: fakeImageID,
		}
		testCase.Applications = append(testCase.Applications, appToBeFound)

		testCase.ExpectedResponse.ApplicationReferences = []models.Application{appToBeFound}
		testCase.ExpectedResponse.IsAnyRefExist = true
	} else if testType == OnlyService || testType == Both {
		srvToBeFound := models.Service{
			Id: expectedSrvID,
			Metadata: []models.Metadata{

				{
					Id:    models.APPLICATION_IMAGE_ADDRESS,
					Value: imageAddress,
				},
			},
		}
		testCase.Offerings = append(testCase.Offerings, srvToBeFound)

		testCase.ExpectedResponse.ServiceReferences = []models.Service{srvToBeFound}
		testCase.ExpectedResponse.IsAnyRefExist = true
	} else if testType == BadImageAddress {
		badSrv := models.Service{
			Id: expectedSrvID,
			Metadata: []models.Metadata{
				{
					Id:    models.APPLICATION_IMAGE_ADDRESS,
					Value: "some trash",
				},
			},
		}
		testCase.Offerings = append(testCase.Offerings, badSrv)
		testCase.IsError = true
	} else if testType == GetListOfDataAppFailure ||
		testType == GetListOfDataSrvFailure ||
		testType == GetListOfDataAppBadData ||
		testType == GetListOfDataSrvBadData {

		testCase.IsError = true
	}

	return testCase
}

func TestGetImageCheckRefs(t *testing.T) {
	testCases := []GetImageCheckRefsTestCaseDefinition{
		NewGetImageCheckRefsTestCase(OnlyApplication),
		NewGetImageCheckRefsTestCase(OnlyService),
		NewGetImageCheckRefsTestCase(Both),
		NewGetImageCheckRefsTestCase(None),
		NewGetImageCheckRefsTestCase(BadImageAddress),
		NewGetImageCheckRefsTestCase(GetListOfDataAppFailure),
		NewGetImageCheckRefsTestCase(GetListOfDataSrvFailure),
		NewGetImageCheckRefsTestCase(GetListOfDataAppBadData),
		NewGetImageCheckRefsTestCase(GetListOfDataSrvBadData),
	}

	Convey("Testing GetImageCheckRefs", t, func() {
		for _, tc := range testCases {
			mockCtrl, c, mocks, catalogClient := prepareMocksAndClient(t)

			Convey(fmt.Sprintf("GetImageRefs with \"%s\" test type", tc.TestType), func() {

				switch tc.TestType {
				case GetListOfDataAppFailure:
					mocks.repositoryMock.EXPECT().GetListOfData(c.getApplicationKey(), models.Application{}).Return([]interface{}{}, errors.New("something bad happened"))
				case GetListOfDataSrvFailure:
					gomock.InOrder(
						mocks.repositoryMock.EXPECT().GetListOfData(c.getApplicationKey(), models.Application{}).Return([]interface{}{}, nil),
						mocks.repositoryMock.EXPECT().GetListOfData(c.getServiceKey(), models.Service{}).Return([]interface{}{}, errors.New("something bad happened")),
					)
				case GetListOfDataAppBadData:
					mocks.repositoryMock.EXPECT().GetListOfData(c.getApplicationKey(), models.Application{}).Return([]interface{}{models.Service{}}, nil)
				case GetListOfDataSrvBadData:
					gomock.InOrder(
						mocks.repositoryMock.EXPECT().GetListOfData(c.getApplicationKey(), models.Application{}).Return([]interface{}{}, nil),
						mocks.repositoryMock.EXPECT().GetListOfData(c.getServiceKey(), models.Service{}).Return([]interface{}{models.Application{}}, nil),
					)
				default:
					gomock.InOrder(
						mocks.repositoryMock.EXPECT().GetListOfData(c.getApplicationKey(), models.Application{}).Return(tc.Applications.ConvertToInterfaces(), nil),
						mocks.repositoryMock.EXPECT().GetListOfData(c.getServiceKey(), models.Service{}).Return(tc.Offerings.ConvertToInterfaces(), nil),
					)
				}

				response, status, err := catalogClient.GetImageRefs(tc.ImageIDToTest)

				if tc.IsError {
					Convey("error should be returned, status should be Internal Server Error", func() {
						So(err, ShouldNotBeNil)
						So(status, ShouldEqual, http.StatusInternalServerError)
					})
				} else {
					Convey("error should be nil, status should be OK, response should resemble expected one", func() {
						So(err, ShouldBeNil)
						So(status, ShouldEqual, http.StatusOK)
						So(response, ShouldResemble, tc.ExpectedResponse)
					})
				}
			})
			mockCtrl.Finish()
		}
	})
}
