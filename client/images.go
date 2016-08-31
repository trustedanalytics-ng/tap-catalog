package client

import (
	"fmt"
	"net/http"

	"github.com/trustedanalytics/tap-catalog/models"
	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *TapCatalogApiConnector) AddImage(image models.Image) (models.Image, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, images))
	result := &models.Image{}
	status, err := brokerHttp.PostModel(connector, image, http.StatusCreated, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) GetImage(imageId string) (models.Image, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, images, imageId))
	result := &models.Image{}
	status, err := brokerHttp.GetModel(connector, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) UpdateImage(imageId string, patches []models.Patch) (models.Image, int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, images, imageId))
	result := &models.Image{}
	status, err := brokerHttp.PatchModel(connector, patches, http.StatusOK, result)
	return *result, status, err
}

func (c *TapCatalogApiConnector) DeleteImage(imageId string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s/%s", c.Address, images, imageId))
	return brokerHttp.DeleteModel(connector, http.StatusNoContent)
}
