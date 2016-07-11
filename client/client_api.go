package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

type TapCatalogApi interface {
	GetInstance(instanceId string) (models.Instance, error)
	UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, error)
	GetService(serviceId string) (models.Service, error)
	UpdateService(serviceId string, patches []models.Patch) (models.Service, error)
	UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, error)
	GetApplication(applicationId string) (models.Application, error)
	UpdateApplication(applicationId string, patches []models.Patch) (models.Application, error)
	AddTemplate(template models.Template) (models.Template, error)
	AddService(service models.Service) (models.Service, error)
	GetImage(imageId string) (models.Image, error)
	UpdateImage(imageId string, patches []models.Patch) (models.Image, error)
	GetServices() ([]models.Service, error)
	AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, error)
}

type TapCatalogApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

const (
	apiVersion   = "v1"
	instances    = apiVersion + "/instances"
	services     = apiVersion + "/services"
	applications = apiVersion + "/applications"
	templates    = apiVersion + "/templates"
	images       = apiVersion + "/images"
)

func NewTapCatalogApiWithBasicAuth(address, username, password string) (*TapCatalogApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithBasicAuth()
	if err != nil {
		return nil, err
	}
	return &TapCatalogApiConnector{address, username, password, client}, nil
}

func NewTapCatalogApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TapCatalogApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapCatalogApiConnector{address, username, password, client}, nil
}

func (c *TapCatalogApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
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

func (c *TapCatalogApiConnector) AddTemplate(template models.Template) (models.Template, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, templates))
	result := &models.Template{}
	err := brokerHttp.PatchModel(connector, template, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) AddService(service models.Service) (models.Service, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, services))
	result := &models.Service{}
	err := brokerHttp.PatchModel(connector, service, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) GetImage(imageId string) (models.Image, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, images, imageId))
	result := &models.Image{}
	err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdateImage(imageId string, patches []models.Patch) (models.Image, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, images, imageId))
	result := &models.Image{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}

func (c *TapCatalogApiConnector) AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s/instances", c.Address, services, serviceId))
	result := &models.Instance{}
	err := brokerHttp.PatchModel(connector, instance, http.StatusCreated, result)
	return *result, err
}
