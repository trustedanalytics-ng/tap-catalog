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
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/models"
)

func TestMakePatch(t *testing.T) {
	field := "metadata"
	operation := models.OperationAdd
	value := models.Metadata{Id: "id", Value: "value"}

	Convey("Test MakePatch", t, func() {
		Convey("Should return proper response", func() {
			byteValue, err := json.Marshal(value)
			So(err, ShouldBeNil)

			expectedPatch := models.Patch{
				Operation: operation,
				Field:     field,
				Value:     byteValue,
			}

			patch, err := MakePatch(field, value, operation)
			So(err, ShouldBeNil)
			So(patch, ShouldResemble, expectedPatch)
		})
	})
}

func TestMakePatchWithPreviousValue(t *testing.T) {
	field := "metadata"
	operation := models.OperationUpdate
	value := models.Metadata{Id: "id", Value: "value"}
	previousValue := models.Metadata{Id: "id", Value: "old"}

	Convey("Test MakePatchWithPreviousValue", t, func() {
		Convey("Should return proper response", func() {
			byteValue, err := json.Marshal(value)
			So(err, ShouldBeNil)

			bytePreviousValue, err := json.Marshal(previousValue)
			So(err, ShouldBeNil)

			expectedPatch := models.Patch{
				Operation: operation,
				Field:     field,
				Value:     byteValue,
				PrevValue: bytePreviousValue,
			}

			patch, err := MakePatchWithPreviousValue(field, value, previousValue, operation)
			So(err, ShouldBeNil)
			So(patch, ShouldResemble, expectedPatch)
		})
	})
}

func TestMakePatchesForInstanceStateAndLastStateMetadata(t *testing.T) {
	Convey("Test MakePatchesForInstanceStateAndLastStateMetadata", t, func() {
		Convey("Should return error if PrevValue not set", func() {
			_, err := MakePatchesForInstanceStateAndLastStateMetadata("", "", models.InstanceStateStarting)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "currentState and stateToSet cannot be empty!")
		})

		Convey("Should return one patch with set PrevValue", func() {
			newState := models.InstanceStateStarting
			byteNewValue, err := json.Marshal(newState)
			So(err, ShouldBeNil)

			oldState := models.InstanceStateStopped
			byteOldValue, err := json.Marshal(oldState)
			So(err, ShouldBeNil)

			patches, err := MakePatchesForInstanceStateAndLastStateMetadata("", oldState, newState)
			So(err, ShouldBeNil)
			So(len(patches), ShouldEqual, 1)
			So(patches[0].Field, ShouldEqual, "State")
			So(patches[0].Operation, ShouldEqual, models.OperationUpdate)
			So(patches[0].Value, ShouldResemble, json.RawMessage(byteNewValue))
			So(patches[0].PrevValue, ShouldResemble, json.RawMessage(byteOldValue))
		})

		Convey("Should return two patches if message set", func() {
			state := models.InstanceStateStarting
			byteStateValue, err := json.Marshal(state)
			So(err, ShouldBeNil)

			oldState := models.InstanceStateStopped
			byteOldValue, err := json.Marshal(oldState)
			So(err, ShouldBeNil)

			message := "test-message"
			byteMessageValue, err := json.Marshal(models.Metadata{
				Id: models.LAST_STATE_CHANGE_REASON, Value: message,
			})
			So(err, ShouldBeNil)

			patches, err := MakePatchesForInstanceStateAndLastStateMetadata(message, oldState, state)
			So(err, ShouldBeNil)
			So(len(patches), ShouldEqual, 2)
			So(patches[0].Field, ShouldEqual, "State")
			So(patches[0].Operation, ShouldEqual, models.OperationUpdate)
			So(patches[0].Value, ShouldResemble, json.RawMessage(byteStateValue))
			So(patches[0].PrevValue, ShouldResemble, json.RawMessage(byteOldValue))
			So(patches[1].Field, ShouldEqual, "Metadata")
			So(patches[1].Operation, ShouldEqual, models.OperationAdd)
			So(patches[1].Value, ShouldResemble, json.RawMessage(byteMessageValue))
			So(patches[1].PrevValue, ShouldResemble, json.RawMessage(nil))
		})
	})
}

func TestMakePatchesForOfferingStateUpdate(t *testing.T) {
	Convey("Test MakePatchesForOfferingStateUpdate", t, func() {
		Convey("Should return error if PrevValue not set", func() {
			_, err := MakePatchesForOfferingStateUpdate("", "")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "currentState and stateToSet cannot be empty!")
		})

		Convey("Should return one patch with set PrevValue", func() {
			newState := models.ServiceStateReady
			byteNewValue, err := json.Marshal(newState)
			So(err, ShouldBeNil)

			oldState := models.ServiceStateDeploying
			byteOldValue, err := json.Marshal(oldState)
			So(err, ShouldBeNil)

			patches, err := MakePatchesForOfferingStateUpdate(oldState, newState)
			So(err, ShouldBeNil)
			So(len(patches), ShouldEqual, 1)
			So(patches[0].Field, ShouldEqual, "State")
			So(patches[0].Operation, ShouldEqual, models.OperationUpdate)
			So(patches[0].Value, ShouldResemble, json.RawMessage(byteNewValue))
			So(patches[0].PrevValue, ShouldResemble, json.RawMessage(byteOldValue))
		})
	})
}
