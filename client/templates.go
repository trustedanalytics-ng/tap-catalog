package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *TapCatalogApiConnector) AddTemplate(template models.Template) (models.Template, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, templates))
	result := &models.Template{}
	status, err := brokerHttp.PostModel(connector, template, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdateTemplate(templateId string, patches []models.Patch) (models.Template, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, templates, templateId))
	result := &models.Template{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}
