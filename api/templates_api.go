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
	"github.com/trustedanalytics/tapng-catalog/webutils"
	"github.com/trustedanalytics/tapng-go-common/logger"
)

var logger = logger_wrapper.InitLogger("templates_api")

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Templates, &models.Template{})
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	result, err := c.repository.GetData(c.buildTemplateKey(templateId), &models.Template{})
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

	webutils.WriteJson(rw, reqTemplate, http.StatusCreated)
}

func (c *Context) DeleteTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	err := c.repository.DeleteData(c.buildTemplateKey(templateId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) PatchTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	template, err := c.repository.GetData(c.buildTemplateKey(templateId), &models.Template{})
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

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildTemplateKey(templateId), &models.Template{}, patches)
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

	template, err = c.repository.GetData(c.buildTemplateKey(templateId), &models.Template{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, template, http.StatusOK)
}

func (c *Context) buildTemplateKey(templateId string) string {
	return c.mapper.ToKey(data.Templates, templateId)
}
