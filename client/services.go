package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) GetServices() ([]models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, services))
	result := &[]models.Service{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) GetService(serviceId string) (models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, serviceId))
	result := &models.Service{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdateService(serviceId string, patches []models.Patch) (models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, serviceId))
	result := &models.Service{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/plans/%s", c.Address, services, serviceId, planId))
	result := &models.ServicePlan{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) AddService(service models.Service) (models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, services))
	result := &models.Service{}
	err := brokerHttp.AddModel(connector, service, http.StatusCreated, result)
	return *result, err
}
