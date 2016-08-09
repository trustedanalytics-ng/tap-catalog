package client

import (
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

type TapCatalogApi interface {
	AddApplication(application models.Application) (models.Application, error)
	AddImage(image models.Image) (models.Image, error)
	AddService(service models.Service) (models.Service, error)
	AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, error)
	AddApplicationInstance(applicationId string, instance models.Instance) (models.Instance, error)
	AddTemplate(template models.Template) (models.Template, error)
	DeleteInstance(instanceId string) error
	DeleteApplication(applicationId string) error
	DeleteImage(imageId string) error
	GetApplication(applicationId string) (models.Application, error)
	GetCatalogHealth() error
	GetImage(imageId string) (models.Image, error)
	GetInstance(instanceId string) (models.Instance, error)
	GetService(serviceId string) (models.Service, error)
	GetServices() ([]models.Service, error)
	ListApplications() ([]models.Application, error)
	ListApplicationsInstances() ([]models.Instance, error)
	ListInstances() ([]models.Instance, error)
	ListServicesInstances() ([]models.Instance, error)
	UpdateApplication(applicationId string, patches []models.Patch) (models.Application, error)
	UpdateImage(imageId string, patches []models.Patch) (models.Image, error)
	UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, error)
	UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, error)
	UpdateService(serviceId string, patches []models.Patch) (models.Service, error)
	UpdateTemplate(templateId string, patches []models.Patch) (models.Template, error)
}

type TapCatalogApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

const (
	apiPrefix    = "api/"
	apiVersion   = "v1"
	instances    = apiPrefix + apiVersion + "/instances"
	services     = apiPrefix + apiVersion + "/services"
	applications = apiPrefix + apiVersion + "/applications"
	templates    = apiPrefix + apiVersion + "/templates"
	images       = apiPrefix + apiVersion + "/images"
	healthz      = apiPrefix + apiVersion + "/healthz"
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
