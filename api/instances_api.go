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
	webutils.WriteJson(rw, "List Instances", http.StatusOK)
}

func (c *Context) GetInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]

	result, err := c.repository.GetData(data.Instances, instanceId)
	if err != nil {
		webutils.Respond500(rw, err)
	}

	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddInstance(rw web.ResponseWriter, req *web.Request) {
	reqInstance := models.Instance{}

	serviceId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
	}
	err = webutils.ReadJson(req, &reqInstance)

	if err != nil {
		webutils.Respond400(rw, err)
	}

	reqInstance.Id = serviceId.String()

	serviceKeyStore := map[string]interface{}{}

	serviceKeyStore = c.mapper.ToKeyValue(data.Instances, reqInstance)

	err = c.repository.StoreData(serviceKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
	}
	webutils.WriteJson(rw, reqInstance, http.StatusCreated)
}

func (c *Context) UpdateInstance(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Update Instance", http.StatusOK)
}

func (c *Context) DeleteInstance(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Instance", http.StatusNoContent)
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
