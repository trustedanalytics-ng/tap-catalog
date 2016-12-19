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

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) Plans(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	services, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		util.HandleError(rw, err)
		return
	}
	plans := []models.ServicePlan{}
	if services.(models.Service).Plans != nil {
		plans = services.(models.Service).Plans
	}
	util.WriteJson(rw, plans, http.StatusOK)
}

func (c *Context) GetPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	key := c.mapper.ToKey(c.getServicePlansDir(serviceId), planId)

	result, err := c.repository.GetData(key, models.ServicePlan{})
	util.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) AddPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	_, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		util.HandleError(rw, err)
		return
	}

	reqPlan := &models.ServicePlan{}
	err = util.ReadJson(req, reqPlan)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqPlan)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	if reqPlan.Id, err = c.reserveID(c.getServicePlansDir(serviceId)); err != nil {
		util.Respond500(rw, err)
		return
	}

	planKeyStore := c.mapper.ToKeyValue(c.getServicePlansDir(serviceId), reqPlan, true)
	err = c.repository.CreateData(planKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	plan, err := c.repository.GetData(c.getServicedPlanIDKey(serviceId, reqPlan.Id), models.ServicePlan{})
	util.WriteJsonOrError(rw, plan, http.StatusCreated, err)
}

func (c *Context) PatchPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	plan, err := c.repository.GetData(c.getServicedPlanIDKey(serviceId, planId), models.ServicePlan{})
	if err != nil {
		util.HandleError(rw, err)
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.getServicedPlanIDKey(serviceId, planId), models.ServicePlan{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	plan, err = c.repository.GetData(c.getServicedPlanIDKey(serviceId, planId), models.ServicePlan{})
	util.WriteJsonOrError(rw, plan, http.StatusOK, err)
}

func (c *Context) DeletePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]
	err := c.repository.DeleteData(c.getServicedPlanIDKey(serviceId, planId))
	util.WriteJsonOrError(rw, "", http.StatusNoContent, err)
}

func (c *Context) getServicedPlanIDKey(serviceId, planId string) string {
	return c.getServicePlansDir(serviceId) + "/" + planId
}

func (c *Context) getServicePlansDir(serviceId string) string {
	return c.getServiceKey() + "/" + serviceId + "/" + data.Plans
}
