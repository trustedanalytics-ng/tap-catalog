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

	"github.com/trustedanalytics-ng/tap-catalog/data"
	"github.com/trustedanalytics-ng/tap-catalog/models"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func (c *Context) Instances(rw web.ResponseWriter, req *web.Request) {
	result, err := c.getInstances()
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) getInstances() ([]models.Instance, error) {
	result := []models.Instance{}
	entities, err := c.repository.GetListOfData(c.getInstanceKey(), models.Instance{})
	if err != nil {
		err = fmt.Errorf("instances retrieval failed: %v", err)
		logger.Warning(err)
		return []models.Instance{}, err
	}

	for _, entity := range entities {
		instance, ok := entity.(models.Instance)
		if !ok {
			err = fmt.Errorf("type assertion for instance failed: object from database: %v", entity)
			logger.Error(err)
			return []models.Instance{}, err
		}
		result = append(result, instance)
	}

	return result, nil
}

func (c *Context) ServicesInstances(rw web.ResponseWriter, req *web.Request) {
	instances, err := c.getFilteredInstances(models.InstanceTypeService, "")
	commonHttp.WriteJsonOrError(rw, instances, http.StatusOK, err)
}

func (c *Context) ServiceInstances(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{}); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	instances, err := c.getFilteredInstances(models.InstanceTypeService, serviceId)
	commonHttp.WriteJsonOrError(rw, instances, http.StatusOK, err)
}

func (c *Context) ApplicationsInstances(rw web.ResponseWriter, req *web.Request) {
	instances, err := c.getFilteredInstances(models.InstanceTypeApplication, "")
	commonHttp.WriteJsonOrError(rw, instances, http.StatusOK, err)
}

func (c *Context) ApplicationInstances(rw web.ResponseWriter, req *web.Request) {
	appId := req.PathParams["applicationId"]

	if _, err := c.repository.GetData(c.buildApplicationKey(appId), models.Application{}); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	instances, err := c.getFilteredInstances(models.InstanceTypeApplication, appId)
	commonHttp.WriteJsonOrError(rw, instances, http.StatusOK, err)
}

func (c *Context) getFilteredInstances(expectedInstanceType models.InstanceType, expectedClassId string) ([]models.Instance, error) {
	return data.GetFilteredInstances(expectedInstanceType, expectedClassId, c.organization, c.repository)
}

func (c *Context) GetApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationId := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationId); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.GetInstance(rw, req)
}

func (c *Context) GetServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.getService(serviceId); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.GetInstance(rw, req)
}

func (c *Context) GetInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]

	result, err := c.repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) GetInstanceBindings(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	result := []models.Instance{}

	instance, err := c.repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	for _, binding := range instance.(models.Instance).Bindings {
		boundInstance, err := c.repository.GetData(c.buildInstanceKey(binding.Id), models.Instance{})
		if err != nil {
			commonHttp.HandleError(rw, err)
			return
		}
		result = append(result, boundInstance.(models.Instance))
	}
	commonHttp.WriteJson(rw, result, http.StatusOK)
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
		_, err := c.repository.GetData(c.buildServiceKey(classId), models.Service{})
		if err != nil {
			commonHttp.Respond404(rw, errors.New("service with id: "+classId+" does not exists!"))
			return
		}
	} else if instanceType == models.InstanceTypeApplication {
		_, err := c.repository.GetData(c.buildApplicationKey(classId), models.Application{})
		if err != nil {
			commonHttp.Respond404(rw, errors.New("application with id: "+classId+" does not exists!"))
			return
		}
	}

	err := commonHttp.ReadJson(req, reqInstance)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	err = reqInstance.ValidateInstanceStructCreate(instanceType)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	for _, binding := range reqInstance.Bindings {
		_, err = c.repository.GetData(c.buildInstanceKey(binding.Id), models.Instance{})
		if err != nil {
			commonHttp.Respond400(rw, errors.New(
				fmt.Sprintf("Field: binding ID has incorrect value: %s!", binding.Id)))
			return
		}
	}

	exists, err := c.repository.IsExistByName(reqInstance.Name, models.Instance{}, c.getInstanceKey())
	if err != nil {
		commonHttp.Respond500(rw, err)
		return
	}
	if exists {
		commonHttp.Respond409(rw, errors.New("instance with name: "+reqInstance.Name+" already exists!"))
		return
	}

	if reqInstance.Id, err = c.reserveID(c.getInstanceKey()); err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	reqInstance.ClassId = classId
	reqInstance.Type = instanceType
	reqInstance.State = models.InstanceStateRequested

	err = c.repository.CreateData(c.mapper.ToKeyValue(c.getInstanceKey(), reqInstance, true))
	if err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	instance, err := c.repository.GetData(c.buildInstanceKey(reqInstance.Id), models.Instance{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}
	commonHttp.WriteJson(rw, instance, http.StatusCreated)
}

