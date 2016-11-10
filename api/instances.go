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
	"github.com/looplab/fsm"
	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) Instances(rw web.ResponseWriter, req *web.Request) {
	result, err := c.Repository.GetListOfData(c.getInstanceKey(), models.Instance{})
	util.WriteJsonOrError(rw, result, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) ServicesInstances(rw web.ResponseWriter, req *web.Request) {
	instances, err := c.getFilteredInstances(models.InstanceTypeService, "")
	util.WriteJsonOrError(rw, instances, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) ServiceInstances(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.Repository.GetData(c.buildServiceKey(serviceId), models.Service{}); err != nil {
		handleGetDataError(rw, err)
		return
	}

	instances, err := c.getFilteredInstances(models.InstanceTypeService, serviceId)
	util.WriteJsonOrError(rw, instances, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) ApplicationsInstances(rw web.ResponseWriter, req *web.Request) {
	instances, err := c.getFilteredInstances(models.InstanceTypeApplication, "")
	util.WriteJsonOrError(rw, instances, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) ApplicationInstances(rw web.ResponseWriter, req *web.Request) {
	appId := req.PathParams["applicationId"]

	if _, err := c.Repository.GetData(c.buildApplicationKey(appId), models.Application{}); err != nil {
		handleGetDataError(rw, err)
		return
	}

	instances, err := c.getFilteredInstances(models.InstanceTypeApplication, appId)
	util.WriteJsonOrError(rw, instances, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) getFilteredInstances(expectedInstanceType models.InstanceType, expectedClassId string) ([]models.Instance, error) {
	return data.GetFilteredInstances(expectedInstanceType, expectedClassId, c.organization, c.Repository)
}

func (c *Context) GetApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationId); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.GetInstance(rw, req)
}

func (c *Context) GetServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.getService(serviceId); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.GetInstance(rw, req)
}

func (c *Context) GetInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]

	result, err := c.Repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	util.WriteJsonOrError(rw, result, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) GetInstanceBindings(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	result := []models.Instance{}

	instance, err := c.Repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	for _, binding := range instance.(models.Instance).Bindings {
		boundInstance, err := c.Repository.GetData(c.buildInstanceKey(binding.Id), models.Instance{})
		if err != nil {
			handleGetDataError(rw, err)
			return
		}
		result = append(result, boundInstance.(models.Instance))
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	c.addInstance(rw, req, req.PathParams["applicationId"], models.InstanceTypeApplication)
}

func (c *Context) AddServiceInstance(rw web.ResponseWriter, req *web.Request) {
	if req.URL.Query().Get("isServiceBroker") == "true" {
		c.addInstance(rw, req, req.PathParams["serviceId"], models.InstanceTypeServiceBroker)
	} else {
		c.addInstance(rw, req, req.PathParams["serviceId"], models.InstanceTypeService)
	}
}

func (c *Context) addInstance(rw web.ResponseWriter, req *web.Request, classId string, instanceType models.InstanceType) {
	reqInstance := &models.Instance{}

	if instanceType == models.InstanceTypeService {
		_, err := c.Repository.GetData(c.buildServiceKey(classId), models.Service{})
		if err != nil {
			util.Respond404(rw, errors.New("service with id: "+classId+" does not exists!"))
			return
		}
	} else if instanceType == models.InstanceTypeApplication {
		_, err := c.Repository.GetData(c.buildApplicationKey(classId), models.Application{})
		if err != nil {
			util.Respond404(rw, errors.New("application with id: "+classId+" does not exists!"))
			return
		}
	}

	err := util.ReadJson(req, reqInstance)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqInstance)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	if instanceType != models.InstanceTypeApplication && models.GetValueFromMetadata(reqInstance.Metadata, models.OFFERING_PLAN_ID) == "" {
		util.Respond400(rw, errors.New(fmt.Sprintf("key %s not found!", models.OFFERING_PLAN_ID)))
		return
	}

	err = data.CheckIfMatchingRegexp(reqInstance.Name, data.RegexpDnsLabelLowercase)
	if err != nil {
		util.Respond400(rw, errors.New("Field: Name has incorrect value: "+reqInstance.Name))
		return
	}

	for _, binding := range reqInstance.Bindings {
		_, err = c.Repository.GetData(c.buildInstanceKey(binding.Id), models.Instance{})
		if err != nil {
			util.Respond400(rw, errors.New(
				fmt.Sprintf("Field: binding ID has incorrect value: %s!", binding.Id)))
			return
		}
		for k, _ := range binding.Data {
			if err = data.CheckIfMatchingRegexp(k, data.RegexpDnsLabel); err != nil {
				util.Respond400(rw, errors.New("Field: data has incorrect value: "+k))
				return
			}
		}
	}

	exists, err := c.Repository.IsExistByName(reqInstance.Name, models.Instance{}, c.getInstanceKey())
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	if exists {
		util.Respond409(rw, errors.New("instance with name: "+reqInstance.Name+" already exists!"))
		return
	}

	reqInstance.ClassId = classId
	reqInstance.Type = instanceType
	reqInstance.State = models.InstanceStateRequested

	err = c.Repository.StoreData(c.mapper.ToKeyValue(c.getInstanceKey(), reqInstance, true))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	instance, err := c.Repository.GetData(c.buildInstanceKey(reqInstance.Id), models.Instance{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, instance, http.StatusCreated)
}

func (c *Context) PatchServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceID := req.PathParams["serviceId"]

	if _, err := c.getService(serviceID); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.PatchInstance(rw, req)
}

func (c *Context) PatchApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationID := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationID); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.PatchInstance(rw, req)
}

