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

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
)

func (c *Context) Plans(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	webutils.WriteJson(rw, serviceId, http.StatusOK)
}

func (c *Context) GetPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	key := c.mapper.ToKey(buildHomeDir(serviceId), planId)

	result, err := c.repository.GetData(data.Plans, key)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	reqPlan := models.ServicePlan{}
	err = webutils.ReadJson(req, &reqPlan)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	reqPlan.Id = planId.String()

	planKeyStore := map[string]interface{}{}

	planKeyStore = c.mapper.ToKeyValue(buildHomeDir(serviceId), reqPlan)

	err = c.repository.StoreData(planKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, reqPlan, http.StatusCreated)
}

func (c *Context) UpdatePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	reqPlan := models.ServicePlan{}
	err := webutils.ReadJson(req, &reqPlan)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	reqPlan.Id = planId

	planKeyStore := map[string]interface{}{}

	planKeyStore = c.mapper.ToKeyValue(buildHomeDir(serviceId), reqPlan)

	err = c.repository.StoreData(planKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, reqPlan, http.StatusOK)
}

func (c *Context) DeletePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]
	err := c.repository.DeleteData(c.buildPlanKey(serviceId, planId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildPlanKey(serviceId, planId string) string {
	return buildHomeDir(serviceId) + "/" + planId
}

func buildHomeDir(serviceId string) string {
	return data.Services + "/" + serviceId + data.Plans
}