func (c *Context) PatchServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceID := req.PathParams["serviceId"]

	if _, err := c.getService(serviceID); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.PatchInstance(rw, req)
}

func (c *Context) PatchApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationID := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationID); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.PatchInstance(rw, req)
}

func (c *Context) PatchInstance(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	instanceInt, err := c.repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	instance, ok := instanceInt.(models.Instance)
	if !ok {
		commonHttp.Respond500(rw, errors.New("Instance retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = commonHttp.ReadJson(req, &patches)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	fsmFunc := func() *fsm.FSM {
		return c.getInstancesFSM(instance.State)
	}
	if err = c.handleFsm(rw, req, patches, fsmFunc); err != nil {
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildInstanceKey(instanceId), models.Instance{}, patches)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	instanceInt, err = c.repository.GetData(c.buildInstanceKey(instanceId), models.Instance{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}
	commonHttp.WriteJson(rw, instanceInt, http.StatusOK)
}

func (c *Context) DeleteServiceInstance(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if _, err := c.getService(serviceId); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.DeleteInstance(rw, req)
}

func (c *Context) DeleteApplicationInstance(rw web.ResponseWriter, req *web.Request) {
	applicationID := req.PathParams["applicationId"]

	if _, err := c.getApplication(applicationID); err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	c.DeleteInstance(rw, req)
}

func (c *Context) DeleteInstance(rw web.ResponseWriter, req *web.Request) {
	instanceID := req.PathParams["instanceId"]
	err := c.repository.DeleteData(c.buildInstanceKey(instanceID))
	commonHttp.WriteJsonOrError(rw, "", http.StatusNoContent, err)
}

func (c *Context) MonitorInstancesStates(rw web.ResponseWriter, req *web.Request) {
	c.monitorSpecificState(rw, req, c.buildInstanceKey(""))
}

func (c *Context) MonitorSpecificInstanceState(rw web.ResponseWriter, req *web.Request) {
	instanceId := req.PathParams["instanceId"]
	c.monitorSpecificState(rw, req, c.buildInstanceKey(instanceId))
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
			makeEventDesc(models.InstanceStateDeploying, models.InstanceStateRequested),
			makeEventDesc(models.InstanceStateFailure, models.InstanceStateDeploying, models.InstanceStateStarting,
				models.InstanceStateRunning, models.InstanceStateStopping, models.InstanceStateDestroying),
			makeEventDesc(models.InstanceStateStopped, models.InstanceStateDeploying, models.InstanceStateStopping,
				models.InstanceStateUnavailable),
			makeEventDesc(models.InstanceStateStartReq, models.InstanceStateStopped),
			makeEventDesc(models.InstanceStateStarting, models.InstanceStateStartReq, models.InstanceStateStopped,
				models.InstanceStateReconfiguration),
			makeEventDesc(models.InstanceStateRunning, models.InstanceStateStarting),
			makeEventDesc(models.InstanceStateReconfiguration, models.InstanceStateRunning, models.InstanceStateStopped),
			makeEventDesc(models.InstanceStateStopReq, models.InstanceStateRunning, models.InstanceStateStarting),
			makeEventDesc(models.InstanceStateStopping, models.InstanceStateStopReq, models.InstanceStateReconfiguration),
			makeEventDesc(models.InstanceStateDestroyReq, models.InstanceStateStopped, models.InstanceStateFailure),
			makeEventDesc(models.InstanceStateDestroying, models.InstanceStateDestroyReq),
			makeEventDesc(models.InstanceStateUnavailable, models.InstanceStateStopped, models.InstanceStateRunning),
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}

func makeEventDesc(destination models.InstanceState, sources ...models.InstanceState) fsm.EventDesc {
	sourceString := []string{}
	for _, source := range sources {
		sourceString = append(sourceString, source.String())
	}

	return fsm.EventDesc{
		Name: destination.String(),
		Src:  sourceString,
		Dst:  destination.String(),
	}
}
