/**
 * Copyright (c) 2017 Intel Corporation
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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestValidateApplicationCreate(t *testing.T) {
	Convey("Test ValidateApplicationStructCreate", t, func() {

		propperApplication := &Application{TemplateId: "1", Name: "propername"}

		Convey("should return error when ID is provided", func() {
			propperApplication.Id = "1"
			err := propperApplication.ValidateApplicationStructCreate()
			So(err, ShouldNotBeNil)
		})

		Convey("should return error when TemplateID is not provided", func() {
			propperApplication.TemplateId = ""
			err := propperApplication.ValidateApplicationStructCreate()
			So(err, ShouldNotBeNil)
		})

		Convey("shouldn't return error for proper application", func() {
			err := propperApplication.ValidateApplicationStructCreate()
			So(err, ShouldBeNil)
		})

		Convey("shouldn't return error when name has - ", func() {
			propperApplication.Name = "no-proper-name"
			err := propperApplication.ValidateApplicationStructCreate()
			So(err, ShouldBeNil)
		})

		Convey("should return error when name has _ ", func() {
			propperApplication.Name = "no_proper_name"
			err := propperApplication.ValidateApplicationStructCreate()
			So(err, ShouldNotBeNil)
		})
	})
}
