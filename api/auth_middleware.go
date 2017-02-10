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
	"os"

	"github.com/gocraft/web"

	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func (c *Context) BasicAuthorizeMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	username, password, is_ok := req.BasicAuth()
	if !is_ok || username != os.Getenv("CATALOG_USER") || password != os.Getenv("CATALOG_PASS") {
		commonHttp.RespondUnauthorized(rw)
		return
	}
	c.mapper.Username = username
	next(rw, req)
}
