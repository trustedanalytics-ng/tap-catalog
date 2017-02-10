/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"net/http"
	"strconv"

	"github.com/gocraft/web"

	"github.com/trustedanalytics-ng/tap-catalog/models"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func (c *Context) LatestIndex(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetLatestIndex(c.mapper.ToKey("", c.organization))
	commonHttp.WriteJsonOrError(rw, models.Index{Latest: result}, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) monitorSpecificState(rw web.ResponseWriter, req *web.Request, key string) {
	afterIndex, err := strconv.ParseUint(req.URL.Query().Get("afterIndex"), 10, 32)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	result, err := c.repository.MonitorObjectsStates(key, afterIndex)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}
	commonHttp.WriteJson(rw, result, http.StatusOK)
}
