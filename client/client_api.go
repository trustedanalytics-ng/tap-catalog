package client

import (
	"net/http"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

type TapCatalogApi interface {
	AddApplication(application models.Application) (models.Application, int, error)
	AddImage(image models.Image) (models.Image, int, error)
	AddService(service models.Service) (models.Service, int, error)
	AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, int, error)
	AddServiceBrokerInstance(serviceId string, instance models.Instance) (models.Instance, int, error)
	AddApplicationInstance(applicationId string, instance models.Instance) (models.Instance, int, error)
	AddTemplate(template models.Template) (models.Template, int, error)
	GetApplication(applicationId string) (models.Application, int, error)
	GetCatalogHealth() (int, error)
	GetImage(imageId string) (models.Image, int, error)
	GetInstance(instanceId string) (models.Instance, int, error)
	GetInstanceBindings(instanceId string) ([]models.Instance, int, error)
	GetService(serviceId string) (models.Service, int, error)
	GetServices() ([]models.Service, int, error)
	ListApplications() ([]models.Application, int, error)
	ListApplicationsInstances() ([]models.Instance, int, error)
	ListInstances() ([]models.Instance, int, error)
	ListServicesInstances() ([]models.Instance, int, error)
	UpdateApplication(applicationId string, patches []models.Patch) (models.Application, int, error)
	UpdateImage(imageId string, patches []models.Patch) (models.Image, int, error)
	UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, int, error)
	UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, int, error)
	UpdateService(serviceId string, patches []models.Patch) (models.Service, int, error)
	UpdateTemplate(templateId string, patches []models.Patch) (models.Template, int, error)
	DeleteApplication(applicationId string) (int, error)
	DeleteImage(imageId string) (int, error)
	DeleteInstance(instanceId string) (int, error)
}

type TapCatalogApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

const (
	apiPrefix        = "api/"
	apiVersion       = "v1"
	instances        = apiPrefix + apiVersion + "/instances"
	instanceBindings = instances + "/bindings"
	services         = apiPrefix + apiVersion + "/services"
	applications     = apiPrefix + apiVersion + "/applications"
	templates        = apiPrefix + apiVersion + "/templates"
	images           = apiPrefix + apiVersion + "/images"
	healthz          = "healthz"
)

func NewTapCatalogApiWithBasicAuth(address, username, password string) (*TapCatalogApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClient()
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
