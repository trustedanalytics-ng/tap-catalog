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

	"github.com/trustedanalytics-ng/tap-catalog/data"
	"github.com/trustedanalytics-ng/tap-catalog/models"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func (c *Context) Templates(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(c.getTemplateKey(), models.Template{})
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) GetTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	result, err := c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) AddTemplate(rw web.ResponseWriter, req *web.Request) {
	reqTemplate := &models.Template{}

	err := commonHttp.ReadJson(req, reqTemplate)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqTemplate)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	if reqTemplate.Id, err = c.reserveID(c.getTemplateKey()); err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	reqTemplate.State = models.TemplateStateInProgress
	templateKeyStore := c.mapper.ToKeyValue(c.getTemplateKey(), reqTemplate, true)
	err = c.repository.CreateData(templateKeyStore)
	if err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	template, err := c.repository.GetData(c.buildTemplateKey(reqTemplate.Id), models.Template{})
	commonHttp.WriteJsonOrError(rw, template, http.StatusCreated, err)
}

func (c *Context) DeleteTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]

	err := c.repository.DeleteData(c.buildTemplateKey(templateId))
	commonHttp.WriteJsonOrError(rw, "", http.StatusNoContent, err)
}

func (c *Context) PatchTemplate(rw web.ResponseWriter, req *web.Request) {
	templateId := req.PathParams["templateId"]
	templateInt, err := c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	template, ok := templateInt.(models.Template)
	if !ok {
		commonHttp.HandleError(rw, errors.New("template retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = commonHttp.ReadJson(req, &patches)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	fsmFunc := func() *fsm.FSM {
		return c.getTemplatesFSM(template.State)
	}
	if err = c.handleFsm(rw, req, patches, fsmFunc); err != nil {
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildTemplateKey(templateId), models.Template{}, patches)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	templateInt, err = c.repository.GetData(c.buildTemplateKey(templateId), models.Template{})
	commonHttp.WriteJsonOrError(rw, templateInt, http.StatusOK, err)
}

func (c *Context) getTemplateKey() string {
	org := c.mapper.ToKey("", c.organization)
	return c.mapper.ToKey(org, data.Templates)
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
