package client

import (
	"errors"
	"fmt"
	"net/http"

	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
	"strconv"
)

func (c *TapCatalogApiConnector) GetCatalogHealth() (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/%s", c.Address, healthz))
	status, _, err := brokerHttp.RestGET(connector.Url, brokerHttp.GetBasicAuthHeader(connector.BasicAuth), connector.Client)
	if status != http.StatusOK {
		err = errors.New("Invalid health status: " + strconv.Itoa(status))
	}
	return status, err
}
