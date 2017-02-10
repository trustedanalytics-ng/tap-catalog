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
	"strings"

	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
)

func ApplyApplicationFilter(dataList []interface{}, filter *commonHttp.ItemFilter) ([]catalogModels.Application, error) {
	applications := []catalogModels.Application{}
	for _, item := range dataList {
		application, ok := item.(catalogModels.Application)
		if !ok {
			err := fmt.Errorf("cannot convert element %v to models.Application", item)
			return applications, err
		}

		if filter.Name != "" && strings.ToUpper(filter.Name) != strings.ToUpper(application.Name) {
			continue
		}
		applications = append(applications, application)
	}

	if filter.Name != "" && len(applications) > 0 {
		return applications, nil
	}

	if filter.Skip > 0 {
		if filter.Skip > len(applications) {
			filter.Skip = len(applications)
		}
		applications = append(applications[filter.Skip:])
	}
	if filter.Limit > 0 {
		if filter.Limit > len(applications) {
			filter.Limit = len(applications)
		}
		applications = append(applications[:filter.Limit])
	}

	return applications, nil
}
