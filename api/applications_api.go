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
	"github.com/tapng/tap-catalog/webutils"
	"github.com/gocraft/web"
)

func (c *Context) Applications(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "List Applications", http.StatusOK)
}

func (c *Context) GetApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Single Application", http.StatusOK)
}

func (c *Context) AddApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Create Application", http.StatusCreated)
}

func (c *Context) UpdateApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Update Application", http.StatusOK)
}

func (c *Context) DeleteApplication(rw web.ResponseWriter, req *web.Request) {
	webutils.WriteJson(rw, "Delete Application", http.StatusNoContent)
}