func (c *Context) PatchInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	instanceInt, err := c.Repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	instance, ok := instanceInt.(models.Instance)
	if !ok {
		util.Respond500(rw, errors.New("Instance retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.allowStateChange(patches, c.getInstancesFSM(instance.State))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildInstanceKey(instanceId), models.Instance{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.Repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	instanceInt, err = c.Repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, instanceInt, http.StatusOK)
}

func (c *Context) DeleteServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.getService(serviceId); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.DeleteInstance(rw, req)
}

func (c *Context) DeleteApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationID := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationID); err != nil {
		handleGetDataError(rw, err)
		return
	}

	c.DeleteInstance(rw, req)
}

func (c *Context) DeleteInstance(rw web.ResponseWriter, req *web.Request) {
	instanceID := req.PathParams["instanceId"]
	err := c.Repository.DeleteData(c.buildInstanceKey(instanceID))
	util.WriteJsonOrError(rw, "", getHttpStatusOrStatusError(http.StatusNoContent, err), err)
}

func (c *Context) getInstanceKey() string {
	return data.GetEntityKey(c.organization, data.Instances)
}

func (c *Context) buildInstanceKey(instanceId string) string {
	return c.mapper.ToKey(c.getInstanceKey(), instanceId)
}

func (c *Context) getInstancesFSM(initialState models.InstanceState) *fsm.FSM {
	return fsm.NewFSM(string(initialState),
		fsm.Events{
			{Name: "DEPLOYING", Src: []string{"REQUESTED"}, Dst: "DEPLOYING"},
			{Name: "FAILURE", Src: []string{"DEPLOYING", "STARTING", "RUNNING", "STOPPING", "DESTROYING"}, Dst: "FAILURE"},
			{Name: "STOPPED", Src: []string{"DEPLOYING", "STOPPING", "UNAVAILABLE"}, Dst: "STOPPED"},
			{Name: "START_REQ", Src: []string{"STOPPED"}, Dst: "START_REQ"},
			{Name: "STARTING", Src: []string{"START_REQ", "STOPPED"}, Dst: "STARTING"},
			{Name: "RUNNING", Src: []string{"STARTING"}, Dst: "RUNNING"},
			{Name: "STOP_REQ", Src: []string{"RUNNING", "STARTING"}, Dst: "STOP_REQ"},
			{Name: "STOPPING", Src: []string{"STOP_REQ"}, Dst: "STOPPING"},
			{Name: "DESTROY_REQ", Src: []string{"STOPPED", "FAILURE", "UNAVAILABLE"}, Dst: "DESTROY_REQ"},
			{Name: "DESTROYING", Src: []string{"DESTROY_REQ"}, Dst: "DESTROYING"},
			{Name: "UNAVAILABLE", Src: []string{"STOPPED"}, Dst: "UNAVAILABLE"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
