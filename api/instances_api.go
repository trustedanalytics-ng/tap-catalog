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

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-catalog/webutils"
)

func (c *Context) Instances(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Instances, data.Instances)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]

	result, err := c.repository.GetData(data.Instances, c.buildInstanceKey(instanceId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddInstance(rw web.ResponseWriter, req *web.Request) {
	reqInstance := models.Instance{}

	instanceId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	err = webutils.ReadJson(req, &reqInstance)

	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	reqInstance.Id = instanceId.String()

	instanceKeyStore := map[string]interface{}{}

	instanceKeyStore = c.mapper.ToKeyValue(data.Instances, reqInstance)

	err = c.repository.StoreData(instanceKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, reqInstance, http.StatusCreated)
}

func (c *Context) PatchInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	instance, err := c.repository.GetData(data.Instances, c.buildInstanceKey(instanceId))
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

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildInstanceKey(instanceId), models.Instance{}, patches)
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

	instance, err = c.repository.GetData(data.Instances, c.buildInstanceKey(instanceId))
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, instance, http.StatusOK)
}

func (c *Context) DeleteInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	err := c.repository.DeleteData(c.buildInstanceKey(instanceId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildInstanceKey(instanceId string) string {
	return c.mapper.ToKey(data.Instances, instanceId)
}

func (c *Context) AddInstanceBinding(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Add Instance Binding", http.StatusCreated)
}

func (c *Context) DeleteInstanceBinding(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Instance Binding", http.StatusNoContent)
}

func (c *Context) AddInstanceMetadata(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Add Instance Binding", http.StatusCreated)
}

func (c *Context) DeleteInstanceMetadata(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Instance Binding", http.StatusNoContent)
}
