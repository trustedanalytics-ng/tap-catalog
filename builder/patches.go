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
package builder

import (
	"encoding/json"

	"github.com/trustedanalytics/tap-catalog/models"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

var logger, _ = commonLogger.InitLogger("builder")

// currentState is optional, if passed then CAS will be used
// message is optional, if pass then LAST_STATE_CHANGE_REASON key will be added/updated in Instance Metadata
func MakePatchesForInstanceStateAndLastStateMetadata(message string, currentState, stateToSet models.InstanceState) ([]models.Patch, error) {
	patches := []models.Patch{}

	statePatch, err := MakePatch("State", stateToSet, models.OperationUpdate)
	if err != nil {
		return patches, err
	}

	if currentState != "" {
		statePatch.PrevValue, err = json.Marshal(currentState)
		if err != nil {
			logger.Error("previousState marshal error:", err)
			return patches, err
		}
	}
	patches = append(patches, statePatch)

	if message != "" {
		lastStateChangeReasonMetadata := models.Metadata{
			Id: models.LAST_STATE_CHANGE_REASON, Value: message,
		}
		metadataPatch, err := MakePatch("Metadata", lastStateChangeReasonMetadata, models.OperationAdd)
		if err != nil {
			return patches, err
		}
		patches = append(patches, metadataPatch)
	}
	return patches, nil
}

func MakePatch(field string, valueToUpdate interface{}, operation models.PatchOperation) (models.Patch, error) {
	instanceBindingByte, err := json.Marshal(valueToUpdate)
	if err != nil {
		logger.Errorf("marshal filed %s error: %v", field, err)
		return models.Patch{}, err
	}

	patch := models.Patch{
		Operation: operation,
		Field:     field,
		Value:     instanceBindingByte,
	}
	return patch, nil
}
