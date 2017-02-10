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
package data

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics-ng/tap-catalog/models"
)

func TestToKeyValueByPatches(t *testing.T) {
	c := DataMapper{}
	Convey("Should return error if field not provided", t, func() {
		patches := []models.Patch{}
		json.Unmarshal([]byte(`[{"value": "IN_PROGRESS", "op": "Update"}]`), &patches)

		_, err := c.ToKeyValueByPatches("a", models.Template{}, patches)
		So(err.Error(), ShouldEqual, "field field is empty!")
	})
	Convey("Should return error if field empty", t, func() {
		patches := []models.Patch{}
		json.Unmarshal([]byte(`[{"field":"", "value": "IN_PROGRESS", "op": "Update"}]`), &patches)

		_, err := c.ToKeyValueByPatches("a", models.Template{}, patches)
		So(err.Error(), ShouldEqual, "field field is empty!")
	})
	Convey("Should return error if value not provided", t, func() {
		patches := []models.Patch{}
		json.Unmarshal([]byte(`[{"field":"Test", "op": "Update"}]`), &patches)

		_, err := c.ToKeyValueByPatches("a", models.Template{}, patches)
		So(err.Error(), ShouldEqual, "field value is empty!")
	})
	Convey("Given proper request", t, func() {
		patches := []models.Patch{}
		json.Unmarshal([]byte(`[{"field":"State", "value":"newTest", "op": "Update"}]`), &patches)

		patchedKeys, err := c.ToKeyValueByPatches("Test/State", models.Template{}, patches)
		Convey("Should not return error", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should return updated value", func() {
			So(patchedKeys.Update[0].Value, ShouldEqual, "newTest")
		})
	})
}
