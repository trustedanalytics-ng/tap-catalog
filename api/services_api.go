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

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/api/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
)

func (c *Context) Services(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "List Services", http.StatusOK)
}

func (c *Context) GetService(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Single Service", http.StatusOK)
}

func (c *Context) AddService(rw web.ResponseWriter, req *web.Request) {
	reqService := models.Service{}

	err := webutils.ReadJson(req, &reqService)
	if err != nil {
		webutils.Respond400(rw, err)
	}
	webutils.WriteJson(rw, reqService, http.StatusCreated)
}

func (c *Context) UpdateService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	reqService := models.Service{}
	reqService.Id = serviceId

	err := webutils.ReadJson(req, &reqService)
	if err != nil {
		webutils.Respond400(rw, err)
	}
	webutils.WriteJson(rw, reqService, http.StatusOK)
}

func (c *Context) DeleteService(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Service", http.StatusNoContent)
}
