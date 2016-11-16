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
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/models"
	mutils "github.com/trustedanalytics/tap-metrics/utils"
)

const (
	applicationsMetricName      = "applications"
	servicesMetricName          = "services"
	servicesInstancesMetricName = "serviceInstances"
)

var tapCounts = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "tap",
		Subsystem: "catalog",
		Name:      "counts",
		Help:      "Count of various TAP components",
	}, []string{"component", "organization"})

var repository data.RepositoryApi

func collectApplicationsCount(org string) (string, float64, error) {
	applicationsCount, err := repository.GetDataCounter(data.GetEntityKey(org, data.Applications), models.Application{})
	if err != nil {
		return "", 0, err
	}
	return applicationsMetricName, float64(applicationsCount), nil
}

func collectServicesCount(org string) (string, float64, error) {
	servicesCount, err := repository.GetDataCounter(data.GetEntityKey(org, data.Services), models.Service{})
	if err != nil {
		return "", 0, err
	}
	return servicesMetricName, float64(servicesCount), nil
}

func collectServiceInstancesCount(org string) (string, float64, error) {
	serviceInstances, err := data.GetFilteredInstances(models.InstanceTypeService, "", org, repository)
	if err != nil {
		return "", 0, err
	}
	return servicesInstancesMetricName, float64(len(serviceInstances)), nil
}

func collectCount() error {
	organizations, err := getAllOrgs()
	if err != nil {
		return err
	}
	for _, org := range organizations {
		metricName, metricValue, err := collectApplicationsCount(org)
		if err != nil {
			return err
		}
		tapCounts.WithLabelValues(metricName, org).Set(metricValue)

		metricName, metricValue, err = collectServicesCount(org)
		if err != nil {
			return err
		}
		tapCounts.WithLabelValues(metricName, org).Set(metricValue)

		metricName, metricValue, err = collectServiceInstancesCount(org)
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
