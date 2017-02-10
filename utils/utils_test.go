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
package utils

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func makeAbstractApps(apps []catalogModels.Application) []interface{} {
	abstractApps := make([]interface{}, len(apps))
	for i, filter := range apps {
		abstractApps[i] = filter
	}
	return abstractApps
}

func TestApplyApplicationFilter(t *testing.T) {

	testCases := []struct {
		apps     []catalogModels.Application
		filter   commonHttp.ItemFilter
		expected []catalogModels.Application
	}{
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
			},
			commonHttp.ItemFilter{},
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
			},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
			},
			commonHttp.ItemFilter{Name: "app"},
			[]catalogModels.Application{},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
			},
			commonHttp.ItemFilter{Name: "app1"},
			[]catalogModels.Application{
				{Name: "app1"},
			},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
				{Name: "app3"},
				{Name: "app4"},
			},
			commonHttp.ItemFilter{Limit: 2},
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
			},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
				{Name: "app3"},
				{Name: "app4"},
			},
			commonHttp.ItemFilter{Skip: 3},
			[]catalogModels.Application{
				{Name: "app4"},
			},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
				{Name: "app3"},
				{Name: "app4"},
			},
			commonHttp.ItemFilter{Limit: 2, Skip: 1},
			[]catalogModels.Application{
				{Name: "app2"},
				{Name: "app3"},
			},
		},
		{
			[]catalogModels.Application{
				{Name: "app1"},
				{Name: "app2"},
				{Name: "app3"},
				{Name: "app4"},
			},
			commonHttp.ItemFilter{Name: "app1", Limit: 2, Skip: 1},
			[]catalogModels.Application{
				{Name: "app1"},
			},
		},
	}

	Convey("For a list of test cases", t, func() {
		for _, testCase := range testCases {
			Convey(fmt.Sprintf("Given an application list '%v' and an instance filter '%v' TAP should return a proper application list", testCase.apps, testCase.filter), func() {
				filteredApplications, err := ApplyApplicationFilter(makeAbstractApps(testCase.apps), &testCase.filter)
				So(err, ShouldBeNil)
				So(filteredApplications, ShouldResemble, testCase.expected)
			})
		}
	})
}
