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
	"log"
	"net/http"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tapng-catalog/data"
	"github.com/trustedanalytics/tapng-catalog/webutils"
)

type Context struct {
	mapper     data.DataMapper
	repository data.RepositoryConnector
}

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "I'm OK", http.StatusOK)
}

func (c *Context) Error(rw web.ResponseWriter, r *web.Request, err interface{}) {
	log.Println("Respond500: reason: error ", err)
	rw.WriteHeader(http.StatusInternalServerError)
}
