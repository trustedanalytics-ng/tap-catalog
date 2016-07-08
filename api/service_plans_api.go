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

	"github.com/trustedanalytics/tapng-catalog/data"
	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/util"
)

func (c *Context) Plans(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	result, err := c.repository.GetListOfData(buildHomeDir(serviceId), models.ServicePlan{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	key := c.mapper.ToKey(buildHomeDir(serviceId), planId)

	result, err := c.repository.GetData(key, models.ServicePlan{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId, err := uuid.NewV4()
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	reqPlan := models.ServicePlan{}
	err = util.ReadJson(req, &reqPlan)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	reqPlan.Id = planId.String()
	planKeyStore := c.mapper.ToKeyValue(buildHomeDir(serviceId), reqPlan, true)

	err = c.repository.StoreData(planKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	plan, err := c.repository.GetData(c.buildPlanKey(serviceId, reqPlan.Id), models.ServicePlan{})
	if err != nil {
		logger.Error(err)
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, plan, http.StatusCreated)
}

func (c *Context) PatchPlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]

	plan, err := c.repository.GetData(c.buildPlanKey(serviceId, planId), models.ServicePlan{})
	if err != nil {
		logger.Error(err)
		util.Respond500(rw, err)
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildPlanKey(serviceId, planId), models.ServicePlan{}, patches)
	if err != nil {
		logger.Error(err)
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		logger.Error(err)
		util.Respond500(rw, err)
		return
	}

	plan, err = c.repository.GetData(c.buildPlanKey(serviceId, planId), models.ServicePlan{})
	if err != nil {
		logger.Error(err)
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, plan, http.StatusOK)
}

func (c *Context) DeletePlan(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	planId := req.PathParams["planId"]
	err := c.repository.DeleteData(c.buildPlanKey(serviceId, planId))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	util.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildPlanKey(serviceId, planId string) string {
	return buildHomeDir(serviceId) + "/" + planId
}

func buildHomeDir(serviceId string) string {
	return data.Services + "/" + serviceId + data.Plans
}
