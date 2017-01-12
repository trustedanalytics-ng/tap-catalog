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

func TestCheckIfMatchingRegexp(t *testing.T) {
	Convey("testing CheckIfMatchingRegexp response", t, func() {
		goodExample := CheckIfMatchingRegexp("GOOD_EXAMPLE", RegexpProperSystemEnvName)
		So(goodExample, ShouldBeNil)

		goodExample = CheckIfMatchingRegexp("good_example", RegexpProperSystemEnvName)
		So(goodExample, ShouldBeNil)

		wrongExample := CheckIfMatchingRegexp("bad-example", RegexpProperSystemEnvName)
		So(wrongExample, ShouldNotBeNil)

		goodExample = CheckIfMatchingRegexp("good-example", RegexpDnsLabelLowercase)
		So(goodExample, ShouldBeNil)

		goodExample = CheckIfMatchingRegexp("good-example", RegexpDnsLabelLowercase)
		So(goodExample, ShouldBeNil)

		wrongExample = CheckIfMatchingRegexp("bad_example", RegexpDnsLabelLowercase)
		So(wrongExample, ShouldNotBeNil)

		wrongExample = CheckIfMatchingRegexp("BAD-example", RegexpDnsLabelLowercase)
		So(wrongExample, ShouldNotBeNil)

	})
}
