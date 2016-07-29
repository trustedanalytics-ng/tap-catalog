package client

import (
	"errors"
	"fmt"
	"net/http"

	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

func (c *TapCatalogApiConnector) GetCatalogHealth() error {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, healthz))
	status, _, err := brokerHttp.RestGET(connector.Url, connector.BasicAuth, connector.Client)
	if status != http.StatusOK {
		err = errors.New("Invalid health status: " + string(status))
	}
	return err
}
