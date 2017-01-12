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
	"strings"

	"github.com/gocraft/web"
	"github.com/looplab/fsm"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	commonHttp "github.com/trustedanalytics/tap-go-common/http"
)

func (c *Context) getServices() ([]models.Service, error) {
	result := []models.Service{}
	entities, err := c.repository.GetListOfData(c.getServiceKey(), models.Service{})
	if err != nil {
		err = fmt.Errorf("services retrieval failed: %v", err)
		logger.Warning(err)
		return []models.Service{}, err
	}

	for _, entity := range entities {
		service, ok := entity.(models.Service)
		if !ok {
			err = fmt.Errorf("type assertion for service failed: object from database: %v", entity)
			logger.Error(err)
			return []models.Service{}, err
		}
		result = append(result, service)
	}

	return result, nil
}

func (c *Context) Services(rw web.ResponseWriter, req *web.Request) {
	result, err := c.getServices()
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) getService(id string) (models.Service, error) {
	entity, err := c.repository.GetData(c.buildServiceKey(id), models.Service{})
	if err != nil {
		err = fmt.Errorf("service %q retrieval failed: %v", id, err)
		logger.Warning(err)
		return models.Service{}, err
	}

	service, ok := entity.(models.Service)
	if !ok {
		err = fmt.Errorf("type assertion for service %q failed: object from database: %v", id, entity)
		logger.Error(err)
		return models.Service{}, err
	}

	return service, nil
}

func (c *Context) GetService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	service, err := c.getService(serviceId)
	commonHttp.WriteJsonOrError(rw, service, http.StatusOK, err)
}

func (c *Context) AddService(rw web.ResponseWriter, req *web.Request) {
	reqService := &models.Service{}

	err := commonHttp.ReadJson(req, reqService)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqService)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	err = data.CheckIfMatchingRegexp(reqService.Name, data.RegexpDnsLabelLowercase)
	if err != nil {
		commonHttp.Respond400(rw, errors.New("Field: Name has incorrect value: "+reqService.Name))
		return
	}

	exists, err := c.repository.IsExistByName(reqService.Name, models.Service{}, c.getServiceKey())
	if err != nil {
		commonHttp.Respond500(rw, err)
		return
	}
	if exists {
		commonHttp.Respond409(rw, errors.New("service with name: "+reqService.Name+" already exists!"))
		return
	}

	if reqService.Id, err = c.reserveID(c.getServiceKey()); err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	reqService.State = models.ServiceStateDeploying
	serviceKeyStore := c.mapper.ToKeyValue(c.getServiceKey(), reqService, true)
	err = c.repository.CreateData(serviceKeyStore)
	if err != nil {
		commonHttp.Respond500(rw, err)
		return
	}

	service, err := c.repository.GetData(c.buildServiceKey(reqService.Id), models.Service{})
	commonHttp.WriteJsonOrError(rw, service, http.StatusCreated, err)
}

func (c *Context) PatchService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	serviceInt, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	service, ok := serviceInt.(models.Service)
	if !ok {
		commonHttp.HandleError(rw, errors.New("Service retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = commonHttp.ReadJson(req, &patches)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	newStateAsString, err := c.getStateChange(patches)
	newState := models.ServiceState(newStateAsString)
	if err != nil && newState == models.ServiceStateOffline && service.State == models.ServiceStateReady {
		if status, err := c.assureOfferingIsNotUsed(serviceId); err != nil {
			err := fmt.Errorf("cannot change offering state from %q to %q: %v", service.State, newState, err)
			logger.Error(err.Error())
			commonHttp.GenericRespond(status, rw, err)
			return
		}
	}

	fsmFunc := func() *fsm.FSM {
		return c.getServicesFSM(service.State)
	}
	if err = c.handleFsm(rw, req, patches, fsmFunc); err != nil {
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildServiceKey(serviceId), models.Service{}, patches)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	serviceInt, err = c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	commonHttp.WriteJsonOrError(rw, serviceInt, http.StatusOK, err)
}

func (c *Context) DeleteService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if status, err := c.assureOfferingIsNotUsed(serviceId); err != nil {
		err := fmt.Errorf("cannot remove offering %q: %v", serviceId, err)
		commonHttp.GenericRespond(status, rw, err)
		return
	}

	err := c.repository.DeleteData(c.buildServiceKey(serviceId))
	commonHttp.WriteJsonOrError(rw, serviceId, http.StatusNoContent, err)
}

func (c *Context) assureOfferingIsNotUsed(serviceID string) (int, error) {
	instances, err := c.getInstances()
	if err != nil {
		return getHttpStatusOrStatusError(http.StatusInternalServerError, err), err
	}

	if instanceExists, message := serviceInstanceExists(serviceID, instances); instanceExists {
		err := fmt.Errorf("offering instance exists: %v", message)
		return http.StatusForbidden, err
	}

	services, err := c.getServices()
	if err != nil {
		return getHttpStatusOrStatusError(http.StatusInternalServerError, err), err
	}

	if isDependence, message := isOfferingDependence(serviceID, services); isDependence {
		err := fmt.Errorf("dependent offerings detected: %v", message)
		return http.StatusForbidden, err
	}

	return http.StatusOK, nil
}

func serviceInstanceExists(serviceID string, instances []models.Instance) (bool, error) {
	instanceNames := []string{}
	for _, instance := range instances {
		if instance.ClassId == serviceID {
			instanceNames = append(instanceNames, instance.Name)
		}
	}

	if len(instanceNames) > 0 {
		joinedInstances := strings.Join(instanceNames, ", ")
		return true, fmt.Errorf("service %q has existing instances: %s", serviceID, joinedInstances)
	}

	return false, nil
}

func isOfferingDependence(serviceID string, services []models.Service) (bool, error) {
	dependentOfferingNames := []string{}
	for _, service := range services {
		if service.Id == serviceID {
			continue
		}
		for _, plan := range service.Plans {
			for _, dependency := range plan.Dependencies {
				if dependency.ServiceId == serviceID {
					dependentOfferingNames = append(dependentOfferingNames, service.Name)
				}
			}
		}
	}

	if len(dependentOfferingNames) > 0 {
		joinedServices := strings.Join(dependentOfferingNames, ", ")
		return true, fmt.Errorf("service %q is a dependence for following offerings: %s", serviceID, joinedServices)
	}

	return false, nil
}

func (c *Context) getServiceKey() string {
	return data.GetEntityKey(c.organization, data.Services)
}

func (c *Context) buildServiceKey(serviceId string) string {
	return c.mapper.ToKey(c.getServiceKey(), serviceId)
}

func (c *Context) getServicesFSM(initialState models.ServiceState) *fsm.FSM {
	return fsm.NewFSM(string(initialState),
		fsm.Events{
			{Name: "READY", Src: []string{"DEPLOYING"}, Dst: "READY"},
			{Name: "OFFLINE", Src: []string{"DEPLOYING", "READY"}, Dst: "OFFLINE"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
