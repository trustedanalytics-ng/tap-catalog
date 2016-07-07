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

	"github.com/trustedanalytics/tapng-catalog/data"
	"github.com/trustedanalytics/tapng-catalog/models"
	"github.com/trustedanalytics/tapng-catalog/webutils"
)

func (c *Context) Images(rw web.ResponseWriter, req *web.Request) {
	result, err := c.repository.GetListOfData(data.Images, &models.Image{})
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) GetImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]

	result, err := c.repository.GetData(c.buildImagesKey(imageId), &models.Image{})
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, result, http.StatusOK)
}

func (c *Context) AddImage(rw web.ResponseWriter, req *web.Request) {
	reqImage := models.Image{}

	err := webutils.ReadJson(req, &reqImage)
	if err != nil {
		webutils.Respond400(rw, err)
		return
	}

	imageId, err := uuid.NewV4()
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	reqImage.Id = imageId.String()
	imageKeyStore := c.mapper.ToKeyValue(data.Images, reqImage, true)
	err = c.repository.StoreData(imageKeyStore)
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}

	image, err := c.repository.GetData(c.buildImagesKey(imageId.String()), &models.Image{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, image, http.StatusCreated)
}

func (c *Context) PatchImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	image, err := c.repository.GetData(c.buildImagesKey(imageId), &models.Image{})
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

	patchedValues, err := c.mapper.ToKeyValueByPatches(c.buildImagesKey(imageId), models.Image{}, patches)
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

	image, err = c.repository.GetData(c.buildImagesKey(imageId), &models.Image{})
	if err != nil {
		logger.Error(err)
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, image, http.StatusOK)
}

func (c *Context) DeleteImage(rw web.ResponseWriter, req *web.Request) {
	imageId := req.PathParams["imageId"]
	err := c.repository.DeleteData(c.buildImagesKey(imageId))
	if err != nil {
		webutils.Respond500(rw, err)
		return
	}
	webutils.WriteJson(rw, "", http.StatusNoContent)
}

func (c *Context) buildImagesKey(imageId string) string {
	return c.mapper.ToKey(data.Images, imageId)
}
