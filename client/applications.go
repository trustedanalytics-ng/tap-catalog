/**
 * Copyright (c) 2017 Intel Corporation
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
package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics-ng/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func (c *TapCatalogApiConnector) AddApplication(application models.Application) (models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, applications))
	result := &models.Application{}
	status, err := brokerHttp.PostModel(connector, application, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) ListApplications(filter *brokerHttp.ItemFilter) ([]models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s?%s", c.Address, applications, filter.BuildQuery()))
	result := &[]models.Application{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) GetApplication(applicationId string) (models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, applicationId))
	result := &models.Application{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdateApplication(applicationId string, patches []models.Patch) (models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, applicationId))
	result := &models.Application{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) DeleteApplication(applicationId string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, applicationId))
	status, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return status, err
}
