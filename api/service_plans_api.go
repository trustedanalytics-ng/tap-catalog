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

	"github.com/nu7hatch/gouuid"
	"github.com/trustedanalytics/tap-catalog/api/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
)

func (c *Context) Plans(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	webutils.WriteJson(rw, serviceId, http.StatusOK)
}

func (c *Context) GetPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]
	webutils.WriteJson(rw, serviceId+planId, http.StatusOK)
}

func (c *Context) AddPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
	}

	reqPlan := models.ServicePlan{}
	webutils.ReadJson(req, &reqPlan)
	reqPlan.Id = planId.String()

	webutils.WriteJson(rw, planId.String()+serviceId, http.StatusCreated)
}

func (c *Context) UpdatePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	webutils.WriteJson(rw, planId+serviceId, http.StatusOK)
}

func (c *Context) DeletePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	webutils.WriteJson(rw, planId+serviceId, http.StatusNoContent)
}
