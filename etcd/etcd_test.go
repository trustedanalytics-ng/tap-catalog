// +build dev_test

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

package etcd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEtcdKVStore(t *testing.T) {
	client := EtcdConnector{}
	Convey("Test EtcdKVStore", t, func() {
		Convey("Should update properly", func() {
			err := client.Set("test", "testValue", "", 0)
			So(err, ShouldBeNil)
		})

		Convey("Should get KV properly", func() {
			value, err := client.GetKeyValue("test")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "testValue")
		})

		Convey("Should get KV into struct properly", func() {
			result := ""
			err := client.GetKeyIntoStruct("test", &result)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "testValue")
		})

		Convey("Should delete KV properly", func() {
			err := client.Delete("test", 0)
			So(err, ShouldBeNil)
		})
	})
}
