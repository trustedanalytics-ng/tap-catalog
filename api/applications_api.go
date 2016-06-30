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
	"github.com/nu7hatch/gouuid"
	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
	"net/http"
)

func (c *Context) Applications(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Applications, data.Applications)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]

	result, err := c.repository.GetData(data.Applications, c.buildApplicationKey(applicationId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddApplication(rw web.ResponseWriter, req *web.Request) {
	reqApplication := models.Application{}

	err := webutils.ReadJson(req, &reqApplication)
	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	applicationId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	reqApplication.Id = applicationId.String()

	applicationKeyStore := map[string]interface{}{}

	applicationKeyStore = c.mapper.ToKeyValue(data.Applications, reqApplication)

	err = c.repository.StoreData(applicationKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
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

	applicationKeyStore := map[string]interface{}{}

	applicationKeyStore = c.mapper.ToKeyValue(data.Applications, reqApplication)

	err = c.repository.StoreData(applicationKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, reqApplication, http.StatusOK)
}

func (c *Context) DeleteApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	err := c.repository.DeleteData(c.buildApplicationKey(applicationId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildApplicationKey(applicationId string) string {
	return c.mapper.ToKey(data.Applications, applicationId)
}
