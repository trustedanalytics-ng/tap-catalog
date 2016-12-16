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
	"os"
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

func prepareMocksAndClient(t *testing.T) (mockCtrl *gomock.Controller, c Context, repositoryMock *data.MockRepositoryApi, catalogClient client.TapCatalogApi) {
	mockCtrl = gomock.NewController(t)
	repositoryMock = data.NewMockRepositoryApi(mockCtrl)
	c = Context{
		repository: repositoryMock,
	}

	router := SetupRouter(c)

	catalogClient = getCatalogClient(router, t)

	return
}

func getCatalogClient(router *web.Router, t *testing.T) client.TapCatalogApi {
	const user = "user"
	const password = "password"

	os.Setenv("CATALOG_USER", user)
	os.Setenv("CATALOG_PASS", password)

	testServer := httptest.NewServer(router)
	catalogClient, err := client.NewTapCatalogApiWithBasicAuth(testServer.URL, user, password)
	if err != nil {
		t.Fatal("Catalog client error: ", err)
	}
	return catalogClient
}
