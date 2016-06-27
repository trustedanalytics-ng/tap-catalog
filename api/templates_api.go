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
	"github.com/trustedanalytics/tap-catalog/webutils"
	"github.com/trustedanalytics/tap-go-common/logger"
)

var logger = logger_wrapper.InitLogger("templates_api")

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "List Templates", http.StatusOK)
}

func (c *Context) GetTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	result, err := c.repository.GetData(data.Templates, templateId)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := models.Template{}

	err := webutils.ReadJson(req, &reqTemplate)
	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	templateKeyStore := map[string]interface{}{}

	templateKeyStore = c.mapper.ToKeyValue(data.Templates, reqTemplate)

	err = c.repository.StoreData(templateKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	//TODO return data saved in etcd not from request
	webutils.WriteJson(rw, reqTemplate, http.StatusCreated)
}

func (c *Context) DeleteTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	webutils.WriteJson(rw, templateId, http.StatusNoContent)
}

func (c *Context) UpdateTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := models.Template{}
	templateId := req.PathParams["templateId"]

	reqTemplate.Id = templateId

	err := webutils.ReadJson(req, &reqTemplate)
	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	webutils.WriteJson(rw, reqTemplate, http.StatusOK)
}
