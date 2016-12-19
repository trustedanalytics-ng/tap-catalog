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

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
)

var stableStates = []models.InstanceState{
	models.InstanceStateStopped,
	models.InstanceStateRunning,
	models.InstanceStateFailure,
}

func (c *Context) CheckStateStability(rw web.ResponseWriter, req *web.Request) {
	instances, err := c.getInstances()
	if err != nil {
		util.WriteJson(rw, err.Error(), getHttpStatusOrStatusError(http.StatusOK, err))
		return
	}

	err = assureInstanceStatesAreStable(instances)
	if err != nil {
		util.WriteJson(rw, models.StateStability{Stable: false, Message: err.Error()}, http.StatusOK)
		return
	}

	util.WriteJson(rw, models.StateStability{Stable: true}, http.StatusOK)
}

func assureInstanceStatesAreStable(instances []models.Instance) error {
	for _, instance := range instances {
		if !isInstanceInStableState(instance) {
			return fmt.Errorf("instance %q state %q is not stable", instance.Id, instance.State.String())
		}
	}

	return nil
}

func isInstanceInStableState(instance models.Instance) bool {
	for _, stateOK := range stableStates {
		if instance.State == stateOK {
			return true
		}
	}
	return false
}
