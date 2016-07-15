package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) AddTemplate(template models.Template) (models.Template, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, templates))
	result := &models.Template{}
	err := brokerHttp.AddModel(connector, template, http.StatusCreated, result)
	return *result, err
}

func (c *TapCatalogApiConnector) UpdateTemplate(templateId string, patches []models.Patch) (models.Template, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, templates, templateId))
	result := &models.Template{}
	err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, err
}
