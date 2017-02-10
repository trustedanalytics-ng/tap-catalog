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
package metrics

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/trustedanalytics-ng/tap-catalog/data"
	"github.com/trustedanalytics-ng/tap-catalog/models"
	mutils "github.com/trustedanalytics-ng/tap-metrics/utils"
)

const (
	applicationsMetricName             = "applications"
	servicesMetricName                 = "services"
	servicesInstancesMetricName        = "serviceInstances"
	servicesInstancesRunningMetricName = "serviceInstancesRunning"
	servicesInstancesDownMetricName    = "serviceInstancesDown"
	applicationsRunningMetricName      = "applicationsRunning"
	applicationsDownMetricName         = "applicationsDown"
)

var tapCounts = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "tap",
		Subsystem: "catalog",
		Name:      "counts",
		Help:      "Count of various TAP components",
	}, []string{"component", "organization"})

var repository data.RepositoryApi

func collectInstancesCount(org string) (runningApplications float64, downApplications float64,
	runningServiceInstances float64, downServiceInstances float64, err error) {
	runningApplications = float64(0)
	downApplications = float64(0)
	runningServiceInstances = float64(0)
	downServiceInstances = float64(0)

	result, err := repository.GetListOfData(data.GetEntityKey(org, data.Instances), models.Instance{})
	if err != nil {
		return
	}
	for _, el := range result {
		instance, ok := el.(models.Instance)
		if !ok {
			elementType := reflect.TypeOf(el).String()
			err = fmt.Errorf("Cannot convert element to models.Instance, element type was: %v", elementType)
			return
		}

		if data.IsInstanceTypeOf(instance, models.InstanceTypeApplication) {
			if data.IsRunningInstance(instance) {
				runningApplications = runningApplications + 1
			} else {
				downApplications = downApplications + 1
			}
		} else if data.IsInstanceTypeOf(instance, models.InstanceTypeService) {
			if data.IsRunningInstance(instance) {
				runningServiceInstances = runningServiceInstances + 1
			} else {
				downServiceInstances = downServiceInstances + 1
			}
		}
	}

	return
}

func collectServicesCount(org string) (string, float64, error) {
	servicesCount, err := repository.GetDataCounter(data.GetEntityKey(org, data.Services), models.Service{})
	if err != nil {
		return "", 0, err
	}
	return servicesMetricName, float64(servicesCount), nil
}

func collectCount() error {
	organizations, err := getAllOrgs()
	if err != nil {
		return err
	}
	for _, org := range organizations {
		runningApplications, downApplications, runningServiceInstances, downServiceInstances, err := collectInstancesCount(org)
		if err != nil {
			return err
		}
		tapCounts.WithLabelValues(applicationsMetricName, org).Set(float64(runningApplications + downApplications))
		tapCounts.WithLabelValues(applicationsRunningMetricName, org).Set(runningApplications)
		tapCounts.WithLabelValues(applicationsDownMetricName, org).Set(downApplications)

		tapCounts.WithLabelValues(servicesInstancesMetricName, org).Set(float64(runningServiceInstances + downServiceInstances))
		tapCounts.WithLabelValues(servicesInstancesRunningMetricName, org).Set(runningServiceInstances)
		tapCounts.WithLabelValues(servicesInstancesDownMetricName, org).Set(downServiceInstances)

		metricName, metricValue, err := collectServicesCount(org)
		if err != nil {
			return err
		}
		tapCounts.WithLabelValues(metricName, org).Set(metricValue)
	}
	return nil
}

func EnableCollection(repo data.RepositoryApi, delay time.Duration) chan<- struct{} {
	repository = repo
	mutils.RegisterMetrics("catalog", tapCounts)
	return mutils.EnableMetricsCollecting(delay,
		collectCount,
	)
}

func getAllOrgs() ([]string, error) {
	/*
		For more orgs, they should be get from user-management
	*/
	coreOrg := os.Getenv("CORE_ORGANIZATION")
	if coreOrg == "" {
		return nil, errors.New("CORE_ORGANIZATION env is empty")
	}
	orgs := []string{coreOrg}

	return orgs, nil
}
