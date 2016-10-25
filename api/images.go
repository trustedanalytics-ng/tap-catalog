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

	"errors"
	"github.com/looplab/fsm"
	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

func (c *Context) Images(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(c.getImagesKey(), models.Image{})
	if err != nil {
		util.Respond500(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]

	result, err := c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddImage(rw web.ResponseWriter, req *web.Request) {
	reqImage := &models.Image{}

	err := util.ReadJson(req, reqImage)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	reqImage.State = models.ImageStatePending
	imageKeyStore := c.mapper.ToKeyValue(c.getImagesKey(), reqImage, true)

	err = c.repository.StoreData(imageKeyStore)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	image, err := c.repository.GetData(c.buildImagesKey(reqImage.Id), models.Image{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, image, http.StatusCreated)
}

func (c *Context) PatchImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	imageInt, err := c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}

	image, ok := imageInt.(models.Image)
	if !ok {
		util.Respond500(rw, errors.New("Image retrieved is in wrong format"))
		return
	}

	patches := []models.Patch{}
	err = util.ReadJson(req, &patches)
	if err != nil {
		util.Respond400(rw, err)
		return
	}

	err = c.allowStateChange(patches, c.getImagesFSM(image.State))
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildImagesKey(imageId), models.Image{}, patches)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	err = c.repository.ApplyPatchedValues(patchedValues)
	if err != nil {
		util.Respond500(rw, err)
		return
	}

	imageInt, err = c.repository.GetData(c.buildImagesKey(imageId), models.Image{})
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, imageInt, http.StatusOK)
}

func (c *Context) DeleteImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	err := c.repository.DeleteData(c.buildImagesKey(imageId))
	if err != nil {
		handleGetDataError(rw, err)
		return
	}
	util.WriteJson(rw, "", http.StatusNoContent)
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
			{Name: "BUILDING", Src: []string{"PENDING"}, Dst: "BUILDING"},
			{Name: "ERROR", Src: []string{"BUILDING"}, Dst: "ERROR"},
			{Name: "READY", Src: []string{"BUILDING"}, Dst: "READY"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				c.enterState(e)
			},
		},
	)
}
