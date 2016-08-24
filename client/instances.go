package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) ListApplicationsInstances() ([]models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, applications, "instances"))
	result := &[]models.Instance{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) ListServicesInstances() ([]models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, services, "instances"))
	result := &[]models.Instance{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) ListInstances() ([]models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, instances))
	result := &[]models.Instance{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) GetInstance(instanceId string) (models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	result := &models.Instance{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) GetInstanceBindings(instanceId string) ([]models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instanceBindings, instanceId))
	result := &[]models.Instance{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	result := &models.Instance{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) AddServiceBrokerInstance(serviceId string, instance models.Instance) (models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances?isServiceBroker=true", c.Address, services, serviceId))
	result := &models.Instance{}
	status, err := brokerHttp.PostModel(connector, instance, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances?isServiceBroker=false", c.Address, services, serviceId))
	result := &models.Instance{}
	status, err := brokerHttp.PostModel(connector, instance, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) AddApplicationInstance(applicationId string, instance models.Instance) (models.Instance, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances", c.Address, applications, applicationId))
	result := &models.Instance{}
	status, err := brokerHttp.PostModel(connector, instance, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) DeleteInstance(instanceId string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, instances, instanceId))
	status, err := brokerHttp.DeleteModel(connector, http.StatusNoContent)
	return status, err
}
