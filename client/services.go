package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *TapCatalogApiConnector) GetServices() ([]models.Service, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, services))
	result := &[]models.Service{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) GetService(serviceId string) (models.Service, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, serviceId))
	result := &models.Service{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdateService(serviceId string, patches []models.Patch) (models.Service, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, serviceId))
	result := &models.Service{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/plans/%s", c.Address, services, serviceId, planId))
	result := &models.ServicePlan{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) AddService(service models.Service) (models.Service, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, services))
	result := &models.Service{}
	status, err := brokerHttp.PostModel(connector, service, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) DeleteService(serviceId string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, serviceId))
	status, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return status, err
}
