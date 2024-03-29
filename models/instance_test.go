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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetValueFromMetadata(t *testing.T) {
	Convey("Test GetValueFromMetadata should return value for existing key", t, func() {
		metadatas := []Metadata{{Id: "mykey", Value: "myvalue"}}
		value := GetValueFromMetadata(metadatas, "mykey")
		So(value, ShouldEqual, "myvalue")
	})

	Convey("Test GetValueFromMetadata should return empty string for not existing key", t, func() {
		metadatas := []Metadata{{Id: "mykey", Value: "myvalue"}}
		value := GetValueFromMetadata(metadatas, "nokey")
		So(value, ShouldEqual, "")
	})
}

func TestValidateInstanceStructCreate(t *testing.T) {
	Convey("Test ValidateInstanceStructCreate", t, func() {

		propperInstance := &Instance{
			Name: "propername",
			Metadata: []Metadata{
				{Id: "PLAN_ID", Value: "1"},
			},
		}

		Convey("should return error when ID is provided", func() {
			propperInstance.Id = "1"
			err := propperInstance.ValidateInstanceStructCreate(InstanceTypeService)
			So(err, ShouldNotBeNil)
		})

		Convey("shouldn't return error for proper instance", func() {
			err := propperInstance.ValidateInstanceStructCreate(InstanceTypeService)
			So(err, ShouldBeNil)
		})

		Convey("shouldn't return error when name has - ", func() {
			propperInstance.Name = "no-proper-name"
			err := propperInstance.ValidateInstanceStructCreate(InstanceTypeService)
			So(err, ShouldBeNil)
		})

		Convey("should return error when name has _ ", func() {
			propperInstance.Name = "no_proper_name"
			err := propperInstance.ValidateInstanceStructCreate(InstanceTypeService)
			So(err, ShouldNotBeNil)
		})
	})
}
