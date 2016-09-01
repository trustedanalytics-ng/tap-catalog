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
	"errors"
	"net/http"

	"github.com/gocraft/web"
	"github.com/looplab/fsm"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(c.getTemplateKey(), models.Template{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	result, err := c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := &models.Template{}

	err := util.ReadJson(req, reqTemplate)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqTemplate)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	reqTemplate.State = models.TemplateStateInProgress
	templateKeyStore := c.mapper.ToKeyValue(c.getTemplateKey(), reqTemplate, true)
	err = c.repository.StoreData(templateKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	template, err := c.repository.GetData(c.buildTemplateKey(reqTemplate.Id), models.Template{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, template, http.StatusCreated)
}

func (c *Context) DeleteTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	err := c.repository.DeleteData(c.buildTemplateKey(templateId))
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	util.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) PatchTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	templateInt, err := c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	template, ok := templateInt.(models.Template)
	if !ok {
		util.Respond500(rw, errors.New("Template retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.allowStateChange(patches, c.getTemplatesFSM(template.State))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildTemplateKey(templateId), models.Template{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	templateInt, err = c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, templateInt, http.StatusOK)
}

func (c *Context) getTemplateKey() string {
	return c.mapper.ToKey(c.organization, data.Templates)
}

func (c *Context) buildTemplateKey(templateId string) string {
	return c.mapper.ToKey(c.getTemplateKey(), templateId)
}

func (c *Context) getTemplatesFSM(initialState models.TemplateState) *fsm.FSM {
	return fsm.NewFSM(string(initialState),
		fsm.Events{
			{Name: "READY", Src: []string{"IN_PROGRESS"}, Dst: "READY"},
			{Name: "UNAVAILABLE", Src: []string{"IN_PROGRESS"}, Dst: "UNAVAILABLE"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
