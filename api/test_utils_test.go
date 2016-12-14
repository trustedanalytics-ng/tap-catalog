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
	"net/http/httptest"
	"testing"

	"github.com/gocraft/web"
	"github.com/golang/mock/gomock"

	"github.com/trustedanalytics/tap-catalog/client"
	"github.com/trustedanalytics/tap-catalog/data"
)

const (
	urlPrefix                  = "/api/v1"
	urlPostImage               = urlPrefix + "/images"
	urlGetImageWatcher         = urlPostImage + "/nextState"
	imageIDWildcard            = ":imageId"
	urlGetSpecificImageWatcher = urlPostImage + "/" + imageIDWildcard + "/nextState"
	serviceIDWildcard          = ":serviceId"
	urlPostServiceInstance     = urlPrefix + "/services/" + serviceIDWildcard + "/instances"
	urlGetLatestIndex          = urlPrefix + "/latestIndex"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router, c Context, repositoryMock *data.MockRepositoryApi) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repositoryMock = data.NewMockRepositoryApi(mockCtrl)
	c = Context{
		repository: repositoryMock,
	}
	router = web.New(c)

	router.Post(urlPostImage, c.AddImage)
	router.Post(urlPostServiceInstance, c.AddServiceInstance)
	router.Get(urlGetImageWatcher, c.MonitorImagesStates)
	router.Get(urlGetSpecificImageWatcher, c.MonitorSpecificImageState)
	router.Get(urlGetLatestIndex, c.LatestIndex)

	return
}

func getCatalogClient(router *web.Router, t *testing.T) *client.TapCatalogApiConnector {
	testServer := httptest.NewServer(router)
	cBroker, err := client.NewTapCatalogApiWithBasicAuth(testServer.URL, "user", "password")
	if err != nil {
		t.Fatal("Container broker error: ", err)
	}
	return cBroker
}
