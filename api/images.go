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
	"github.com/trustedanalytics-ng/tap-go-common/util"
)

func (c *Context) Images(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(c.getImagesKey(), models.Image{})
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) GetImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]

	result, err := c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	commonHttp.WriteJsonOrError(rw, result, http.StatusOK, err)
}

func (c *Context) AddImage(rw web.ResponseWriter, req *web.Request) {
	reqImage := &models.Image{}

	err := commonHttp.ReadJson(req, reqImage)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	reqImage.State = models.ImageStateRequested
	imageKeyStore := c.mapper.ToKeyValue(c.getImagesKey(), reqImage, true)

	err = c.repository.CreateData(imageKeyStore)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	image, err := c.repository.GetData(c.buildImagesKey(reqImage.Id), models.Image{})
	commonHttp.WriteJsonOrError(rw, image, http.StatusCreated, err)
}

func (c *Context) PatchImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	imageInt, err := c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	image, ok := imageInt.(models.Image)
	if !ok {
		commonHttp.HandleError(rw, errors.New("Image retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = commonHttp.ReadJson(req, &patches)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	fsmFunc := func() *fsm.FSM {
		return c.getImagesFSM(image.State)
	}
	if err = c.handleFsm(rw, req, patches, fsmFunc); err != nil {
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildImagesKey(imageId), models.Image{}, patches)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		commonHttp.HandleError(rw, err)
		return
	}

	imageInt, err = c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	commonHttp.WriteJsonOrError(rw, imageInt, http.StatusOK, err)
}

func (c *Context) DeleteImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	err := c.repository.DeleteData(c.buildImagesKey(imageId))
	commonHttp.WriteJsonOrError(rw, "", http.StatusNoContent, err)
}

func (c *Context) GetImageCheckRefs(rw web.ResponseWriter, req *web.Request) {
	imageID := req.PathParams["imageId"]
	var response models.ImageRefsResponse
	var err error

	response.ApplicationReferences, err = c.applicationImageRefs(imageID)
	if err != nil {
		logger.Errorf("applicationImageRefs returned error: %v", err)
		commonHttp.GenericRespond(http.StatusInternalServerError, rw, err)
		return
	}

	response.ServiceReferences, err = c.servicesImageRefs(imageID)
	if err != nil {
		logger.Errorf("servicesImageRefs returned error: %v", err)
		commonHttp.GenericRespond(http.StatusInternalServerError, rw, err)
		return
	}

	if len(response.ApplicationReferences) > 0 || len(response.ServiceReferences) > 0 {
		response.IsAnyRefExist = true
	}
	commonHttp.WriteJson(rw, response, http.StatusOK)
}

func (c *Context) applicationImageRefs(imageID string) ([]models.Application, error) {
	var appsWhichUsesThisImage []models.Application
	applications, err := c.repository.GetListOfData(c.getApplicationKey(), models.Application{})
	if err != nil {
		return appsWhichUsesThisImage, err
	}
	for _, appInterface := range applications {
		app, ok := appInterface.(models.Application)
		if !ok {
			err = errors.New("type assertion error for application!")
			return appsWhichUsesThisImage, err
		}

		if app.ImageId == imageID {
			appsWhichUsesThisImage = append(appsWhichUsesThisImage, app)
		}
	}
	return appsWhichUsesThisImage, nil
}

func (c *Context) servicesImageRefs(imageID string) ([]models.Service, error) {
	var servicesWhichUsesThisImage []models.Service
	services, err := c.repository.GetListOfData(c.getServiceKey(), models.Service{})
	if err != nil {
		return servicesWhichUsesThisImage, err
	}
	for _, servInterface := range services {
		srv, ok := servInterface.(models.Service)
		if !ok {
			return servicesWhichUsesThisImage, errors.New("type assertion error for service!")
		}

		for _, meta := range srv.Metadata {
			if meta.Id == models.APPLICATION_IMAGE_ADDRESS {
				_, imageIDInService, _, err := util.ParseImageAddress(meta.Value)
				if err != nil {
					return servicesWhichUsesThisImage, fmt.Errorf("%s malformed in service: %v , err: %v", models.APPLICATION_IMAGE_ADDRESS, srv.Name, err)
				}

				if imageIDInService == models.GenerateImageId(imageID) {
					servicesWhichUsesThisImage = append(servicesWhichUsesThisImage, srv)
				}
			}
		}
	}
	return servicesWhichUsesThisImage, nil
}

func (c *Context) MonitorImagesStates(rw web.ResponseWriter, req *web.Request) {
	c.monitorSpecificState(rw, req, c.buildImagesKey(""))
}

func (c *Context) MonitorSpecificImageState(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	c.monitorSpecificState(rw, req, c.buildImagesKey(imageId))
}

func (c *Context) getImagesKey() string {
	org := c.mapper.ToKey("", c.organization)
	return c.mapper.ToKey(org, data.Images)
}

func (c *Context) buildImagesKey(imageId string) string {
	return c.mapper.ToKey(c.getImagesKey(), imageId)
}

func (c *Context) getImagesFSM(initialState models.ImageState) *fsm.FSM {
	return fsm.NewFSM(string(initialState),
		fsm.Events{
			{Name: "PENDING", Src: []string{"REQUESTED"}, Dst: "PENDING"},
			{Name: "BUILDING", Src: []string{"PENDING"}, Dst: "BUILDING"},
			{Name: "ERROR", Src: []string{"BUILDING"}, Dst: "ERROR"},
			{Name: "READY", Src: []string{"BUILDING"}, Dst: "READY"},
			{Name: "REMOVING", Src: []string{"READY"}, Dst: "REMOVING"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
