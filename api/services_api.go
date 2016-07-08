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
	"github.com/trustedanalytics/tapng-catalog/webutils"
)

func (c *Context) Services(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Services, models.Service{})
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	result, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddService(rw web.ResponseWriter, req *web.Request) {
	reqService := models.Service{}

	serviceId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	err = webutils.ReadJson(req, &reqService)

	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	reqService.Id = serviceId.String()
	if reqService.Plans != nil {
		for i, plan := range reqService.Plans {
			planId, err := uuid.NewV4()
			if err != nil {
				webutils.Respond500(rw, err)
			}
			plan.Id = planId.String()
			reqService.Plans[i] = plan
		}
	}

	serviceKeyStore := c.mapper.ToKeyValue(data.Services, reqService, true)
	err = c.repository.StoreData(serviceKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	service, err := c.repository.GetData(c.buildServiceKey(serviceId.String()), models.Service{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, service, http.StatusCreated)
}

func (c *Context) PatchService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	service, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}

	patches, err := webutils.ReadPatch(req)
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildServiceKey(serviceId), models.Service{}, patches)
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}

	service, err = c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, service, http.StatusOK)
}

func (c *Context) DeleteService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	err := c.repository.DeleteData(c.buildServiceKey(serviceId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, serviceId, http.StatusNoContent)
}

func (c *Context) buildServiceKey(serviceId string) string {
	return c.mapper.ToKey(data.Services, serviceId)
}
