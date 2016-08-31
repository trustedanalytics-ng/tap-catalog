package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *TapCatalogApiConnector) AddApplication(application models.Application) (models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, applications))
	result := &models.Application{}
	status, err := brokerHttp.PostModel(connector, application, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) ListApplications() ([]models.Application, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, applications))
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
