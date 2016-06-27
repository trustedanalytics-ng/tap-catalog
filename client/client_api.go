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
	UpdateInstance(instanceId string, instance models.Instance) (models.Instance, error)
	GetService(serviceId string) (models.Service, error)
	UpdateService(serviceId string, instance models.Service) (models.Service, error)
	UpdatePlan(serviceId, planId string, instance models.ServicePlan) (models.ServicePlan, error)
	GetApplication(applicationId string) (models.Application, error)
	UpdateApplication(applicationId string, instance models.Application) (models.Application, error)
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

	url := fmt.Sprintf("%s/instances/%s", c.Address, instanceId)
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

func (c *TapCatalogApiConnector) UpdateInstance(instanceId string, instance models.Instance) (models.Instance, error) {
	result := models.Instance{}

	reqBody, err := json.Marshal(instance)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/instances/%s", c.Address, instanceId)
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

	url := fmt.Sprintf("%s/services/%s", c.Address, serviceId)
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

func (c *TapCatalogApiConnector) UpdateService(serviceId string, instance models.Service) (models.Service, error) {
	result := models.Service{}

	reqBody, err := json.Marshal(instance)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/services/%s", c.Address, serviceId)
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

func (c *TapCatalogApiConnector) UpdatePlan(serviceId, planId string, instance models.ServicePlan) (models.ServicePlan, error) {
	result := models.ServicePlan{}

	reqBody, err := json.Marshal(instance)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/services/%s/plans/%s", c.Address, serviceId, planId)
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

	url := fmt.Sprintf("%s/applications/%s", c.Address, applicationId)
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

func (c *TapCatalogApiConnector) UpdateApplication(applicationId string, instance models.Application) (models.Application, error) {
	result := models.Application{}

	reqBody, err := json.Marshal(instance)
	if err != nil {
		return result, err
	}

	url := fmt.Sprintf("%s/applications/%s", c.Address, applicationId)
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
