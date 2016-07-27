package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tapng-catalog/models"
	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) AddImage(image models.Image) (models.Image, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, images))
	result := &models.Image{}
	err := brokerHttp.AddModel(connector, image, http.StatusCreated, result)
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
