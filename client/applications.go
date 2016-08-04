package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) AddApplication(application models.Application) (models.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, applications))
	result := &models.Application{}
	err := brokerHttp.PostModel(connector, application, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) ListApplications() ([]models.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, applications))
	result := &[]models.Application{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) GetApplication(applicationId string) (models.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, applicationId))
	result := &models.Application{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdateApplication(applicationId string, patches []models.Patch) (models.Application, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, applicationId))
	result := &models.Application{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}
