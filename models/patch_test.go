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
package models

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidatePatch(t *testing.T) {
	Convey("Should return error when Field is null or empty", t, func() {
		patch := Patch{}
		json.Unmarshal([]byte(`{"value": "IN_PROGRESS", "op": "Update"}`), &patch)

		err := ValidatePatchStructure(patch)
		So(err.Error(), ShouldEqual, "field field is empty!")
	})
	Convey("Should return error when Value is null", t, func() {
		patch := Patch{}
		json.Unmarshal([]byte(`{"field": "State", "op": "Update"}`), &patch)

		err := ValidatePatchStructure(patch)
		So(err.Error(), ShouldEqual, "field value is empty!")
	})
	Convey("Should return error when Value is empty", t, func() {
		patch := Patch{}
		json.Unmarshal([]byte(`{"field": "State", "value":"","op": "Update"}`), &patch)

		err := ValidatePatchStructure(patch)
		So(err, ShouldEqual, nil)
	})
	Convey("Should not return error when all fields are provided", t, func() {
		patch := Patch{}
		json.Unmarshal([]byte(`{"value": "IN_PROGRESS", "op": "Update", "field":"state"}`), &patch)

		err := ValidatePatchStructure(patch)
		So(err, ShouldEqual, nil)
	})

}
