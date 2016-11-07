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
	"fmt"
	"net/http"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

const keyNotFoundMessage = "Key not found"

func (c *Context) Applications(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(c.getApplicationKey(), models.Application{})
	util.WriteJsonOrError(rw, result, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) getApplication(id string) (models.Application, error) {
	entity, err := c.repository.GetData(c.buildApplicationKey(id), models.Application{})
	if err != nil {
		err = fmt.Errorf("application %q retrieval failed: %v", id, err)
		logger.Warning(err)
		return models.Application{}, err
	}

	application, ok := entity.(models.Application)
	if !ok {
		err = fmt.Errorf("type assertion for application %q failed: object from database: %v", id, entity)
		logger.Error(err)
		return models.Application{}, err
	}

	return application, nil
}

func (c *Context) GetApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]

	app, err := c.getApplication(applicationId)
	util.WriteJsonOrError(rw, app, getHttpStatusOrStatusError(http.StatusOK, err), err)
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

	err = data.CheckIfMatchingRegexp(reqApplication.Name, data.RegexpDnsLabelLowercase)
	if err != nil {
		util.Respond400(rw, errors.New("Field: Name has incorrect value:"+reqApplication.Name))
		return
	}

	exists, err := c.repository.IsExistByName(reqApplication.Name, models.Application{}, c.getApplicationKey())
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	if exists {
		util.Respond409(rw, errors.New("application with name: "+reqApplication.Name+" already exists!"))
		return
	}

	applicationKeyStore := c.mapper.ToKeyValue(c.getApplicationKey(), reqApplication, true)
	err = c.repository.StoreData(applicationKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	application, err := c.repository.GetData(c.buildApplicationKey(reqApplication.Id), models.Application{})
	util.WriteJsonOrError(rw, application, getHttpStatusOrStatusError(http.StatusCreated, err), err)
}

func (c *Context) PatchApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	application, err := c.repository.GetData(c.buildApplicationKey(applicationId), models.Application{})
	if err != nil {
		handleGetDataError(rw, err)
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
	util.WriteJsonOrError(rw, application, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) DeleteApplication(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]
	err := c.repository.DeleteData(c.buildApplicationKey(applicationId))
	util.WriteJsonOrError(rw, "", getHttpStatusOrStatusError(http.StatusNoContent, err), err)
}

func (c *Context) getApplicationKey() string {
	org := c.mapper.ToKey("", c.organization)
	return c.mapper.ToKey(org, data.Applications)
}

func (c *Context) buildApplicationKey(applicationId string) string {
	return c.mapper.ToKey(c.getApplicationKey(), applicationId)
}
