package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
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
}

type TapCatalogApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

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

func (c *TapCatalogApiConnector) GetInstance(instanceId string) (models.Instance, error) {
	result := models.Instance{}

	url := fmt.Sprintf("%s/v1/instances/%s", c.Address, instanceId)
	status, body, err := brokerHttp.RestGET(url, &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, error) {
	result := models.Instance{}

	reqBody, err := json.Marshal(patches)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/v1/instances/%s", c.Address, instanceId)
	status, body, err := brokerHttp.RestPATCH(url, string(reqBody), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) GetService(serviceId string) (models.Service, error) {
	result := models.Service{}

	url := fmt.Sprintf("%s/v1/services/%s", c.Address, serviceId)
	status, body, err := brokerHttp.RestGET(url, &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) UpdateService(serviceId string, patches []models.Patch) (models.Service, error) {
	result := models.Service{}

	reqBody, err := json.Marshal(patches)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/v1/services/%s", c.Address, serviceId)
	status, body, err := brokerHttp.RestPATCH(url, string(reqBody), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) UpdatePlan(serviceId, planId string, patches []models.Patch) (models.ServicePlan, error) {
	result := models.ServicePlan{}

	reqBody, err := json.Marshal(patches)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/v1/services/%s/plans/%s", c.Address, serviceId, planId)
	status, body, err := brokerHttp.RestPATCH(url, string(reqBody), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) GetApplication(applicationId string) (models.Application, error) {
	result := models.Application{}

	url := fmt.Sprintf("%s/v1/applications/%s", c.Address, applicationId)
	status, body, err := brokerHttp.RestGET(url, &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) UpdateApplication(applicationId string, patches []models.Patch) (models.Application, error) {
	result := models.Application{}

	reqBody, err := json.Marshal(patches)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/v1/applications/%s", c.Address, applicationId)
	status, body, err := brokerHttp.RestPATCH(url, string(reqBody), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusOK {
		return result, errors.New("Bad response status: " + strconv.Itoa(status) + ". Body: " + string(body))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) AddTemplate(template models.Template) (models.Template, error) {

	result := models.Template{}

	url := fmt.Sprintf("%s/v1/templates", c.Address)
	b, err := json.Marshal(&template)
	if err != nil {
		return result, err
	}
	status, body, err := brokerHttp.RestPOST(url, string(b), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusCreated {
		return result, errors.New("Bad response status: " + strconv.Itoa(status))
	}
	return result, nil
}

func (c *TapCatalogApiConnector) AddService(service models.Service) (models.Service, error) {

	result := models.Service{}

	url := fmt.Sprintf("%s/v1/services", c.Address)
	b, err := json.Marshal(&service)
	if err != nil {
		return result, err
	}
	status, body, err := brokerHttp.RestPOST(url, string(b), &brokerHttp.BasicAuth{c.Username, c.Password}, c.Client)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if status != http.StatusCreated {
		return result, errors.New("Bad response status: " + strconv.Itoa(status))
	}
	return result, nil
}