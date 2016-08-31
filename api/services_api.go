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
	"github.com/looplab/fsm"

	"errors"
	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) Services(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Services, models.Service{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	result, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddService(rw web.ResponseWriter, req *web.Request) {
	reqService := &models.Service{}

	err := util.ReadJson(req, reqService)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqService)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfDNSLabelLowercaseCompatible(reqService.Name, "Name")
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	exists, err := c.repository.IsExistByName(reqService.Name, models.Service{}, data.Services)
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	if exists {
		util.Respond409(rw, errors.New("service with name: "+reqService.Name+" already exists!"))
		return
	}

	reqService.State = models.ServiceStateDeploying
	serviceKeyStore := c.mapper.ToKeyValue(data.Services, reqService, true)
	err = c.repository.StoreData(serviceKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	service, err := c.repository.GetData(c.buildServiceKey(reqService.Id), models.Service{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, service, http.StatusCreated)
}

func (c *Context) PatchService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	serviceInt, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	service, ok := serviceInt.(models.Service)
	if !ok {
		util.Respond500(rw, errors.New("Service retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.allowStateChange(patches, c.getServicesFSM(service.State))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildServiceKey(serviceId), models.Service{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	serviceInt, err = c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, serviceInt, http.StatusOK)
}

func (c *Context) DeleteService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	err := c.repository.DeleteData(c.buildServiceKey(serviceId))
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, serviceId, http.StatusNoContent)
}

func (c *Context) buildServiceKey(serviceId string) string {
	return c.mapper.ToKey(data.Services, serviceId)
}

func (c *Context) getServicesFSM(initialState models.ServiceState) *fsm.FSM {
	return fsm.NewFSM(string(initialState),
		fsm.Events{
			{Name: "READY", Src: []string{"DEPLOYING"}, Dst: "READY"},
			{Name: "OFFLINE", Src: []string{"DEPLOYING"}, Dst: "OFFLINE"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
