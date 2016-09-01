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
	"github.com/looplab/fsm"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-go-common/util"
	"strings"
)

var logger = logger_wrapper.InitLogger("api")

type Context struct {
	mapper     data.DataMapper
	repository data.RepositoryConnector
	organization string
}

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	util.WriteJson(rw, "I'm OK", http.StatusOK)
}

func (c *Context) Error(rw web.ResponseWriter, r *web.Request, err interface{}) {
	logger.Error("Respond500: reason: error ", err)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (c *Context) enterState(e *fsm.Event) {
	logger.Debugf("State changed from %s to %s", e.Src, e.Dst)
}

func (c *Context) allowStateChange(patches []models.Patch, stateMachine *fsm.FSM) error {
	for _, patch := range patches {
		if strings.EqualFold(patch.Field, "state") {
			value := c.removeQuotes(string(patch.Value))
			return stateMachine.Event(value)
		}
	}
	return nil
}

func (c *Context) removeQuotes(value string) string {
	return value[1 : len(value)-1]
}

func handleGetDataError(rw web.ResponseWriter, err error) {
	if err != nil {
		if strings.Contains(err.Error(), keyNotFoundMessage) {
			util.Respond404(rw, err)
		} else {
			util.Respond500(rw, err)
		}
	}
}
