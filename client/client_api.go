package client

import (
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
	AddApplication(application models.Application) (models.Application, error)
	GetApplication(applicationId string) (models.Application, error)
	UpdateApplication(applicationId string, patches []models.Patch) (models.Application, error)
	ListApplications() ([]models.Application, error)
	AddTemplate(template models.Template) (models.Template, error)
	AddService(service models.Service) (models.Service, error)
	AddImage(image models.Image) (models.Image, error)
	GetImage(imageId string) (models.Image, error)
	UpdateImage(imageId string, patches []models.Patch) (models.Image, error)
	GetServices() ([]models.Service, error)
	AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, error)
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
