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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateServiceStructCreate(t *testing.T) {
	Convey("Test ValidateServiceStructCreate", t, func() {

		propperService := &Service{
			Name: "propername",
			Plans: []ServicePlan{
				{Name: "Test Plan"},
			}}

		Convey("should return error when ID is provided", func() {
			propperService.Id = "1"
			err := propperService.ValidateServiceStructCreate()
			So(err, ShouldNotBeNil)
		})

		Convey("should return error when Plan is not provided", func() {
			propperService.Plans = []ServicePlan{}
			err := propperService.ValidateServiceStructCreate()
			So(err, ShouldNotBeNil)
		})

		Convey("shouldn't return error for proper service", func() {
			err := propperService.ValidateServiceStructCreate()
			So(err, ShouldBeNil)
		})

		Convey("shouldn't return error when name has - ", func() {
			propperService.Name = "no-proper-name"
			err := propperService.ValidateServiceStructCreate()
			So(err, ShouldBeNil)
		})

		Convey("should return error when name has _ ", func() {
			propperService.Name = "no_proper_name"
			err := propperService.ValidateServiceStructCreate()
			So(err, ShouldNotBeNil)
		})
	})
}
