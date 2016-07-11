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

	"github.com/trustedanalytics/tapng-catalog/data"
	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-go-common/util"
)

func (c *Context) Applications(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Applications, models.Application{})
	if err != nil {
		util.Respond500(rw, err)
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]

	result, err := c.repository.GetData(c.buildApplicationKey(applicationId), models.Application{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddApplication(rw web.ResponseWriter, req *web.Request) {
	reqApplication := &models.Application{}

	err := util.ReadJson(req, reqApplication)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqApplication)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	applicationKeyStore := c.mapper.ToKeyValue(data.Applications, reqApplication, true)
	err = c.repository.StoreData(applicationKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	application, err := c.repository.GetData(c.buildApplicationKey(reqApplication.Id), models.Application{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, application, http.StatusCreated)
}

func (c *Context) PatchApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	application, err := c.repository.GetData(c.buildApplicationKey(applicationId), models.Application{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildApplicationKey(applicationId), models.Application{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	application, err = c.repository.GetData(c.buildApplicationKey(applicationId), models.Application{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, application, http.StatusOK)
}

func (c *Context) DeleteApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	err := c.repository.DeleteData(c.buildApplicationKey(applicationId))
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildApplicationKey(applicationId string) string {
	return c.mapper.ToKey(data.Applications, applicationId)
}
