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
	"github.com/trustedanalytics/tap-go-common/util"
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
	util.WriteJsonOrError(rw, result, getHttpStatusOrStatusError(http.StatusOK, err), err)
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
	util.WriteJsonOrError(rw, service, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) AddService(rw web.ResponseWriter, req *web.Request) {
	reqService := &models.Service{}

	err := util.ReadJson(req, reqService)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfIdFieldIsEmpty(reqService)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = data.CheckIfMatchingRegexp(reqService.Name, data.RegexpDnsLabelLowercase)
	if err != nil {
		util.Respond400(rw, errors.New("Field: Name has incorrect value: "+reqService.Name))
		return
	}

	exists, err := c.repository.IsExistByName(reqService.Name, models.Service{}, c.getServiceKey())
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	if exists {
		util.Respond409(rw, errors.New("service with name: "+reqService.Name+" already exists!"))
		return
	}

	if reqService.Id, err = c.reserveID(c.getServiceKey()); err != nil {
		util.Respond500(rw, err)
		return
	}

	reqService.State = models.ServiceStateDeploying
	serviceKeyStore := c.mapper.ToKeyValue(c.getServiceKey(), reqService, true)
	err = c.repository.CreateData(serviceKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	service, err := c.repository.GetData(c.buildServiceKey(reqService.Id), models.Service{})
	util.WriteJsonOrError(rw, service, getHttpStatusOrStatusError(http.StatusCreated, err), err)
}

func (c *Context) PatchService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]
	serviceInt, err := c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	service, ok := serviceInt.(models.Service)
	if !ok {
		util.Respond500(rw, errors.New("Service retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.allowStateChange(patches, c.getServicesFSM(service.State))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildServiceKey(serviceId), models.Service{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	serviceInt, err = c.repository.GetData(c.buildServiceKey(serviceId), models.Service{})
	util.WriteJsonOrError(rw, serviceInt, getHttpStatusOrStatusError(http.StatusOK, err), err)
}

func (c *Context) DeleteService(rw web.ResponseWriter, req *web.Request) {
	serviceId := req.PathParams["serviceId"]

	if status, err := c.assureOfferingCanBeDeleted(serviceId); err != nil {
		err := fmt.Errorf("cannot remove offering %q: %v", serviceId, err)
		util.GenericRespond(status, rw, err)
		return
	}

	err := c.repository.DeleteData(c.buildServiceKey(serviceId))
	util.WriteJsonOrError(rw, serviceId, getHttpStatusOrStatusError(http.StatusNoContent, err), err)
}

func (c *Context) assureOfferingCanBeDeleted(serviceID string) (int, error) {
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
			{Name: "OFFLINE", Src: []string{"DEPLOYING"}, Dst: "OFFLINE"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
