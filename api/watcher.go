package api

import (
	"net/http"
	"strconv"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) LatestIndex(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetLatestIndex(c.mapper.ToKey("", c.organization))
	util.WriteJsonOrError(rw, models.Index{Latest: result}, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) monitorSpecificState(rw web.ResponseWriter, req *web.Request, key string) {
	afterIndex, err := strconv.ParseUint(req.URL.Query().Get("afterIndex"), 10, 32)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	result, err := c.repository.MonitorObjectsStates(key, afterIndex)
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}
