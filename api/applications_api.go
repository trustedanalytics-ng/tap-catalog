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
	"github.com/gocraft/web"
	"github.com/trustedanalytics/tap-catalog/api/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
	"net/http"
)

func (c *Context) Applications(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "List Applications", http.StatusOK)
}

func (c *Context) GetApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Single Application", http.StatusOK)
}

func (c *Context) AddApplication(rw web.ResponseWriter, req *web.Request) {
	reqApplication := models.Application{}

	err := webutils.ReadJson(req, &reqApplication)
	if err != nil {
		webutils.Respond400(rw, err)
	}
	webutils.WriteJson(rw, reqApplication, http.StatusCreated)
}

func (c *Context) UpdateApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	reqApplication := models.Application{}
	reqApplication.Id = applicationId

	err := webutils.ReadJson(req, &reqApplication)
	if err != nil {
		webutils.Respond400(rw, err)
	}
	webutils.WriteJson(rw, reqApplication, http.StatusOK)
}

func (c *Context) DeleteApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Application", http.StatusNoContent)
}
