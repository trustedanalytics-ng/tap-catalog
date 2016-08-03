package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) ListApplicationsInstances() ([]models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, "instances"))
	result := &[]models.Instance{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) ListServicesInstances() ([]models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, "instances"))
	result := &[]models.Instance{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) ListInstances() ([]models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, instances))
	result := &[]models.Instance{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) GetInstance(instanceId string) (models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	result := &models.Instance{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	result := &models.Instance{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances", c.Address, services, serviceId))
	result := &models.Instance{}
	err := brokerHttp.PostModel(connector, instance, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) AddApplicationInstance(applicationId string, instance models.Instance) (models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances", c.Address, applications, applicationId))
	result := &models.Instance{}
	err := brokerHttp.PostModel(connector, instance, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) DeleteInstance(instanceId string) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return err
}
