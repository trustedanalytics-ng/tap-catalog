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
	"fmt"
	"net/http"
	"strings"

	"github.com/gocraft/web"
	"github.com/looplab/fsm"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	commonHttp "github.com/trustedanalytics/tap-go-common/http"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

const (
	maxUUIDGenerationTrials = 10
)

var logger, _ = commonLogger.InitLogger("api")

type Context struct {
	mapper       data.DataMapper
	repository   data.RepositoryApi
	organization string
}

func NewContext(r data.RepositoryApi, org string) (Context, error) {
	ctx := Context{
		repository:   r,
		organization: org,
	}
	return ctx, ctx.initDB(org)
}

func (c *Context) initDB(org string) error {
	err := c.repository.CreateDirs(org)
	if err != nil {
		return fmt.Errorf("cannot create directories in ETCD for organization %s: %v", org, err)
	}
	return nil
}

func (c *Context) enterState(e *fsm.Event) {
	logger.Debugf("State changed from %s to %s", e.Src, e.Dst)
}

func (c *Context) getStateChange(patches []models.Patch) (string, error) {
	for _, patch := range patches {
		if err := models.ValidatePatchStructure(patch); err != nil {
			return "", err
		} else if strings.EqualFold(*patch.Field, "state") {
			value := c.removeQuotes(string(*patch.Value))
			return value, nil
		}
	}
	return "", nil
}

func (c *Context) allowStateChange(newState string, stateMachine *fsm.FSM) error {
	if newState != "" {
		return stateMachine.Event(newState)
	}
	return nil
}

type FsmFunc func() *fsm.FSM

func (c *Context) handleFsm(rw web.ResponseWriter, req *web.Request, patches []models.Patch, fsmFunc FsmFunc) error {
	newState, err := c.getStateChange(patches)
	if err != nil {
		commonHttp.Respond400(rw, err)
		return err
	}
	fsm := fsmFunc()
	err = c.allowStateChange(newState, fsm)
	if err != nil {
		if fsm.Current() == newState {
			commonHttp.Respond409(rw, err)
		} else {
			commonHttp.Respond400(rw, err)
		}
		return err
	}
	return nil
}

func (c *Context) removeQuotes(value string) string {
	return value[1 : len(value)-1]
}

func (c *Context) reserveID(path string) (string, error) {
	id := ""
	var err error
	idCreated := false
	for i := 0; i < maxUUIDGenerationTrials; i++ {
		id, err = data.GenerateID()
		if err != nil {
			return "", fmt.Errorf("generation ID failed: %v", err)
		}
		dir := fmt.Sprintf("%s/%s", path, id)
		if err := c.repository.CreateDir(dir); err == nil {
			idCreated = true
			break
		}
	}

	if !idCreated {
		return "", fmt.Errorf("cannot create entity with generated ID after %d trials (notice: this is incredibly unlikely)", maxUUIDGenerationTrials)
	}
	return id, nil
}

func getHttpStatusOrStatusError(status int, err error) int {
	if err != nil {
		if commonHttp.IsNotFoundError(err) {
			return http.StatusNotFound
		} else if commonHttp.IsAlreadyExistsError(err) {
			return http.StatusConflict
		}
		return http.StatusInternalServerError
	}
	return status
}
